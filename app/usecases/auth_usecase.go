package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"linkshare/app/dto"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/models"
	"linkshare/app/repositories"
	"linkshare/app/repositories/mongo_repo"
	"linkshare/app/repositories/sql_repo"
	"linkshare/app/security"
	"linkshare/generated"
	"net/http"
	"sync"
	"time"
)

type IAuthUseCase interface {
	Register(ctx context.Context, req *generated.RegisterRequest) *model.ErrorLog
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, *model.ErrorLog)
	validateUser(ctx context.Context, users *dto.LoginRequest) (*models.Users, *model.ErrorLog)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, *model.ErrorLog)
}

type authUseCase struct {
	generalRepository repositories.IGeneralRepository
	authRepository    sql_repo.IAuthRepository
	accessRepository  mongo_repo.IAccessTokenSessionsRepository
	refreshRepository mongo_repo.IRefreshTokenSessionsRepository
	jwtSecurity       security.IJwtSecurity
}

func NewAuthUseCase(generalRepository repositories.IGeneralRepository, authRepository sql_repo.IAuthRepository, accessRepository mongo_repo.IAccessTokenSessionsRepository, refreshRepository mongo_repo.IRefreshTokenSessionsRepository, jwtSecurity security.IJwtSecurity) IAuthUseCase {
	return &authUseCase{
		generalRepository: generalRepository,
		authRepository:    authRepository,
		accessRepository:  accessRepository,
		refreshRepository: refreshRepository,
		jwtSecurity:       jwtSecurity,
	}
}

func (u *authUseCase) Register(ctx context.Context, req *generated.RegisterRequest) *model.ErrorLog {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error while hashing password")
		return errLog
	}
	req.Password = string(bcryptPassword)
	tx := u.generalRepository.BeginTransaction(ctx)
	if tx == nil {
		return helper.WriteLog(errors.New("error while begin transaction"), http.StatusInternalServerError, "fail to register")
	}

	errLog := u.authRepository.Create(tx, req, ctx)
	if errLog != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return helper.WriteLog(err, http.StatusInternalServerError, fmt.Sprintf("error while rollback transaction, after got error create: %s", errLog.Err.Error()))
		}
		return errLog
	}
	err = tx.Commit(ctx)
	if err != nil {
		return helper.WriteLog(err, http.StatusInternalServerError, "error while commit transaction")
	}
	return nil
}

func (u *authUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, *model.ErrorLog) {
	response := &dto.LoginResponse{}
	user, errLog := u.validateUser(ctx, req)
	if errLog != nil {
		if errors.Is(errLog.Err, pgx.ErrNoRows) {
			return nil, helper.WriteLog(errors.New("username or email or password is not valid"), http.StatusUnauthorized, "")
		}
		return nil, errLog
	}
	if user == nil {
		return nil, helper.WriteLog(errors.New("username or email or password is not valid"), http.StatusUnauthorized, "")
	}
	userRequest := &models.UserDataOnJWT{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
	tokenResult, err := u.jwtSecurity.GenerateToken(userRequest)
	if err != nil {
		errLog = helper.WriteLog(err, http.StatusInternalServerError, "")
		return nil, errLog
	}

	// insert access token to database
	requestAccessToken := &models.AccessTokenSession{
		AccessToken: tokenResult.AccessToken,
		Expired:     tokenResult.AccessTokenExpired,
		UserData:    userRequest,
	}
	go u.accessRepository.Insert(requestAccessToken, ctx)

	redisKeyForAccessToken := helper.GenRedisKeyAccessTokenSessionByUserID(user.Id)
	u.generalRepository.SetRedisCache(ctx, redisKeyForAccessToken, requestAccessToken, requestAccessToken.Expired.Sub(time.Now()))

	// insert refresh token to database
	requestRefreshToken := &models.RefreshTokenSession{
		RefreshToken: tokenResult.RefreshToken,
		Expired:      tokenResult.RefreshTokenExpired,
		UserData:     userRequest,
	}
	go u.refreshRepository.Insert(requestRefreshToken, ctx)

	// set response
	response.AccessToken = tokenResult.AccessToken
	response.RefreshToken = tokenResult.RefreshToken

	// set redis cache for refresh token
	redisKeyForRefreshToken := helper.GenRedisKeyRefreshTokenSessionByUserID(user.Id)
	u.generalRepository.SetRedisCache(ctx, redisKeyForRefreshToken, requestRefreshToken, requestRefreshToken.Expired.Sub(time.Now()))

	return response, nil
}

func (u *authUseCase) validateUser(ctx context.Context, request *dto.LoginRequest) (*models.Users, *model.ErrorLog) {
	user, errLog := u.authRepository.GetUserByUsernameOrEmail(request.Username, request.Email, ctx)
	if errLog != nil {
		return nil, errLog
	}

	// validate password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		errLog = helper.WriteLog(errors.New("username or email or password is not valid"), 401, "")
		return nil, errLog
	}
	return user, nil
}

func (u *authUseCase) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, *model.ErrorLog) {
	response := &dto.LoginResponse{
		RefreshToken: req.RefreshToken,
	}

	userData, err := helper.GetUserDataOnCtx(ctx)
	if err != nil {
		return nil, helper.WriteLog(err, http.StatusInternalServerError, "error while get user data on context")
	}

	res := &dto.GetByRefreshTokenChan{}
	redisKey := helper.GenRedisKeyRefreshTokenSessionByUserID(userData.Id)
	u.generalRepository.GetRedisCache(ctx, redisKey, &res.Data)
	if res.Data == nil {
		res.Data, res.ErrLog = u.refreshRepository.GetByRefreshToken(req.RefreshToken, ctx)
		if res.ErrLog != nil {
			if res.ErrLog.StatusCode == http.StatusNotFound {
				res.ErrLog.Err = errors.New("unauthorized")
				res.ErrLog.Message = "refresh token expired"
				res.ErrLog.SystemMessage = ""
			}
			return nil, res.ErrLog
		}
	}
	if res.Data.Expired.Before(time.Now()) {
		return nil, helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "refresh token expired")
	}

	token, err := u.jwtSecurity.GenerateAccessToken(userData)
	if err != nil {
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error generating access token")
		return nil, errLog
	}
	response.AccessToken = token.AccessToken
	requestAccessToken := &models.AccessTokenSession{
		AccessToken: token.AccessToken,
		Expired:     token.AccessTokenExpired,
		UserData:    userData,
	}
	u.clearAndInsertAccessToken(ctx, userData.Id, requestAccessToken)
	return response, nil
}

func (u *authUseCase) clearAndInsertAccessToken(ctx context.Context, userId int64, dataAccess *models.AccessTokenSession) {
	wg := &sync.WaitGroup{}
	// delete access session on redis and mongo
	wg.Add(2)
	go func() {
		defer wg.Done()
		u.generalRepository.DelRedisCache(ctx, helper.GenRedisKeyAccessTokenSessionByUserID(userId))
	}()
	go func() {
		defer wg.Done()
		u.accessRepository.DeleteByUserId(ctx, userId)
	}()
	wg.Wait()

	// set access session on redis and mongo
	wg.Add(2)
	go func() {
		defer wg.Done()
		redisKeyForAccessToken := helper.GenRedisKeyAccessTokenSessionByUserID(userId)
		u.generalRepository.SetRedisCache(ctx, redisKeyForAccessToken, dataAccess, dataAccess.Expired.Sub(time.Now()))
	}()
	go func() {
		defer wg.Done()
		u.accessRepository.Insert(dataAccess, ctx)
	}()
	wg.Wait()
}

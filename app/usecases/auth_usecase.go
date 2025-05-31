package usecases

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"linkshare/app/dto"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/models"
	"linkshare/app/repositories/mongo_repo"
	"linkshare/app/repositories/sql_repo"
	"linkshare/app/security"
	"net/http"
)

type IAuthUseCase interface {
	Login(req *dto.LoginRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog)
	ValidateUser(users *dto.LoginRequest, ctx context.Context) (*models.Users, *model.ErrorLog)
	RefreshToken(req *dto.RefreshTokenRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog)
}

type authUseCase struct {
	authRepository    sql_repo.IAuthRepository
	accessRepository  mongo_repo.IAccessTokenSessionsRepository
	refreshRepository mongo_repo.IRefreshTokenSessionsRepository
	jwtSecurity       security.IJwtSecurity
}

func NewAuthUseCase(authRepository sql_repo.IAuthRepository, accessRepository mongo_repo.IAccessTokenSessionsRepository, refreshRepository mongo_repo.IRefreshTokenSessionsRepository, jwtSecurity security.IJwtSecurity) IAuthUseCase {
	return &authUseCase{
		authRepository:    authRepository,
		accessRepository:  accessRepository,
		refreshRepository: refreshRepository,
		jwtSecurity:       jwtSecurity,
	}
}

func (u *authUseCase) Login(req *dto.LoginRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog) {
	response := &dto.LoginResponse{}
	user, errLog := u.ValidateUser(req, ctx)
	if errLog != nil {
		return nil, errLog
	}
	if user == nil {
		errLog = helper.WriteLog(errors.New("username or email or password is not valid"), http.StatusUnauthorized, "")
		return nil, errLog
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
	requestAccessToken := &models.AccessTokenSession{
		AccessToken: tokenResult.AccessToken,
		Expired:     tokenResult.AccessTokenExpired,
		UserData:    userRequest,
	}
	go u.accessRepository.Insert(requestAccessToken, ctx)
	requestRefreshToken := &models.RefreshTokenSession{
		RefreshToken: tokenResult.RefreshToken,
		Expired:      tokenResult.RefreshTokenExpired,
		UserData:     userRequest,
	}
	go u.refreshRepository.Insert(requestRefreshToken, ctx)
	response.AccessToken = tokenResult.AccessToken
	response.RefreshToken = tokenResult.RefreshToken
	return response, nil
}

func (u *authUseCase) ValidateUser(request *dto.LoginRequest, ctx context.Context) (*models.Users, *model.ErrorLog) {
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

func (u *authUseCase) RefreshToken(req *dto.RefreshTokenRequest, ctx context.Context) (*dto.LoginResponse, *model.ErrorLog) {
	response := &dto.LoginResponse{
		RefreshToken: req.RefreshToken,
	}
	responseChan := make(chan *dto.GetByRefreshTokenChan)
	go func(responseChan chan *dto.GetByRefreshTokenChan) {
		res := &dto.GetByRefreshTokenChan{}
		res.Data, res.ErrLog = u.refreshRepository.GetByRefreshToken(req.RefreshToken, ctx)
		responseChan <- res
	}(responseChan)

	resultValidate := u.jwtSecurity.ValidateRefreshToken(req.RefreshToken)
	if resultValidate.Error != nil {
		if errors.Is(resultValidate.Error, jwt.ErrTokenExpired) {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, resultValidate.Error.Error())
			return nil, errLog
		}
		errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "invalid refresh token")
		return nil, errLog
	}
	responseGetSession := <-responseChan
	if responseGetSession.ErrLog != nil {
		return nil, responseGetSession.ErrLog
	}
	token, err := u.jwtSecurity.GenerateAccessToken(responseGetSession.Data.UserData)
	if err != nil {
		errLog := helper.WriteLog(err, http.StatusInternalServerError, "error generating access token")
		return nil, errLog
	}
	response.AccessToken = token.AccessToken

	requestRefreshToken := &models.AccessTokenSession{
		AccessToken: token.AccessToken,
		Expired:     token.AccessTokenExpired,
		UserData:    responseGetSession.Data.UserData,
	}
	go u.accessRepository.Insert(requestRefreshToken, ctx)
	return response, nil
}

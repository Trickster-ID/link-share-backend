package usecases

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"linkshare/app/configuration"
	"linkshare/app/dto"
	"linkshare/app/global/helper"
	mocksMongo "linkshare/app/mocks/app/repositories/mongo_repo"
	mocksSql "linkshare/app/mocks/app/repositories/sql_repo"
	mockSecurity "linkshare/app/mocks/app/security"
	"linkshare/app/models"
	"net/http"
	"time"

	"testing"
)

func TestLogin(t *testing.T) {
	configuration.InitialConfigForUnitTest()
	mockAuthRepo := mocksSql.NewMockIAuthRepository(t)
	mockAccessRepo := mocksMongo.NewMockIAccessTokenSessionsRepository(t)
	mockRefreshRepo := mocksMongo.NewMockIRefreshTokenSessionsRepository(t)
	mockJwtSecurity := mockSecurity.NewMockIJwtSecurity(t)
	useCase := NewAuthUseCase(mockAuthRepo, mockAccessRepo, mockRefreshRepo, mockJwtSecurity)
	ctx := context.TODO()
	timeNow := time.Now()
	mockAccessRepo.On("Insert", mock.AnythingOfType("*models.AccessTokenSession"), ctx).Return()
	mockRefreshRepo.On("Insert", mock.AnythingOfType("*models.RefreshTokenSession"), ctx).Return()

	t.Run("+: admin username", func(t *testing.T) {
		returnValue := &models.Users{
			Id:           1,
			Username:     "admin",
			Email:        "admin@live.com",
			PasswordHash: "$2a$04$ZP1.DVAdR677eHTBUpDzE.0hHnp31JcyRK/eMF9Z7Y.iOWJkE/JNi",
			RoleID:       1,
			CreatedAt:    nil,
			UpdatedAt:    nil,
		}
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("admin", "", ctx).Return(returnValue, nil)
		userRequest := &models.UserDataOnJWT{
			Id:       1,
			Username: "admin",
			Email:    "admin@live.com",
		}
		generateTokenResponse := &dto.GenerateTokenResponse{
			AccessToken:         "accessToken123",
			AccessTokenExpired:  timeNow,
			RefreshToken:        "refreshToken123",
			RefreshTokenExpired: timeNow,
		}
		mockJwtSecurity.EXPECT().GenerateToken(userRequest).Return(generateTokenResponse, nil)
		request := &dto.LoginRequest{
			Username: "admin",
			Email:    "",
			Password: "admin",
		}
		response, errLog := useCase.Login(request, ctx)
		assert.Nil(t, errLog)
		assert.NotNil(t, response)
		assert.NotNil(t, response.AccessToken)
		assert.NotNil(t, response.RefreshToken)
	})
	t.Run("+: admin email", func(t *testing.T) {
		returnValue := &models.Users{
			Id:           1,
			Username:     "admin",
			Email:        "admin@live.com",
			PasswordHash: "$2a$04$ZP1.DVAdR677eHTBUpDzE.0hHnp31JcyRK/eMF9Z7Y.iOWJkE/JNi",
			RoleID:       1,
			CreatedAt:    nil,
			UpdatedAt:    nil,
		}
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("", "admin@live.com", ctx).Return(returnValue, nil)
		userRequest := &models.UserDataOnJWT{
			Id:       1,
			Username: "admin",
			Email:    "admin@live.com",
		}
		generateTokenResponse := &dto.GenerateTokenResponse{
			AccessToken:         "accessToken123",
			AccessTokenExpired:  timeNow,
			RefreshToken:        "refreshToken123",
			RefreshTokenExpired: timeNow,
		}
		mockJwtSecurity.EXPECT().GenerateToken(userRequest).Return(generateTokenResponse, nil)
		request := &dto.LoginRequest{
			Username: "",
			Email:    "admin@live.com",
			Password: "admin",
		}
		response, errLog := useCase.Login(request, ctx)
		assert.Nil(t, errLog)
		assert.NotNil(t, response)
		assert.NotNil(t, response.AccessToken)
		assert.NotNil(t, response.RefreshToken)
	})
	t.Run("-: wrong password", func(t *testing.T) {
		returnValue := &models.Users{
			Id:           1,
			Username:     "admin",
			Email:        "admin@live.com",
			PasswordHash: "$2a$04$ZP1.DVAdR677eHTBUpDzE.0hHnp31JcyRK/eMF9Z7Y.iOWJkE/JNi",
			RoleID:       1,
			CreatedAt:    nil,
			UpdatedAt:    nil,
		}
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("admin", "", ctx).Return(returnValue, nil)
		request := &dto.LoginRequest{
			Username: "admin",
			Email:    "",
			Password: "admin123",
		}
		response, errLog := useCase.Login(request, ctx)
		assert.NotNil(t, errLog)
		assert.Nil(t, response)
		assert.Equal(t, errors.New("username or email or password is not valid"), errLog.Err)
		assert.Equal(t, http.StatusUnauthorized, errLog.StatusCode)
		assert.Equal(t, "", errLog.Message)
	})
	t.Run("-: wrong username", func(t *testing.T) {
		errorLogExpected := helper.WriteLogWoP(pgx.ErrNoRows, 404, "please enter valid username or email or password")
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("admin123", "", ctx).Return(nil, errorLogExpected)
		request := &dto.LoginRequest{
			Username: "admin123",
			Email:    "",
			Password: "admin",
		}
		response, errLog := useCase.Login(request, ctx)
		assert.NotNil(t, errLog)
		assert.Nil(t, response)
		assert.Equal(t, pgx.ErrNoRows, errLog.Err)
		assert.Equal(t, http.StatusNotFound, errLog.StatusCode)
		assert.Equal(t, "please enter valid username or email or password", errLog.Message)
	})
	assert.True(t, mockAuthRepo.AssertExpectations(t))
	assert.True(t, mockAccessRepo.AssertExpectations(t))
	assert.True(t, mockRefreshRepo.AssertExpectations(t))
}

func TestValidateUser(t *testing.T) {
	configuration.InitialConfigForUnitTest()
	mockAuthRepo := mocksSql.NewMockIAuthRepository(t)
	mockAccessRepo := mocksMongo.NewMockIAccessTokenSessionsRepository(t)
	mockRefreshRepo := mocksMongo.NewMockIRefreshTokenSessionsRepository(t)
	mockJwtSecurity := mockSecurity.NewMockIJwtSecurity(t)
	useCase := NewAuthUseCase(mockAuthRepo, mockAccessRepo, mockRefreshRepo, mockJwtSecurity)
	ctx := context.Background()

	t.Run("+: admin", func(t *testing.T) {
		returnValue := &models.Users{
			Id:           1,
			Username:     "admin",
			Email:        "admin@live.com",
			PasswordHash: "$2a$04$ZP1.DVAdR677eHTBUpDzE.0hHnp31JcyRK/eMF9Z7Y.iOWJkE/JNi",
			RoleID:       1,
			CreatedAt:    nil,
			UpdatedAt:    nil,
		}
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("admin", "", context.Background()).Return(returnValue, nil)
		request := &dto.LoginRequest{
			Username: "admin",
			Email:    "",
			Password: "admin",
		}
		response, errLog := useCase.validateUser(ctx, request)
		assert.Nil(t, errLog)
		assert.NotNil(t, response)
		assert.Equal(t, returnValue, response)
	})

	t.Run("-: does not exist user", func(t *testing.T) {
		mockAuthRepo.EXPECT().GetUserByUsernameOrEmail("userNotExist", "", context.Background()).Return(nil, helper.WriteLogWoP(pgx.ErrNoRows, 404, "please enter valid username"))
		request := &dto.LoginRequest{
			Username: "userNotExist",
			Email:    "",
			Password: "userNotExist",
		}
		response, errLog := useCase.validateUser(ctx, request)

		expectedError := helper.WriteLogWoP(pgx.ErrNoRows, 404, "please enter valid username")
		assert.Nil(t, response)
		assert.NotNil(t, errLog)
		assert.Equal(t, expectedError, errLog)
	})
}

func TestRefreshToken(t *testing.T) {
	configuration.InitialConfigForUnitTest()
	mockAuthRepo := mocksSql.NewMockIAuthRepository(t)
	mockAccessRepo := mocksMongo.NewMockIAccessTokenSessionsRepository(t)
	mockRefreshRepo := mocksMongo.NewMockIRefreshTokenSessionsRepository(t)
	mockJwtSecurity := mockSecurity.NewMockIJwtSecurity(t)
	useCase := NewAuthUseCase(mockAuthRepo, mockAccessRepo, mockRefreshRepo, mockJwtSecurity)
	ctx := context.Background()
	timeNow := time.Now()
	t.Run("+: admin refresh", func(t *testing.T) {
		objectID, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
		if err != nil {
			t.Fatal("create mongo id error: ", err)
		}
		getRefreshResponse := &models.RefreshTokenSession{
			ObjectID:     objectID,
			RefreshToken: "refreshToken123",
			Expired:      timeNow,
			UserData: &models.UserDataOnJWT{
				Id:       1,
				Username: "admin",
				Email:    "admin@live.com",
			},
		}
		mockRefreshRepo.EXPECT().GetByRefreshToken("token123", ctx).Return(getRefreshResponse, nil)
		validateRefreshTokenResponse := &dto.ValidateTokenResponse{
			User: &models.UserDataOnJWT{
				Id:       1,
				Username: "admin",
				Email:    "admin@live.com",
			},
			Error: nil,
		}
		mockJwtSecurity.EXPECT().ValidateRefreshToken("token123").Return(validateRefreshTokenResponse)
		generateAccessTokenResponse := &dto.GenerateAccessTokenResponse{
			AccessToken:        "accessToken123",
			AccessTokenExpired: timeNow,
		}
		mockJwtSecurity.EXPECT().GenerateAccessToken(validateRefreshTokenResponse.User).Return(generateAccessTokenResponse, nil)
		request := &dto.RefreshTokenRequest{
			RefreshToken: "token123",
		}
		mockAccessRepo.On("Insert", mock.AnythingOfType("*models.AccessTokenSession"), ctx).Return()
		refreshResponse, errLog := useCase.RefreshToken(request, ctx)
		assert.NotNil(t, refreshResponse)
		assert.Nil(t, errLog)
		assert.NotNil(t, refreshResponse.AccessToken)
		assert.NotNil(t, refreshResponse.RefreshToken)
	})
	t.Run("-: wrong refresh token", func(t *testing.T) {
		getByRefreshTokenErrLogExpected := helper.WriteLogWoP(mongo.ErrNoDocuments, http.StatusNotFound, "")
		mockRefreshRepo.EXPECT().GetByRefreshToken("token1234", ctx).Return(nil, getByRefreshTokenErrLogExpected)
		validateRefreshTokenResponse := &dto.ValidateTokenResponse{
			User:  nil,
			Error: errors.New("invalid token"),
		}
		mockJwtSecurity.EXPECT().ValidateRefreshToken("token1234").Return(validateRefreshTokenResponse)
		request := &dto.RefreshTokenRequest{
			RefreshToken: "token1234",
		}
		refreshResponse, errLog := useCase.RefreshToken(request, ctx)
		expextedErrLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "invalid refresh token")
		assert.Nil(t, refreshResponse)
		assert.NotNil(t, errLog)
		assert.Equal(t, expextedErrLog, errLog)
	})
	//time.Sleep(50 * time.Millisecond)
}

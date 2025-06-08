package middleware

import (
	"encoding/base64"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"linkshare/app/dto"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/security"
	"net/http"
	"os"
	"strings"
)

func AccessTokenMiddleware() fiber.Handler {
	return validateToken(true)
}
func RefreshTokenMiddleware() fiber.Handler {
	return validateToken(false)
}

func validateToken(isAccessToken bool) fiber.Handler {
	return func(f *fiber.Ctx) error {
		authHeader := f.Get(fiber.HeaderAuthorization)
		errLog := helper.WriteLogWoP(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
		if authHeader == "" {
			logrus.Trace("auth header is empty")
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog}, http.StatusUnauthorized)
		}

		// Check if the token is in the format "Bearer <token>"
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		jwtSecurity := security.NewJwtSecurity()
		var response *dto.ValidateTokenResponse
		if isAccessToken {
			response = jwtSecurity.ValidateAccessToken(tokenString)
		} else {
			response = jwtSecurity.ValidateRefreshToken(tokenString)
		}
		if response.Error != nil {
			if errors.Is(response.Error, jwt.ErrTokenExpired) {
				errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, response.Error.Error())
				return helper.Response(f, &model.BaseResponse{ErrorLog: errLog}, http.StatusUnauthorized)
			}
			logrus.Trace(response.Error)
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog}, http.StatusUnauthorized)
		}
		err := helper.SetUserDataOnCtx(f, response.User)
		if err != nil {
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog}, http.StatusUnauthorized)
		}
		return f.Next()
	}
}

func BasicAuthMiddleware() fiber.Handler {
	return func(f *fiber.Ctx) error {
		// Get the Authorization header
		auth := f.Get("Authorization")
		if auth == "" {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "No Authorization header")
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
		}

		// Check if it starts with "Basic"
		if !strings.HasPrefix(auth, "Basic ") {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid Authorization header")
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
		}

		// Decode the Base64 encoded credentials
		encodedCredentials := strings.TrimPrefix(auth, "Basic ")
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid credentials encoding")
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
		}

		// Split the decoded credentials into username and password
		credentials := strings.SplitN(string(decodedCredentials), ":", 2)
		if len(credentials) != 2 || credentials[0] != os.Getenv("BASICAUTH_USERNAME") || credentials[1] != os.Getenv("BASICAUTH_PASSWORD") {
			errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid username or password")
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
		}

		// Continue to the next handler if authentication is successful
		return f.Next()
	}
}

package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/security"
	"net/http"
	"os"
	"strings"
)

func TokenMiddleware() fiber.Handler {
	return func(f fiber.Ctx) error {
		authHeader := f.Get("Authorization")
		errLog := helper.WriteLogWoP(errors.New("unauthorized"), http.StatusUnauthorized, "Invalid token")
		if authHeader == "" {
			logrus.Trace("auth header is empty")
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
		}

		// Check if the token is in the format "Bearer <token>"
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		jwtSecurity := security.NewJwtSecurity()
		response := jwtSecurity.ValidateAccessToken(tokenString)
		if response.Error != nil {
			if errors.Is(response.Error, jwt.ErrTokenExpired) {
				errLog := helper.WriteLog(errors.New("unauthorized"), http.StatusUnauthorized, response.Error.Error())
				return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
			}
			logrus.Trace(response.Error)
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
		}
		marshalledData, err := json.Marshal(response.User)
		if err != nil {
			logrus.Trace(err)
			return helper.Response(f, &model.BaseResponse{ErrorLog: errLog})
		}
		f.Set("user_data", string(marshalledData))
		return f.Next()
	}
}

func BasicAuthMiddleware() fiber.Handler {
	return func(f fiber.Ctx) error {
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

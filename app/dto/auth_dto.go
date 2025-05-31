package dto

import (
	"linkshare/app/models"
	"time"
)

type GenerateTokenResponse struct {
	AccessToken         string    `json:"access_token"`
	AccessTokenExpired  time.Time `json:"access_token_expired"`
	RefreshToken        string    `json:"refresh_token"`
	RefreshTokenExpired time.Time `json:"refresh_token_expired"`
}

type GenerateAccessTokenResponse struct {
	AccessToken        string    `json:"access_token"`
	AccessTokenExpired time.Time `json:"access_token_expired"`
}

type LoginRequest struct {
	Username string `form:"username" json:"username"`
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidateTokenResponse struct {
	User  *models.UserDataOnJWT `json:"user"`
	Error error                 `json:"error"`
}

type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token"`
}

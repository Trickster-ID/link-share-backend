package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"linkshare/app/constants"
	"linkshare/app/dto"
	"linkshare/app/models"
	"os"
	"time"
)

type IJwtSecurity interface {
	GenerateToken(userData *models.UserDataOnJWT) (*dto.GenerateTokenResponse, error)
	GenerateAccessToken(userData *models.UserDataOnJWT) (*dto.GenerateAccessTokenResponse, error)
	ValidateAccessToken(accessToken string) *dto.ValidateTokenResponse
	ValidateRefreshToken(refreshToken string) *dto.ValidateTokenResponse
	RefreshToken(refreshToken string) (string, error)
}

type JwtSecurity struct {
	JwtKeyAccessToken   string
	AccessTokenExpired  time.Duration
	JwtKeyRefreshToken  string
	RefreshTokenExpired time.Duration
}

func NewJwtSecurity() IJwtSecurity {
	return &JwtSecurity{
		JwtKeyAccessToken:   os.Getenv("JWT_KEY_ACCESS_TOKEN"),
		AccessTokenExpired:  constants.ACCESS_TOKEN_EXPIRED,
		JwtKeyRefreshToken:  os.Getenv("JWT_KEY_REFRESH_TOKEN"),
		RefreshTokenExpired: constants.REFRESH_TOKEN_EXPIRED,
	}
}

// getToken will generate token with JWT lib
func (s *JwtSecurity) getToken(userData *models.UserDataOnJWT, expiredTime time.Time, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"id":       userData.Id,
		"username": userData.Username,
		"email":    userData.Email,
		"exp":      expiredTime.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "oauth2-pikri",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(secretKey))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return token, nil
}

func (s *JwtSecurity) GenerateToken(userData *models.UserDataOnJWT) (*dto.GenerateTokenResponse, error) {
	response := &dto.GenerateTokenResponse{
		AccessTokenExpired:  time.Now().Add(s.AccessTokenExpired),
		RefreshTokenExpired: time.Now().Add(s.RefreshTokenExpired),
	}
	accessToken, err := s.getToken(userData, response.AccessTokenExpired, s.JwtKeyAccessToken)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	response.AccessToken = accessToken
	refreshToken, err := s.getToken(userData, response.RefreshTokenExpired, s.JwtKeyRefreshToken)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	response.RefreshToken = refreshToken
	return response, nil
}

func (s *JwtSecurity) GenerateAccessToken(userData *models.UserDataOnJWT) (*dto.GenerateAccessTokenResponse, error) {
	response := &dto.GenerateAccessTokenResponse{
		AccessTokenExpired: time.Now().Add(s.AccessTokenExpired),
	}
	accessToken, err := s.getToken(userData, response.AccessTokenExpired, s.JwtKeyAccessToken)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	response.AccessToken = accessToken
	return response, nil
}

func (s *JwtSecurity) ValidateAccessToken(accessToken string) *dto.ValidateTokenResponse {
	return s.validateToken(accessToken, s.JwtKeyAccessToken)
}

func (s *JwtSecurity) ValidateRefreshToken(refreshToken string) *dto.ValidateTokenResponse {
	return s.validateToken(refreshToken, s.JwtKeyRefreshToken)
}

// ValidateToken will validate jwt token for middleware function
func (s *JwtSecurity) validateToken(tokenString, secretKey string) *dto.ValidateTokenResponse {
	response := &dto.ValidateTokenResponse{}
	// Parse the token and validate its signature
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			// Make sure the signing method is ECDSA (ES256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil || !token.Valid {
		response.Error = errors.New("invalid token")
		return response
	}
	expTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		response.Error = err
		return response
	}
	if expTime.Before(time.Now()) {
		response.Error = jwt.ErrTokenExpired
		return response
	}
	userData := &models.UserDataOnJWT{
		Id:       int64(token.Claims.(jwt.MapClaims)["id"].(float64)),
		Username: token.Claims.(jwt.MapClaims)["username"].(string),
		Email:    token.Claims.(jwt.MapClaims)["email"].(string),
	}
	response.User = userData
	return response
}

func (s *JwtSecurity) RefreshToken(refreshToken string) (string, error) {
	response := s.ValidateRefreshToken(refreshToken)
	if response.Error != nil {
		logrus.Error(response.Error)
		return "", response.Error
	}
	newToken, err := s.getToken(response.User, time.Now().Add(time.Hour*24), s.JwtKeyAccessToken)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newToken, nil
}

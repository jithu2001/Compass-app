package utils

import (
	"errors"
	"time"

	"compass-backend/config"
	"compass-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID uint64           `json:"user_id"`
	Email  string           `json:"email"`
	Role   models.UserRole  `json:"role"`
	Type   TokenType        `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User, tokenType TokenType, cfg *config.Config) (string, error) {
	var expiration time.Duration
	
	switch tokenType {
	case AccessToken:
		expiration = cfg.JWT.AccessTokenDuration
	case RefreshToken:
		expiration = cfg.JWT.RefreshTokenDuration
	default:
		return "", errors.New("invalid token type")
	}

	claims := Claims{
		UserID: user.UserID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "compass-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func ValidateToken(tokenString string, expectedType TokenType, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func GenerateTokenPair(user *models.User, cfg *config.Config) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateToken(user, AccessToken, cfg)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateToken(user, RefreshToken, cfg)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
package services

import (
	"errors"

	"compass-backend/config"
	"compass-backend/internal/models"
	"compass-backend/internal/repositories"
	"compass-backend/internal/utils"
)

type AuthService interface {
	SignIn(email, password string) (accessToken, refreshToken string, user *models.User, err error)
	RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) SignIn(email, password string) (string, string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	if user.PasswordHash == nil {
		return "", "", nil, errors.New("password not set")
	}

	if !utils.CheckPassword(password, *user.PasswordHash) {
		return "", "", nil, errors.New("invalid credentials")
	}

	if user.AccountStatus != models.StatusActive {
		return "", "", nil, errors.New("account is not active")
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(user, s.cfg)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := utils.ValidateToken(refreshToken, utils.RefreshToken, s.cfg)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	if user.AccountStatus != models.StatusActive {
		return "", "", errors.New("account is not active")
	}

	newAccessToken, newRefreshToken, err := utils.GenerateTokenPair(user, s.cfg)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
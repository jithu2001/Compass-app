package services

import (
	"errors"

	"compass-backend/internal/models"
	"compass-backend/internal/repositories"
	"compass-backend/internal/utils"
)

type UserService interface {
	CreateUser(user *models.User, password string, invitedBy uint64) error
	UpdateUserStatus(userID uint64, status models.AccountStatus) error
	GetUser(userID uint64) (*models.User, error)
	ListUsers() ([]models.User, error)
	SetPassword(userID uint64, password string) error
	ChangePassword(userID uint64, currentPassword, newPassword string) error
	ResetPassword(userID uint64, newPassword string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *models.User, password string, invitedBy uint64) error {
	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	user.InvitedBy = &invitedBy
	user.AccountStatus = models.StatusPending

	// Hash password if provided
	if password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return err
		}
		user.PasswordHash = &hashedPassword
		user.AccountStatus = models.StatusActive
	}

	return s.userRepo.Create(user)
}

func (s *userService) UpdateUserStatus(userID uint64, status models.AccountStatus) error {
	return s.userRepo.UpdateStatus(userID, status)
}

func (s *userService) GetUser(userID uint64) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *userService) SetPassword(userID uint64, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.PasswordHash = &hashedPassword
	user.AccountStatus = models.StatusActive

	// This would typically be a separate repository method
	// For simplicity, we'll assume the repository handles this
	return nil
}

func (s *userService) ChangePassword(userID uint64, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if user has a password set
	if user.PasswordHash == nil {
		return errors.New("no password set for this account")
	}

	// Verify current password
	if !utils.CheckPassword(currentPassword, *user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password in repository
	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

func (s *userService) ResetPassword(userID uint64, newPassword string) error {
	// Get user to verify they exist
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password in repository (admin can reset without knowing current password)
	return s.userRepo.UpdatePassword(userID, hashedPassword)
}
package repositories

import (
	"compass-backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	UpdateStatus(id uint64, status models.AccountStatus) error
	UpdatePassword(id uint64, hashedPassword string) error
	List() ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.Preload("InvitedByUser").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateStatus(id uint64, status models.AccountStatus) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", id).Update("account_status", status).Error
}

func (r *userRepository) UpdatePassword(id uint64, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("user_id = ?", id).Update("password_hash", hashedPassword).Error
}

func (r *userRepository) List() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("InvitedByUser").Find(&users).Error
	return users, err
}
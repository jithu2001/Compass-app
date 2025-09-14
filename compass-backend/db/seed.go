package db

import (
	"log"

	"compass-backend/config"
	"compass-backend/internal/models"
	"compass-backend/internal/utils"

	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB, cfg *config.Config) error {
	// Check if admin already exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", cfg.Admin.Email).Count(&count)
	
	if count > 0 {
		log.Println("Admin user already exists")
		return nil
	}

	// Create admin user
	hashedPassword, err := utils.HashPassword(cfg.Admin.Password)
	if err != nil {
		return err
	}

	admin := &models.User{
		FullName:      cfg.Admin.Name,
		Email:         cfg.Admin.Email,
		PasswordHash:  &hashedPassword,
		Role:          models.RoleAdmin,
		AccountStatus: models.StatusActive,
	}

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	log.Printf("Admin user created successfully: %s", cfg.Admin.Email)
	return nil
}
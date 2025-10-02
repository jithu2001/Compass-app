package db

import (
	"fmt"
	"log"

	"compass-backend/config"
	"compass-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if cfg.Server.Mode == "debug" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models in correct order
	// First migrate User since it has no dependencies
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("failed to migrate User model: %w", err)
	}
	
	// Then migrate Project which depends on User
	err = DB.AutoMigrate(&models.Project{})
	if err != nil {
		return fmt.Errorf("failed to migrate Project model: %w", err)
	}
	
	// Finally migrate dependent models
	err = DB.AutoMigrate(
		&models.ProjectSpecification{},
		&models.ProjectRFI{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connection established and migrations completed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
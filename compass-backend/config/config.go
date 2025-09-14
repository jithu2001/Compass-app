package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	Admin    AdminConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type ServerConfig struct {
	Port string
	Mode string
}

type AdminConfig struct {
	Email    string
	Password string
	Name     string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "compass_db"),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", "default-secret-change-this"),
			AccessTokenDuration:  parseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m")),
			RefreshTokenDuration: parseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "7d")),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Admin: AdminConfig{
			Email:    getEnv("ADMIN_EMAIL", "admin@compass.com"),
			Password: getEnv("ADMIN_PASSWORD", "AdminPassword123!"),
			Name:     getEnv("ADMIN_NAME", "System Administrator"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	// Handle days suffix
	if strings.HasSuffix(s, "d") {
		days := strings.TrimSuffix(s, "d")
		if d, err := strconv.Atoi(days); err == nil {
			return time.Duration(d) * 24 * time.Hour
		}
	}
	
	duration, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Error parsing duration %s: %v", s, err)
		return 15 * time.Minute
	}
	return duration
}
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

type ServerConfig struct {
	Host string
	Port int
	Mode string
}

type JWTConfig struct {
	Secret     string
	ExpireTime int
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "gin_learning"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "Asia/Ho_Chi_Minh"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_POST", 8080),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "jwt123"),
			ExpireTime: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err != nil {
			return v
		}
	}
	return defaultValue
}

func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode, db.TimeZone)
}

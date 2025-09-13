package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"hari-donganh/gin-learning/internal/infrastructure/config"
)

type DB struct {
	*gorm.DB
}

func NewConnection(cfg *config.Config) *DB {
	database := config.NewDatabase(&cfg.Database)

	if err := autoMigrate(database.DB); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	log.Println("Database connection established successfully")
	return &DB{DB: database.DB}
}

func autoMigrate(db *gorm.DB) error {
	log.Printf("Starting database auto migration ... ")
	err := db.AutoMigrate()
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}
	log.Println("Database auto migration completed successfully")
	return nil
}

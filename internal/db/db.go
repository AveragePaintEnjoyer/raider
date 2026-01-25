package db

import (
	"log"
	"raider/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&models.MainCategory{}, &models.Category{}, &models.Tag{}, &models.Entry{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

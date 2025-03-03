package database

import (
	"live-studio-api/models"
	"log"

	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	
	// Migrasi sesuai dengan Models
	err := db.AutoMigrate(&models.Product{})
	
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
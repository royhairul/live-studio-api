package database

import (
	"live-studio-api/models"
	"log"

	"gorm.io/gorm"
)

// List All Model
var modelsList = []interface{}{
	&models.Product{},
}

func MigrateDatabase(db *gorm.DB) {
	
	// Run AutoMigrate for all models
	if error := db.AutoMigrate(modelsList...); error != nil {
		log.Fatalf("Failed to migrate database: %v", error)
	}

	// Succesfully migration
	log.Println("✅ Database migration completed successfully!")
}
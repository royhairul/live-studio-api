package database

import (
	"log"

	"github.com/royhairul/live-studio-api/models"
	"gorm.io/gorm"

	AccountEntity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	AttendanceEntity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
)

// List All Model
var modelsList = []interface{}{
	&models.Product{},
	&models.User{},
	&models.Studio{},
	&models.Host{},
	&models.Account{},
	&models.ResetPassword{},

	// User Relation model
	&models.UserRelation{},

	// Schedule model
	&models.Schedule{},
	&models.ScheduleShift{},

	// Role and Permission
	&models.Role{},
	&models.Permission{},

	&AccountEntity.Account{},
	&AttendanceEntity.Attendance{},
}

func MigrateDatabase(db *gorm.DB) {

	// Run AutoMigrate for all models
	if error := db.AutoMigrate(modelsList...); error != nil {
		log.Fatalf("Failed to migrate database: %v", error)
	}

	// Succesfully migration
	log.Println("✅ Database migration completed successfully!")
}

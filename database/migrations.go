package database

import (
	"log"

	"gorm.io/gorm"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	attendanceentity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	orderentity "github.com/royhairul/live-studio-api/internal/domains/order/entity"
	permissionentity "github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	productentity "github.com/royhairul/live-studio-api/internal/domains/product/entity"
	roleentity "github.com/royhairul/live-studio-api/internal/domains/role/entity"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	transactionentity "github.com/royhairul/live-studio-api/internal/domains/transaction/entity"

	"github.com/royhairul/live-studio-api/models"
)

// List All Model
var modelsList = []interface{}{
	// &models.Product{},
	&models.User{},
	// &models.Studio{},
	&models.ResetPassword{},

	// // User Relation model
	&models.UserRelation{},

	// // Schedule model
	// &models.ScheduleShift{},

	// // Role and Permission
	// &models.Role{},
	// &models.Permission{},

	&permissionentity.Permission{},
	&roleentity.Role{},
	&hostentity.Host{},
	&accountentity.Account{},
	&studioentity.Studio{},
	&scheduleentity.Schedule{},
	&attendanceentity.Attendance{},

	&productentity.Product{},
	&transactionentity.Transaction{},
	&orderentity.Order{},
}

func MigrateDatabase(db *gorm.DB) {
	// Run AutoMigrate for all models
	if error := db.AutoMigrate(modelsList...); error != nil {
		log.Fatalf("Failed to migrate database: %v", error)
	}

	// Succesfully migration
	log.Println("✅ Database migration completed successfully!")
}

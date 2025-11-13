package database

import (
	"log"

	"gorm.io/gorm"

	// Entities
	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	accountadsentity "github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
	accountsessionentity "github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	attendanceentity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	authentity "github.com/royhairul/live-studio-api/internal/domains/auth/entity"
	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	orderentity "github.com/royhairul/live-studio-api/internal/domains/order/entity"
	permissionentity "github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	productentity "github.com/royhairul/live-studio-api/internal/domains/product/entity"
	roleentity "github.com/royhairul/live-studio-api/internal/domains/role/entity"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	shiftentity "github.com/royhairul/live-studio-api/internal/domains/shift/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	targetentity "github.com/royhairul/live-studio-api/internal/domains/target/entity"
	transactionentity "github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	userentity "github.com/royhairul/live-studio-api/internal/domains/user/entity"
)

// modelsList contains all entities for migration
var modelsList = []interface{}{
	// User & Auth
	&userentity.User{},
	&authentity.ResetPassword{},

	// Role & Permission
	&permissionentity.Permission{},
	&roleentity.Role{},

	// Host
	&hostentity.Host{},

	// Account related
	&accountentity.Account{},
	&accountsessionentity.Accountsession{},
	&accountadsentity.Accountads{},

	// Studio
	&studioentity.Studio{},

	// Attendance & Schedule
	&attendanceentity.Attendance{},
	&scheduleentity.Schedule{},
	&shiftentity.Shift{},

	// Targets
	&targetentity.Target{},

	// Transactions & Products
	&transactionentity.Transaction{},
	&productentity.Product{},
	&orderentity.Order{},
}

// MigrateDatabase runs auto-migration for all models
func MigrateDatabase(db *gorm.DB) {
	if err := db.AutoMigrate(modelsList...); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("✅ Database migration completed successfully!")
}

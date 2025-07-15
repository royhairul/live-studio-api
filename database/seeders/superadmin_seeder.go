package seeders

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/royhairul/live-studio-api/database"
	userentity "github.com/royhairul/live-studio-api/internal/domains/user/entity"
)

func SuperadminSeeder() {
	SUPERADMIN_NAME := os.Getenv("SUPERADMIN_NAME")
	SUPERADMIN_EMAIL := os.Getenv("SUPERADMIN_EMAIL")
	SUPERADMIN_PASS := os.Getenv("SUPERADMIN_PASS")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(SUPERADMIN_PASS), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}

	superadminUser := userentity.User{
		Name:     SUPERADMIN_NAME,
		Email:    SUPERADMIN_EMAIL,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&superadminUser).Error; err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}

	log.Println("✅ Created superadmin user")
}

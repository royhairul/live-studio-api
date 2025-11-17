package database

import (
	"fmt"
	"log"

	"github.com/royhairul/live-studio-api/config"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Buat DSN dari konfigurasi
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	// Koneksi ke database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	log.Println("Connected to database succesfully")

	tenantdb.RegisterTenantCallback(database)
	database.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	DB = database

	return database, nil
}

package main

import (
	"fmt"
	"live-studio-api/config"
	"live-studio-api/database"
	"live-studio-api/routes"
	"live-studio-api/validators"
)

func main() {
	// Inisialisasi validator
	validators.InitValidator()

	// Load konfigurasi
	cfg := config.LoadConfig()

	// Koneksi Database
	database.ConnectDatabase(cfg)

	// Migration
	database.MigrateDatabase(database.DB)

	// Inisialisasi router
	r := routes.SetupRouter()

	// Jalankan server
	r.Run(fmt.Sprintf(":%s", cfg.ServerPort))
}
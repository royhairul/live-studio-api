package main

import (
	"fmt"

	"github.com/royhairul/live-studio-api/config"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/routes"
	"github.com/royhairul/live-studio-api/validators"
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
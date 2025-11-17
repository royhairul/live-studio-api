package main

import (
	"fmt"
	"os"

	"github.com/royhairul/live-studio-api/config"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/database/seeders"
)

func main() {
	cfg := config.LoadConfig()
	database.InitDatabase(cfg)
	database.MigrateDatabase(database.DB)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("⚠️  Harap masukkan nama seeder, contoh: users atau all")
		return
	}

	switch args[1] {
	case "host":
		seeders.HostSeeder()

	case "role":
		seeders.PermissionSeeder()

	case "permission":
		seeders.RoleSeeder()

	case "role-permission":
		seeders.PermissionSeeder()
		seeders.RoleSeeder()

	case "attendance":
		seeders.AttendanceSeeder()

	case "accountsession":
		seeders.AccountSessionSeeder()

	case "accountads":
		seeders.AccountAdsSeeder()

	case "transaction":
		seeders.TransactionSeeder()

	// Handling superadmin
	case "superadmin":
		seeders.SuperadminSeeder()

	default:
		fmt.Println("Seeder tidak dikenal:", args[1])
	}
}

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
	database.ConnectDatabase(cfg)
	database.MigrateDatabase(database.DB)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("⚠️  Harap masukkan nama seeder, contoh: users atau all")
		return
	}

	switch args[1] {
	case "host":
		seeders.HostSeeder()
		break
	case "role":
		seeders.PermissionSeeder()
		break
	case "permission":
		seeders.RoleSeeder()
		break
	case "role-permission":
		seeders.PermissionSeeder()
		seeders.RoleSeeder()
		break

	default:
		fmt.Println("Seeder tidak dikenal:", args[1])
	}
}

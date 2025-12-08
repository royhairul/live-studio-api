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

	case "transaction-month":
		// Usage: go run cmd/seed/main.go transaction-month <tenant_id|all> <year> <month>
		// Example: go run cmd/seed/main.go transaction-month tenant-123 2024 11
		// Example: go run cmd/seed/main.go transaction-month all 2024 11
		seeders.TransactionMonthSeederFromArgs(args[1:])

	case "transaction-range":
		// Usage: go run cmd/seed/main.go transaction-range <tenant_id|all> <start_year> <start_month> <end_year> <end_month>
		// Example: go run cmd/seed/main.go transaction-range tenant-123 2024 1 2024 12
		// Example: go run cmd/seed/main.go transaction-range all 2024 1 2024 12
		seeders.TransactionRangeSeederFromArgs(args[1:])

	// Handling superadmin
	case "superadmin":
		seeders.SuperadminSeeder()

	default:
		fmt.Println("Seeder tidak dikenal:", args[1])
		fmt.Println("")
		fmt.Println("Available seeders:")
		fmt.Println("  - host")
		fmt.Println("  - role")
		fmt.Println("  - permission")
		fmt.Println("  - role-permission")
		fmt.Println("  - attendance")
		fmt.Println("  - accountsession")
		fmt.Println("  - accountads")
		fmt.Println("  - transaction")
		fmt.Println("  - transaction-month <tenant_id|all> <year> <month>")
		fmt.Println("  - transaction-range <tenant_id|all> <start_year> <start_month> <end_year> <end_month>")
		fmt.Println("  - superadmin")
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/royhairul/live-studio-api/database"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("⚠️ Perintah tidak dikenali.")
		fmt.Println("Gunakan:")
		fmt.Println("  migrate  → migrate database")
		fmt.Println("  refresh  → drop semua tabel + migrate ulang")
		return
	}

	command := args[1]

	switch command {
	case "migrate":
		fmt.Println("🚀 Running migration...")
		database.MigrateDatabase(database.DB)

	case "refresh":
		fmt.Println("🔄 Refreshing database (DROP + MIGRATE)...")
		database.RefreshDatabase(database.DB)

	default:
		fmt.Println("⚠️ Perintah tidak dikenali:", command)
		fmt.Println("Gunakan:")
		fmt.Println("  migrate")
		fmt.Println("  refresh")
	}
}

package seeders

import (
	"fmt"

	"github.com/royhairul/live-studio-api/database"
	permissionentity "github.com/royhairul/live-studio-api/internal/domains/permission/entity"
)

func PermissionSeeder() {
	permissions := []permissionentity.Permission{
		// Host
		{Name: "host_view", Group: "Host", Description: "Melihat data host"},
		{Name: "host_create", Group: "Host", Description: "Membuat host baru"},
		{Name: "host_edit", Group: "Host", Description: "Mengedit data host"},
		{Name: "host_delete", Group: "Host", Description: "Menghapus host"},

		// Schedule Host
		{Name: "schedule_host_view", Group: "Schedule Host", Description: "Melihat jadwal host"},
		{Name: "schedule_host_create", Group: "Schedule Host", Description: "Menambahkan jadwal host"},
		{Name: "schedule_host_edit", Group: "Schedule Host", Description: "Mengedit jadwal host"},
		{Name: "schedule_host_delete", Group: "Schedule Host", Description: "Menghapus jadwal host"},

		// Shift
		{Name: "shift_view", Group: "Shift", Description: "Melihat data shift"},
		{Name: "shift_create", Group: "Shift", Description: "Menambahkan shift baru"},
		{Name: "shift_edit", Group: "Shift", Description: "Mengedit data shift"},
		{Name: "shift_delete", Group: "Shift", Description: "Menghapus shift"},

		// Account
		{Name: "account_view", Group: "Account", Description: "Melihat akun pengguna"},
		{Name: "account_create", Group: "Account", Description: "Menambahkan akun pengguna"},
		{Name: "account_edit", Group: "Account", Description: "Mengedit akun pengguna"},
		{Name: "account_delete", Group: "Account", Description: "Menghapus akun pengguna"},

		// Live
		{Name: "live_view", Group: "Live", Description: "Melihat data live"},
		{Name: "live_report", Group: "Live", Description: "Melihat laporan live"},

		// Finance
		{Name: "finance_view", Group: "Finance", Description: "Melihat data keuangan"},
		{Name: "finance_research_view", Group: "Finance", Description: "Melihat hasil riset keuangan"},
		{Name: "finance_research_create", Group: "Finance", Description: "Membuat riset keuangan"},

		// Studio
		{Name: "studio_view", Group: "Studio", Description: "Melihat studio"},
		{Name: "studio_create", Group: "Studio", Description: "Menambahkan studio"},
		{Name: "studio_edit", Group: "Studio", Description: "Mengedit studio"},
		{Name: "studio_delete", Group: "Studio", Description: "Menghapus studio"},

		// User
		{Name: "user_view", Group: "User", Description: "Melihat user"},
		{Name: "user_create", Group: "User", Description: "Menambahkan user"},
		{Name: "user_edit", Group: "User", Description: "Mengedit user"},
		{Name: "user_delete", Group: "User", Description: "Menghapus user"},

		// Role
		{Name: "role_view", Group: "Role", Description: "Melihat role"},
		{Name: "role_create", Group: "Role", Description: "Menambahkan role"},
		{Name: "role_edit", Group: "Role", Description: "Mengedit role"},
		{Name: "role_delete", Group: "Role", Description: "Menghapus role"},
	}

	for _, p := range permissions {
		result := database.DB.Where("name = ?", p.Name).FirstOrCreate(&p)
		if result.Error != nil {
			fmt.Printf("Failed to seed permission '%s': %v\n", p.Name, result.Error)
		} else if result.RowsAffected > 0 {
			fmt.Printf("Permission '%s' created\n", p.Name)
		}
	}
}

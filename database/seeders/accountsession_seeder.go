package seeders

import (
	"fmt"

	"github.com/royhairul/live-studio-api/database"
	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	accountsessionentity "github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	attendanceentity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
)

func AccountSessionSeeder() {
	var accounts []accountentity.Account
	var attendances []attendanceentity.Attendance

	// Ambil data accounts
	if err := database.DB.Find(&accounts).Error; err != nil {
		fmt.Println("❌ Gagal ambil akun:", err)
		return
	}

	// Ambil data attendances
	if err := database.DB.Find(&attendances).Error; err != nil {
		fmt.Println("❌ Gagal ambil attendance:", err)
		return
	}

	if len(accounts) == 0 || len(attendances) == 0 {
		fmt.Println("⚠️ Tidak ada akun atau attendance ditemukan, seeder dilewati")
		return
	}

	var sessions []accountsessionentity.Accountsession

	for _, att := range attendances {
		for idx, acc := range accounts {
			sessions = append(sessions, accountsessionentity.Accountsession{
				AccountID:     acc.ID,
				AttendanceID:  att.ID,
				GMVSalesStart: uint(400000 + idx*1000),
				GMVSalesEnd:   uint(400000 + idx*2500),
				GMVPaidStart:  uint(580000 + idx*1100),
				GMVPaidEnd:    uint(580000 + idx*2080),
				StudioID:      att.StudioID,
			})
		}
	}

	if err := database.DB.Create(&sessions).Error; err != nil {
		fmt.Println("❌ Gagal insert accountsession:", err)
	} else {
		fmt.Printf("✅ AccountSession seeder berhasil, total data: %d\n", len(sessions))
	}
}

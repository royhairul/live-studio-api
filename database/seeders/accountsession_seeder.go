package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/royhairul/live-studio-api/database"
	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	accountsessionentity "github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	attendanceentity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
)

func AccountSessionSeeder() {
	var accounts []accountentity.Account
	var attendances []attendanceentity.Attendance

	// Seed random biar hasilnya beda tiap jalan
	rand.Seed(time.Now().UnixNano())

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
			startSales := uint(400000 + idx*1000)
			startPaid := uint(580000 + idx*1100)

			// selisih random antara 200.000 – 800.000
			diffSales := uint(rand.Intn(600001) + 200000)
			diffPaid := uint(rand.Intn(600001) + 200000)

			sessions = append(sessions, accountsessionentity.Accountsession{
				AccountID:     acc.ID,
				AttendanceID:  att.ID,
				GMVSalesStart: startSales,
				GMVSalesEnd:   startSales + diffSales,
				GMVPaidStart:  startPaid,
				GMVPaidEnd:    startPaid + diffPaid,
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

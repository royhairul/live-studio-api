package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/royhairul/live-studio-api/database"
	attendanceentity "github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	shiftentity "github.com/royhairul/live-studio-api/internal/domains/shift/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
)

func AttendanceSeeder() {
	rand.Seed(time.Now().UnixNano())

	var hosts []hostentity.Host
	var shifts []shiftentity.Shift
	var studios []studioentity.Studio

	// Ambil data dari DB
	if err := database.DB.Find(&hosts).Error; err != nil || len(hosts) == 0 {
		fmt.Println("⚠️ Tidak ada data host di database")
		return
	}
	if err := database.DB.Find(&shifts).Error; err != nil || len(shifts) == 0 {
		fmt.Println("⚠️ Tidak ada data shift di database")
		return
	}
	if err := database.DB.Find(&studios).Error; err != nil || len(studios) == 0 {
		fmt.Println("⚠️ Tidak ada data studio di database")
		return
	}

	// Helper random
	randomHost := func() hostentity.Host {
		return hosts[rand.Intn(len(hosts))]
	}
	randomShift := func() shiftentity.Shift {
		return shifts[rand.Intn(len(shifts))]
	}
	randomStudio := func() studioentity.Studio {
		return studios[rand.Intn(len(studios))]
	}

	now := time.Now()
	var attendances []attendanceentity.Attendance

	// Generate untuk hari ini, kemarin, 2 hari lalu (bisa ditambah sesuai kebutuhan)
	for i := 0; i < 3; i++ {
		// tanggal mundur i hari
		day := now.AddDate(0, 0, -i)

		// buat check-in 2 jam setelah jam 7 pagi, checkout 4 jam setelah checkin
		date := time.Date(day.Year(), day.Month(), day.Day(), 7, 0, 0, 0, day.Location())
		checkIn := date.Add(2 * time.Hour)
		checkOut := checkIn.Add(4 * time.Hour)

		attendances = append(attendances, attendanceentity.Attendance{
			Date:         &date,
			CheckedInAt:  &checkIn,
			CheckedOutAt: &checkOut,
			Status:       "present",
			Note:         "Generated auto seeder",
			HostID:       randomHost().ID,
			ShiftID:      randomShift().ID,
			StudioID:     randomStudio().ID,
		})
	}

	if err := database.DB.Create(&attendances).Error; err != nil {
		fmt.Println("❌ Gagal insert attendance:", err)
	} else {
		fmt.Printf("✅ Attendance seeder berhasil, total data: %d\n", len(attendances))
	}
}

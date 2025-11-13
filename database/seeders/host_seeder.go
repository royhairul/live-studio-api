package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/royhairul/live-studio-api/database"
	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
)

func HostSeeder() {
	var studios []studioentity.Studio

	// Get All Studios
	if err := database.DB.Find(&studios); err != nil {
		fmt.Println("❌ Gagal mengambil data studio:", err)
		return
	}

	if len(studios) == 0 {
		fmt.Println("⚠️  Tidak ada studio yang tersedia. Seeder host dilewati.")
	}

	rand.Seed(time.Now().UnixNano())

	for _, studio := range studios {
		// Buat 1–2 host untuk setiap studio
		hostCount := rand.Intn(2) + 1
		for i := 0; i < hostCount; i++ {
			name := faker.Name()
			phone := faker.Phonenumber()

			host := hostentity.Host{
				Name:     name,
				Phone:    phone,
				StudioID: uint(studio.ID),
			}

			if err := database.DB.Create(&host).Error; err != nil {
				fmt.Printf("❌ Gagal menambahkan host %s untuk Studio %s: %v\n", name, studio.Name, err)
			} else {
				fmt.Printf("✅ Host %s berhasil ditambahkan ke Studio %s\n", name, studio.Name)
			}
		}
	}
}

package seeders

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/database"
	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	accountadsentity "github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
)

func AccountAdsSeeder() {
	var accounts []accountentity.Account

	// Ambil semua akun
	if err := database.DB.Find(&accounts).Error; err != nil {
		fmt.Println("❌ Gagal ambil akun:", err)
		return
	}

	if len(accounts) == 0 {
		fmt.Println("⚠️  Tidak ada akun ditemukan, seeder dilewati")
		return
	}

	now := time.Now()
	var ads []accountadsentity.Accountads

	for _, acc := range accounts {
		// contoh: generate ads untuk hari ini, kemarin, 2 hari yang lalu
		dates := []time.Time{
			now,
			now.AddDate(0, 0, -1),
			now.AddDate(0, 0, -2),
		}

		for i, d := range dates {
			ads = append(ads, accountadsentity.Accountads{
				Spend:     uint(200000 + uint(i*50000)), // contoh variasi spend
				Date:      &d,
				AccountID: acc.ID,
			})
		}
	}

	if err := database.DB.Create(&ads).Error; err != nil {
		fmt.Println("❌ Gagal insert account_ads:", err)
	} else {
		fmt.Printf("✅ AccountAds seeder berhasil, total data: %d\n", len(ads))
	}
}

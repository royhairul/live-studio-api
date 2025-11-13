package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/royhairul/live-studio-api/database"
	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	transactionentity "github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/snowflakeid"
)

// helper untuk pointer
func tptr(t time.Time) *time.Time {
	return &t
}

func TransactionSeeder() {
	now := time.Now()
	rand.Seed(time.Now().UnixNano())

	// 1. Ambil semua akun dari database
	var accounts []accountentity.Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		fmt.Println("❌ Gagal ambil accounts:", err)
		return
	}
	if len(accounts) == 0 {
		fmt.Println("⚠️ Tidak ada account di database, transaksi tidak dibuat")
		return
	}

	// 2. Helper bikin transaksi
	makeTransaction := func(uniqueID string, status string, commission int64, purchaseTime time.Time, accountID uint) transactionentity.Transaction {
		snowflakeid.InitSnowflake()
		return transactionentity.Transaction{
			UniqueID:               uniqueID,
			Status:                 status,
			TotalCommission:        commission,
			TotalCommissionWithMCN: commission,
			PurchaseTime:           tptr(purchaseTime),
			CompleteTime:           tptr(purchaseTime.Add(48 * time.Hour)), // default complete 2 hari setelah purchase
			AccountID:              accountID,
			CreatedAt:              now,
			UpdatedAt:              now,
		}
	}

	// 3. Generate transaksi random
	var transactions []transactionentity.Transaction
	total := 100 // jumlah transaksi (bisa kamu ubah)
	statuses := []string{"Pending", "Waiting for payment", "Completed", "Cancelled"}

	for i := 0; i < total; i++ {
		// random pilih account
		acc := accounts[rand.Intn(len(accounts))]

		// random status
		status := statuses[rand.Intn(len(statuses))]

		// random commission (100rb - 2jt)
		commission := int64(100000 + rand.Intn(2000000))

		// random tanggal (dari 0 sampai 30 hari lalu)
		purchaseTime := now.AddDate(0, 0, -rand.Intn(30))

		transactions = append(transactions,
			makeTransaction(
				fmt.Sprintf("TRX-%d", time.Now().UnixNano()+int64(i)),
				status,
				commission,
				purchaseTime,
				acc.ID,
			))
	}

	// 4. Insert transaksi
	if err := database.DB.Create(&transactions).Error; err != nil {
		fmt.Println("❌ Gagal insert transaction:", err)
	} else {
		fmt.Printf("✅ Transaction auto seeder berhasil, %d transaksi ditambahkan\n", len(transactions))
	}
}

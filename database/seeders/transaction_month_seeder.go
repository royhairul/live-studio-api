package seeders

import (
	"context"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/database"
	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/service"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"

	shopeeclient "github.com/royhairul/live-studio-api/internal/clients/shopee"
	shopeeservice "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	"github.com/royhairul/live-studio-api/internal/pkg/httpclient"
	"github.com/royhairul/live-studio-api/internal/pkg/snowflakeid"

	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	orderservice "github.com/royhairul/live-studio-api/internal/domains/order/service"
	productservice "github.com/royhairul/live-studio-api/internal/domains/product/service"

	accountrepo "github.com/royhairul/live-studio-api/internal/domains/account/repository"
	orderrepo "github.com/royhairul/live-studio-api/internal/domains/order/repository"
	productrepo "github.com/royhairul/live-studio-api/internal/domains/product/repository"
	transactionrepo "github.com/royhairul/live-studio-api/internal/domains/transaction/repository"
)

// TransactionMonthSeederParams contains parameters for seeding transactions for a specific month
type TransactionMonthSeederParams struct {
	TenantID string // Tenant ID (optional - if empty, processes all accounts)
	Year     int    // Year (e.g., 2024)
	Month    int    // Month (1-12)
}

// TransactionMonthSeeder seeds transactions for a specific month
// If TenantID is provided, only processes accounts for that tenant
// If TenantID is empty, processes all accounts (each with their own tenant_id)
// Similar to the cronjob but for seeding historical data per month
func TransactionMonthSeeder(seedParams TransactionMonthSeederParams) {
	// Validate parameters
	if seedParams.Year < 2020 || seedParams.Year > 2100 {
		log.Fatalf("❌ Invalid year: %d", seedParams.Year)
		return
	}
	if seedParams.Month < 1 || seedParams.Month > 12 {
		log.Fatalf("❌ Invalid month: %d (must be 1-12)", seedParams.Month)
		return
	}

	if seedParams.TenantID != "" {
		log.Printf("🌱 Starting transaction seeder for Tenant: %s, Period: %d-%02d\n",
			seedParams.TenantID, seedParams.Year, seedParams.Month)
	} else {
		log.Printf("🌱 Starting transaction seeder for ALL accounts, Period: %d-%02d\n",
			seedParams.Year, seedParams.Month)
	}

	// Create context (with or without tenant filter)
	ctx := context.Background()
	if seedParams.TenantID != "" {
		ctx = tenantdb.AttachTenant(ctx, seedParams.TenantID)
	}

	// Initialize snowflake for ID generation
	if _, err := snowflakeid.InitSnowflake(); err != nil {
		log.Fatalf("❌ Failed to initialize snowflake: %v", err)
		return
	}

	// Initialize repositories
	accountRepo := accountrepo.NewAccountRepository(database.DB)
	transactionRepo := transactionrepo.NewTransactionRepository(database.DB)
	orderRepo := orderrepo.NewOrderRepository(database.DB)
	productRepo := productrepo.NewProductRepository(database.DB)

	// Initialize services
	httpClient := httpclient.NewClient(httpclient.NewHttpClient())
	shopeeAffiliateClient := shopeeclient.NewShopeeAffiliateClient(httpClient)
	shopeeCheckoutSvc := shopeeservice.NewShopeeCheckoutService(shopeeservice.ShopeeCheckoutServiceDeps{
		ShopeeClient: shopeeAffiliateClient,
	})
	accountSvc := accountservice.NewAccountService(accountRepo, nil)
	productSvc := productservice.NewProductService(productRepo)
	orderSvc := orderservice.NewOrderService(orderRepo, productSvc)
	transactionSvc := service.NewTransactionService(
		transactionRepo,
		shopeeCheckoutSvc,
		accountSvc,
		orderSvc,
		productSvc,
	)

	// Get accounts (filtered by tenant if specified, otherwise all accounts)
	var accounts []accountentity.Account
	db := database.DB

	// If no tenant specified, unscoped query to get all accounts regardless of tenant
	if seedParams.TenantID == "" {
		// Use Unscoped to bypass tenant filtering and get all accounts
		db = db.Session(&gorm.Session{})
	}

	if err := db.WithContext(ctx).Find(&accounts).Error; err != nil {
		if seedParams.TenantID != "" {
			log.Fatalf("❌ Failed to fetch accounts for tenant %s: %v", seedParams.TenantID, err)
		} else {
			log.Fatalf("❌ Failed to fetch accounts: %v", err)
		}
		return
	}

	if len(accounts) == 0 {
		if seedParams.TenantID != "" {
			log.Printf("⚠️  No accounts found for tenant %s, skipping seeding", seedParams.TenantID)
		} else {
			log.Printf("⚠️  No accounts found, skipping seeding")
		}
		return
	}

	if seedParams.TenantID != "" {
		log.Printf("📊 Found %d account(s) for tenant %s", len(accounts), seedParams.TenantID)
	} else {
		log.Printf("📊 Found %d account(s) across all tenants", len(accounts))
	}

	// Seed transactions day by day for the entire month
	startDate := time.Date(seedParams.Year, time.Month(seedParams.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // First day of next month

	totalCreated := 0
	totalDays := 0

	// Iterate through each day of the month
	for currentDate := startDate; currentDate.Before(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		dateStr := currentDate.Format("2006-01-02")
		totalDays++

		log.Printf("📅 Processing date: %s", dateStr)

		dayTotal := 0

		// Process each account individually with its own tenant context
		for _, account := range accounts {
			// Create context with the account's tenant_id
			accountCtx := context.Background()
			if account.TenantID != "" {
				accountCtx = tenantdb.AttachTenant(accountCtx, account.TenantID)
			}

			req := params.CreateTransactionRequest{
				Date: dateStr,
			}

			// Create transaction service with account-specific context
			result, err := transactionSvc.Create(accountCtx, req)
			if err != nil {
				log.Printf("   ❌ Error processing account '%s' (Tenant: %s): %v",
					account.Name, account.TenantID, err)
				continue
			}

			// Find result for this specific account
			for _, r := range result {
				if r.AccountID == account.ID {
					if r.NewTransaction > 0 {
						dayTotal += r.NewTransaction
						log.Printf("   ✅ Account '%s' (ID: %d, Tenant: %s): %d new transaction(s)",
							r.AccountName, r.AccountID, account.TenantID, r.NewTransaction)
					}
					break
				}
			}
		}

		totalCreated += dayTotal
		if dayTotal > 0 {
			log.Printf("   📈 Day total: %d transaction(s)", dayTotal)
		} else {
			log.Printf("   ℹ️  No new transactions for this day")
		}
	}

	log.Printf("\n✅ Transaction seeding completed!")
	if seedParams.TenantID != "" {
		log.Printf("   Tenant ID: %s", seedParams.TenantID)
	} else {
		log.Printf("   Processed all tenants")
	}
	log.Printf("   Period: %d-%02d", seedParams.Year, seedParams.Month)
	log.Printf("   Days processed: %d", totalDays)
	log.Printf("   Total transactions created: %d", totalCreated)
}

// TransactionMonthSeederFromArgs is a helper function to parse command-line arguments
// Usage:
//   - With tenant: go run cmd/seed/main.go transaction-month <tenant_id> <year> <month>
//   - All accounts: go run cmd/seed/main.go transaction-month all <year> <month>
//
// Examples:
//   - go run cmd/seed/main.go transaction-month tenant-123 2024 11
//   - go run cmd/seed/main.go transaction-month all 2024 11
func TransactionMonthSeederFromArgs(args []string) {
	if len(args) < 4 {
		log.Println("❌ Usage: transaction-month <tenant_id|all> <year> <month>")
		log.Println("   Example with tenant: transaction-month tenant-123 2024 11")
		log.Println("   Example all accounts: transaction-month all 2024 11")
		return
	}

	tenantID := args[1]
	// If tenant_id is "all", process all accounts
	if tenantID == "all" {
		tenantID = ""
	}

	year, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("❌ Invalid year: %s", args[2])
	}

	month, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatalf("❌ Invalid month: %s", args[3])
	}

	TransactionMonthSeeder(TransactionMonthSeederParams{
		TenantID: tenantID,
		Year:     year,
		Month:    month,
	})
}

// TransactionRangeSeeder seeds transactions for a date range
// This is useful for seeding multiple months at once
func TransactionRangeSeeder(tenantID string, startYear, startMonth, endYear, endMonth int) {
	if tenantID != "" {
		log.Printf("🌱 Starting transaction range seeder")
		log.Printf("   Tenant: %s", tenantID)
	} else {
		log.Printf("🌱 Starting transaction range seeder for ALL accounts")
	}
	log.Printf("   From: %d-%02d", startYear, startMonth)
	log.Printf("   To: %d-%02d", endYear, endMonth)

	currentYear := startYear
	currentMonth := startMonth

	for {
		TransactionMonthSeeder(TransactionMonthSeederParams{
			TenantID: tenantID,
			Year:     currentYear,
			Month:    currentMonth,
		})

		// Move to next month
		currentMonth++
		if currentMonth > 12 {
			currentMonth = 1
			currentYear++
		}

		// Check if we've reached the end
		if currentYear > endYear || (currentYear == endYear && currentMonth > endMonth) {
			break
		}
	}

	log.Printf("\n✅ Transaction range seeding completed!")
}

// TransactionRangeSeederFromArgs is a helper function to parse command-line arguments for range seeding
// Usage:
//   - With tenant: go run cmd/seed/main.go transaction-range <tenant_id> <start_year> <start_month> <end_year> <end_month>
//   - All accounts: go run cmd/seed/main.go transaction-range all <start_year> <start_month> <end_year> <end_month>
//
// Examples:
//   - go run cmd/seed/main.go transaction-range tenant-123 2024 1 2024 12
//   - go run cmd/seed/main.go transaction-range all 2024 1 2024 12
func TransactionRangeSeederFromArgs(args []string) {
	if len(args) < 6 {
		log.Println("❌ Usage: transaction-range <tenant_id|all> <start_year> <start_month> <end_year> <end_month>")
		log.Println("   Example with tenant: transaction-range tenant-123 2024 1 2024 12")
		log.Println("   Example all accounts: transaction-range all 2024 1 2024 12")
		return
	}

	tenantID := args[1]
	// If tenant_id is "all", process all accounts
	if tenantID == "all" {
		tenantID = ""
	}

	startYear, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("❌ Invalid start year: %s", args[2])
	}

	startMonth, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatalf("❌ Invalid start month: %s", args[3])
	}

	endYear, err := strconv.Atoi(args[4])
	if err != nil {
		log.Fatalf("❌ Invalid end year: %s", args[4])
	}

	endMonth, err := strconv.Atoi(args[5])
	if err != nil {
		log.Fatalf("❌ Invalid end month: %s", args[5])
	}

	TransactionRangeSeeder(tenantID, startYear, startMonth, endYear, endMonth)
}

package service

import (
	"context"
	"log"
	"sync"
	"time"

	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	accountparams "github.com/royhairul/live-studio-api/internal/domains/account/params"
	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	"github.com/royhairul/live-studio-api/internal/domains/finance/params"
)

type FinanceServiceImpl struct {
	shopeeSvc  ShopeeService.ShopeeFinanceService
	accountSvc accountservice.AccountService
	options    params.FinanceFilter
}

// --- Builder Pattern (WithX) ---

func (f *FinanceServiceImpl) WithAccountUniqueID(uniqueID string) FinanceService {
	instance := *f
	instance.options.UniqueID = &uniqueID
	return &instance
}

func (f *FinanceServiceImpl) WithStudioID(studioID string) FinanceService {
	instance := *f
	instance.options.StudioID = &studioID
	return &instance
}

func (f *FinanceServiceImpl) WithStatus(status string) FinanceService {
	instance := *f
	instance.options.Status = &status
	return &instance
}

func (f *FinanceServiceImpl) WithPaymentMethod(method string) FinanceService {
	instance := *f
	instance.options.PaymentMethod = &method
	return &instance
}

// --- Constructor ---

func NewFinanceService(
	shopeeSvc ShopeeService.ShopeeFinanceService,
	accountSvc accountservice.AccountService,
) FinanceService {
	return &FinanceServiceImpl{
		shopeeSvc:  shopeeSvc,
		accountSvc: accountSvc,
		options:    params.FinanceFilter{},
	}
}

// --- Core Logic ---

func (f *FinanceServiceImpl) FindAll(ctx context.Context, startDate *time.Time, endDate *time.Time) (*params.FinanceResponse, error) {
	accountQuery := f.accountSvc

	// Apply optional filters for account
	if f.options.UniqueID != nil {
		accountQuery = accountQuery.WithUniqueID(*f.options.UniqueID)
	}
	if f.options.StudioID != nil {
		accountQuery = accountQuery.WithStudioID(*f.options.StudioID)
	}

	// Get all filtered accounts
	accounts, err := accountQuery.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		sem   = make(chan struct{}, 5)
		items []params.FinanceItem
	)

	for _, acc := range accounts {
		wg.Add(1)
		sem <- struct{}{}

		go func(account *accountparams.AccountResponse) {
			defer wg.Done()
			defer func() { <-sem }()

			commissions, err := f.shopeeSvc.GetPaymentCommission(account.Cookie, startDate, endDate)
			if err != nil {
				log.Printf("⚠️ gagal mengambil komisi untuk %s: %v", account.Name, err)
				return
			}

			var localItems []params.FinanceItem
			for _, comm := range commissions.List {
				item := params.NewFinanceItem(*account, comm)

				// --- Internal filtering (status + method) ---
				if f.options.Status != nil && item.PaymentStatus != *f.options.Status {
					continue
				}
				if f.options.PaymentMethod != nil && item.PaymentMethod != *f.options.PaymentMethod {
					continue
				}

				localItems = append(localItems, *item)
			}

			// Thread-safe append
			mu.Lock()
			items = append(items, localItems...)
			mu.Unlock()
		}(acc)
	}

	wg.Wait()

	return params.NewFinanceResponse(items), nil
}

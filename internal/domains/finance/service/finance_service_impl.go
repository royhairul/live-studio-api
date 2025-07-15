package service

import (
	"log"

	ShopeeParams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	"github.com/royhairul/live-studio-api/internal/domains/account/repository"
	"github.com/royhairul/live-studio-api/internal/domains/finance/params"
)

type FinanceServiceImpl struct {
	shopeeSvc   ShopeeService.ShopeeFinanceService
	accountRepo repository.AccountRepository
}

func NewFinanceService(shopeeSvc ShopeeService.ShopeeFinanceService, accountRepo repository.AccountRepository) FinanceService {
	return &FinanceServiceImpl{shopeeSvc, accountRepo}
}

func (f *FinanceServiceImpl) FindAll(financeReq ShopeeParams.ShopeeLiveFinanceRequest) ([]*params.FinanceResponse, error) {
	accounts, err := f.accountRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*params.FinanceResponse
	for _, account := range accounts {
		report, err := f.shopeeSvc.GetShopeeLiveSalesReportRange(financeReq, account.Cookie)
		if err != nil {
			log.Println("failed to get report for account %s: %v\n", account.Name, err)
			continue
		}

		result = append(result, &params.FinanceResponse{
			AccountName: account.Name,
			Total:       len(report),
			ReportLive:  report,
		})
	}

	return result, nil
}

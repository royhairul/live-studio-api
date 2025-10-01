package service

import (
	"fmt"
	"log"

	"github.com/royhairul/live-studio-api/internal/domains/finance/params"

	ShopeeParams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"
)

type FinanceServiceImpl struct {
	shopeeSvc         ShopeeService.ShopeeFinanceService
	studioSvc         studioservice.StudioService
	accountSvc        accountservice.AccountService
	accountSessionSvc accountsessionservice.AccountsessionService
}

func NewFinanceService(
	shopeeSvc ShopeeService.ShopeeFinanceService,
	studioSvc studioservice.StudioService,
	accountSvc accountservice.AccountService,
	accountSessionSvc accountsessionservice.AccountsessionService,
) FinanceService {
	return &FinanceServiceImpl{shopeeSvc, studioSvc, accountSvc, accountSessionSvc}
}

func (f *FinanceServiceImpl) FindAll(financeReq ShopeeParams.ShopeeLiveFinanceRequest) ([]*params.FinanceResponse, error) {
	accounts, err := f.accountSvc.FindAll()
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

// FindAllCommission implements FinanceService.
func (f *FinanceServiceImpl) FindAllCommission() (*params.CommissionTotalResponse, error) {
	studios, err := f.studioSvc.FindAll()
	if err != nil {
		return nil, err
	}

	var results params.CommissionTotalResponse
	for _, studio := range studios {
		accounts, err := f.accountSvc.WithStudioID(fmt.Sprintf("%d", studio.ID)).FindAll()
		if err != nil {
			return nil, err
		}

		var list params.CommissionStudioResponse
		list.StudioName = studio.Name
		list.TotalGMV = 0
		list.TotalCommission = 0
		list.TotalIncome = 0

		for _, acc := range accounts {
			sessions, err := f.accountSessionSvc.WithAccountID(fmt.Sprintf("%d", acc.ID)).FindAll()
			if err != nil {
				continue
			}
			for _, session := range sessions {
				list.TotalGMV += session.GMVSales
			}
		}
	}

	return &results, nil
}

// FindByStudioCommission implements FinanceService.
func (f *FinanceServiceImpl) FindByStudioCommission() (*params.CommissionStudioDetailResponse, error) {
	panic("unimplemented")
}

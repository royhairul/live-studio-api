package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"
	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/repository"
)

type ShopeeFinanceServiceImpl struct {
	repo repository.FinanceRepository
}

func NewShopeeFinanceService(repo repository.FinanceRepository) ShopeeFinanceService {
	return &ShopeeFinanceServiceImpl{repo: repo}
}

func (s *ShopeeFinanceServiceImpl) GetShopeeLiveSalesReport(financeReq params.ShopeeLiveFinanceRequest, cookie string) (*params.ShopeeLiveFinanceResponse, error) {
	return s.repo.GetShopeeLiveSalesReport(financeReq, cookie)
}

func (s *ShopeeFinanceServiceImpl) GetShopeeLiveSalesReportRange(financeReq params.ShopeeLiveFinanceRequest, cookie string) ([]params.ShopeeLiveReportItem, error) {
	return s.repo.GetShopeeLiveSalesReportRange(financeReq, cookie)
}

// GetShopeeLiveRealTime implements ShopeeFinanceService.
func (s *ShopeeFinanceServiceImpl) GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItem, error) {
	return s.repo.GetShopeeLiveRealTime(cookie)
}

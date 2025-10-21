package service

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"go.uber.org/fx"
)

type ShopeeFinanceServiceDeps struct {
	fx.In
	ShopeeClient *shopee.ShopeeClient `name:"affiliateShopeeClient"`
}

type ShopeeFinanceServiceImpl struct {
	ShopeeClient *shopee.ShopeeClient
}

func NewShopeeFinanceService(deps ShopeeFinanceServiceDeps) ShopeeFinanceService {
	return &ShopeeFinanceServiceImpl{ShopeeClient: deps.ShopeeClient}
}

// GetPaymentCommission implements ShopeeFinanceService.
func (s *ShopeeFinanceServiceImpl) GetPaymentCommission(cookie string, startDate *time.Time, endDate *time.Time) (*params.ShopeeFinanceCommissionReportResponse, error) {
	endpoint := "/api/v3/payment/billing_list"
	query := map[string]string{
		"order_completed_start_time": fmt.Sprint(startDate.Unix()),
		"order_completed_end_time":   fmt.Sprint(endDate.Unix()),
		"version":                    "1",
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeFinanceCommissionReportResponse]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return &result.Data, nil
}

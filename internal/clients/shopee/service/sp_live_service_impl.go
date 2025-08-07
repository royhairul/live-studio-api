package service

import (
	"fmt"

	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"go.uber.org/fx"
)

type ShopeeLiveServiceDeps struct {
	fx.In
	ShopeeClient *shopee.ShopeeClient `name:"creatorShopeeClient"`
}

type ShopeeLiveServiceImpl struct {
	ShopeeClient *shopee.ShopeeClient
}

func NewShopeeLiveService(deps ShopeeLiveServiceDeps) ShopeeLiveService {
	return &ShopeeLiveServiceImpl{ShopeeClient: deps.ShopeeClient}
}

// GetShopeeLiveRealTime implements ShopeeLiveService.
func (s ShopeeLiveServiceImpl) GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItemRT, error) {
	endpoint := "/supply/api/lm/sellercenter/realtime/sessionList"
	query := map[string]string{
		"page":     "1",
		"pageSize": "10",
		"name":     "",
		"orderBy":  "",
		"sort":     "desc",
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeLiveRealTimeResponse]

	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data.List, nil
}

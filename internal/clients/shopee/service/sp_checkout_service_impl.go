package service

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"go.uber.org/fx"
)

type ShopeeCheckoutServiceDeps struct {
	fx.In
	ShopeeClient *shopee.ShopeeClient `name:"affiliateShopeeClient"`
}

type ShopeeCheckoutServiceImpl struct {
	ShopeeClient *shopee.ShopeeClient
}

func NewShopeeCheckoutService(deps ShopeeCheckoutServiceDeps) ShopeeCheckoutService {
	return &ShopeeCheckoutServiceImpl{ShopeeClient: deps.ShopeeClient}
}

// GetAll implements ShopeeCheckoutService.
func (s *ShopeeCheckoutServiceImpl) GetAll(cookie string) {
	panic("unimplemented")
}

// GetByRangeTime implements ShopeeCheckoutService.
func (s *ShopeeCheckoutServiceImpl) GetByRangeTime(cookie string, startTime time.Time, endTime time.Time) (*params.ShopeeCheckoutResponse, error) {
	// Initialize the request to the Shopee API
	pageNum := 1
	pageSize := 500

	var allData []*params.ShopeeCheckoutList

	for {
		endpoint := "/api/v3/report/list"
		query := map[string]string{
			"page_size":       fmt.Sprintf("%d", pageSize),
			"page_num":        fmt.Sprintf("%d", pageNum),
			"purchase_time_s": fmt.Sprintf("%d", startTime.Unix()),
			"purchase_time_e": fmt.Sprintf("%d", endTime.Unix()),
			"version":         "1",
		}

		req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
		if err != nil {
			return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
		}

		var result params.ShopeeApiResponse[params.ShopeeCheckoutResponse]
		if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
			return nil, fmt.Errorf("failed to execute request to %s: %w", endpoint, err)
		}

		if result.Error != 0 {
			return nil, fmt.Errorf(result.ErrorMsg)
		}

		for i := range result.Data.List {
			allData = append(allData, result.Data.List[i])
		}

		total := result.Data.TotalCount
		if pageSize*pageNum > total {
			break
		}

		pageNum++
	}

	response := params.ShopeeCheckoutResponse{
		TotalCount: len(allData),
		List:       allData,
	}

	return &response, nil
}

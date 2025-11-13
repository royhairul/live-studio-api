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

// GetLiveSessionRT implements ShopeeLiveService.
func (s ShopeeLiveServiceImpl) GetLiveSessionRT(cookie string) ([]params.ShopeeLiveReportItemRT, error) {
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

	var result params.ShopeeApiResponse[params.ShopeeApiPaginationResult[params.ShopeeLiveReportItemRT]]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data.List, nil
}

// GetDashboardOverviewRT implements ShopeeLiveService.
func (s *ShopeeLiveServiceImpl) GetDashboardOverviewRT(cookie string, sessionID string) (params.ShopeeLiveOverviewResponse, error) {
	endpoint := "/supply/api/lm/sellercenter/realtime/dashboard/overview"
	query := map[string]string{
		"sessionId": sessionID,
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return params.ShopeeLiveOverviewResponse{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeLiveOverviewResponse]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return params.ShopeeLiveOverviewResponse{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return params.ShopeeLiveOverviewResponse{}, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data, nil
}

// GetDashboardProductListRT implements ShopeeLiveService.
func (s *ShopeeLiveServiceImpl) GetDashboardProductListRT(cookie string, sessionID string, page int, pageSize int) (params.ShopeeApiPaginationResult[params.ShopeeLiveProductResponse], error) {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		page = 10
	}

	endpoint := "/supply/api/lm/sellercenter/realtime/dashboard/productList"
	query := map[string]string{
		"sessionId":            sessionID,
		"productName":          "",
		"productListTimeRange": "0",
		"productListOrderBy":   "productClicks",
		"sort":                 "desc",
		"page":                 fmt.Sprint(page),
		"pageSize":             fmt.Sprint(pageSize),
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return params.ShopeeApiPaginationResult[params.ShopeeLiveProductResponse]{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeApiPaginationResult[params.ShopeeLiveProductResponse]]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return params.ShopeeApiPaginationResult[params.ShopeeLiveProductResponse]{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return params.ShopeeApiPaginationResult[params.ShopeeLiveProductResponse]{}, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data, nil
}

// GetDashboardBuyerRT implements ShopeeLiveService.
func (s *ShopeeLiveServiceImpl) GetDashboardBuyerRT(cookie string, sessionID string) ([]params.ShopeeLiveAudienceAnalyticsResponse, error) {
	endpoint := "/supply/api/lm/sellercenter/realtime/dashboard/buyer-profile"
	query := map[string]string{
		"sessionId": sessionID,
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return []params.ShopeeLiveAudienceAnalyticsResponse{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[[]params.ShopeeLiveAudienceAnalyticsResponse]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return []params.ShopeeLiveAudienceAnalyticsResponse{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return []params.ShopeeLiveAudienceAnalyticsResponse{}, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data, nil
}

// GetDashboardViewerRT implements ShopeeLiveService.
func (s *ShopeeLiveServiceImpl) GetDashboardViewerRT(cookie string, sessionID string) ([]params.ShopeeLiveAudienceAnalyticsResponse, error) {
	endpoint := "/supply/api/lm/sellercenter/realtime/dashboard/viewer-profile"
	query := map[string]string{
		"sessionId": sessionID,
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[[]params.ShopeeLiveAudienceAnalyticsResponse]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data, nil
}

// GetDashboardViewerSourceRT implements ShopeeLiveService.
func (s *ShopeeLiveServiceImpl) GetDashboardViewerSourceRT(cookie string, sessionID string) (params.ShopeeLiveAudienceAnalyticsResponse, error) {
	endpoint := "/supply/api/lm/sellercenter/realtime/dashboard/viewer-source"
	query := map[string]string{
		"sessionId": sessionID,
	}

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return params.ShopeeLiveAudienceAnalyticsResponse{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeLiveAudienceAnalyticsResponse]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return params.ShopeeLiveAudienceAnalyticsResponse{}, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return params.ShopeeLiveAudienceAnalyticsResponse{}, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data, nil
}

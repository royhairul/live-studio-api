package service

import "github.com/royhairul/live-studio-api/internal/clients/shopee/params"

type ShopeeLiveService interface {
	GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItemRT, error)

	GetDashboardOverviewRT(cookie, sessionID string) (params.ShopeeLiveOverviewResponse, error)
	GetDashboardBuyerRT(cookie, sessionID string) ([]params.ShopeeLiveAudienceAnalyticsResponse, error)
	GetDashboardViewerRT(cookie string, sessionID string) ([]params.ShopeeLiveAudienceAnalyticsResponse, error)
	GetDashboardProductListRT(cookie, sessionID string) ([]params.ShopeeLiveProductResponse, error)
}

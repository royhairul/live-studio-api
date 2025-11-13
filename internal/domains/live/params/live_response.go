package params

import (
	shopeeparam "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
)

type LiveResponse struct {
	AccountID   string                               `json:"id"`
	AccountName string                               `json:"name"`
	Relive      int                                  `json:"relive"`
	Total       int                                  `json:"total"`
	ReportLive  []shopeeparam.ShopeeLiveReportItemRT `json:"reportLive"`
}

type LiveDetailResponse struct {
	AccountID     string                                                                       `json:"id"`
	AccountName   string                                                                       `json:"name"`
	Overview      shopeeparam.ShopeeLiveOverviewResponse                                       `json:"overview"`
	BuyerProfile  []shopeeparam.ShopeeLiveAudienceAnalyticsResponse                            `json:"buyer_profile"`
	ViewerProfile []shopeeparam.ShopeeLiveAudienceAnalyticsResponse                            `json:"viewer_profile"`
	ViewerSource  shopeeparam.ShopeeLiveAudienceAnalyticsResponse                              `json:"viewer_source"`
	Products      shopeeparam.ShopeeApiPaginationResult[shopeeparam.ShopeeLiveProductResponse] `json:"products"`
}

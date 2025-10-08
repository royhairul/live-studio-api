package params

import (
	shopeeparam "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
)

type LiveResponse struct {
	AccountName string `json:"name"`
	Relive      int    `json:"relive"`
	Total       int    `json:"total"`
	ReportLive  any    `json:"reportLive"`
}

type LiveDetailResponse struct {
	AccountID     string                                                                       `json:"account_id"`
	AccountName   string                                                                       `json:"account_name"`
	Overview      shopeeparam.ShopeeLiveOverviewResponse                                       `json:"overview"`
	BuyerProfile  []shopeeparam.ShopeeLiveAudienceAnalyticsResponse                            `json:"buyer_profile"`
	ViewerProfile []shopeeparam.ShopeeLiveAudienceAnalyticsResponse                            `json:"viewer_profile"`
	Products      shopeeparam.ShopeeApiPaginationResult[shopeeparam.ShopeeLiveProductResponse] `json:"products"`
}

package params

// Shopee Live Report
type ShopeeLiveReportItemRT struct {
	SessionID         int64   `json:"sessionId"`
	Title             string  `json:"title"`
	StartTime         int64   `json:"startTime"`
	Duration          int64   `json:"duration"`
	Views             int     `json:"viewers"`
	Comments          int     `json:"comments"`
	Atc               int     `json:"atc"`
	EngagedUV         int     `json:"engagedUv"`
	AvgEngagedCCU     int     `json:"avgEngagedCcu"`
	PlacedOrders      int     `json:"placedOrders"`
	PlacedItemSold    int     `json:"placedItemSold"`
	PlacedSales       float64 `json:"placedSales"`
	ConfirmedOrders   int     `json:"confirmedOrders"`
	ConfirmedItemSold int     `json:"confirmedItemSold"`
	ConfirmedSales    float64 `json:"confirmedSales"`

	OmsetPerHour float64 `json:"omsetPerHour"`
}

// Shopee Live Overview
type ShopeeLiveOverviewResponse struct {
	Status                   int                      `json:"status"`
	CCU                      int                      `json:"ccu"`
	Viewers                  int                      `json:"viewers"`
	PCU                      int                      `json:"pcu"`
	Views                    int                      `json:"views"`
	Buyers                   int                      `json:"buyers"`
	ConfirmedBuyers          int                      `json:"confirmedBuyers"`
	AvgViewTime              float64                  `json:"avgViewTime"`
	Atc                      int                      `json:"atc"`
	PlacedItemsSold          int                      `json:"placedItemsSold"`
	PlacedOrder              int                      `json:"placedOrder"`
	PlacedGMV                float64                  `json:"placedGmv"`
	ConfirmedItemsSold       int                      `json:"confirmedItemsSold"`
	ConfirmedOrder           int                      `json:"confirmedOrder"`
	ConfirmedGMV             float64                  `json:"confirmedGmv"`
	CTR                      float64                  `json:"ctr"`
	CO                       float64                  `json:"co"`
	ConfirmedCO              float64                  `json:"confirmedCo"`
	EngagedCCU               int                      `json:"engagedCcu"`
	CommentsWithinLastMinute int                      `json:"commentsWhithinLastOneMinute"`
	AtcWithinLastMinute      int                      `json:"atcWhithinLastOneMinute"`
	GPM                      float64                  `json:"gpm"`
	ConfirmedGPM             float64                  `json:"confirmedGpm"`
	ABS                      float64                  `json:"abs"`
	ConfirmedABS             float64                  `json:"confirmedAbs"`
	CommentsRate             float64                  `json:"commentsRate"`
	AvgEngagedCCU            int                      `json:"avgEngagedCcu"`
	EngagedViewers           int                      `json:"engagedViewers"`
	EngagementData           ShopeeLiveEngagementData `json:"engagementData"`
}

type ShopeeLiveEngagementData struct {
	Likes              int     `json:"likes"`
	Comments           int     `json:"comments"`
	Viewers            int     `json:"viewers"`
	Shares             int     `json:"shares"`
	Views              int     `json:"views"`
	NewFollowers       int     `json:"newFollowers"`
	AvgViewingDuration float64 `json:"avgViewingDuration"`
}

// Shopee Live Product List
type ShopeeLiveProductResponse struct {
	ItemID            int64   `json:"itemId"`
	Title             string  `json:"title"`
	CoverImage        string  `json:"coverImage"`
	MaxPrice          float64 `json:"maxPrice"`
	MinPrice          float64 `json:"minPrice"`
	ProductClicks     int     `json:"productClicks"`
	Ctr               float64 `json:"ctr"`
	Atc               int     `json:"atc"`
	OrdersCreated     int     `json:"ordersCreated"`
	Revenue           float64 `json:"revenue"`
	ItemSold          int     `json:"itemSold"`
	Cor               float64 `json:"cor"`
	ConfirmedOrderCnt int     `json:"confirmedOrderCnt"`
	ConfirmedRevenue  float64 `json:"confirmedRevenue"`
	ConfirmedItemSold int     `json:"confirmedItemSold"`
	ConfirmedCor      float64 `json:"confirmedCor"`
}

// Shopee Live Buyer and Viewer
type ShopeeLiveAudienceAnalyticsResponse struct {
	Type         string                            `json:"type"`
	Distribution []ShopeeLiveAudienceAnalyticsItem `json:"distribution"`
}

type ShopeeLiveAudienceAnalyticsItem struct {
	ItemName   string  `json:"itemName"`
	Value      int     `json:"value"`
	Percentage float64 `json:"percentage,omitempty"`
}

func (r *ShopeeLiveAudienceAnalyticsResponse) CalculatePercentage() {
	var total int
	for _, item := range r.Distribution {
		total += item.Value
	}

	if total == 0 {
		return
	}

	for i := range r.Distribution {
		r.Distribution[i].Percentage = (float64(r.Distribution[i].Value) / float64(total)) * 100
	}
}

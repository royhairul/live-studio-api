package params

type ShopeeLiveFinanceResponse struct {
	Total int                    `json:"total"`
	List  []ShopeeLiveReportItem `json:"list"`
	// Page      int                    `json:"page"`
	// PageSize  int                    `json:"pageSize"`
	// TotalPage int                    `json:"totalPage"`
}

type ShopeeLiveReportItem struct {
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
}

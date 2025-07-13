package params

type ShopeeLiveRealTimeResponse struct {
	Total int                      `json:"total"`
	List  []ShopeeLiveReportItemRT `json:"list"`
}
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

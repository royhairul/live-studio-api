package params

type ShopeeCheckoutRequest struct {
	Cookie    string `json:"cookie"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

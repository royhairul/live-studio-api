package params

// Performa Host
type PerformaHostResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	TotalDuration int64  `json:"total_duration"`
	TotalSales    uint   `json:"total_sales"`
	TotalPaid     uint   `json:"total_paid"`
}

type PerformaHostDetailResponse struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	TotalDuration int64                      `json:"total_duration"`
	TotalSales    uint                       `json:"total_sales"`
	TotalPaid     uint                       `json:"total_paid"`
	AvgSales      uint                       `json:"avg_sales"`
	AvgPaid       uint                       `json:"avg_paid"`
	List          []PerformaHostItemResponse `json:"list"`
}

type PerformaHostItemResponse struct {
	AccountName string `json:"account_name"`
	Duration    int64  `json:"duration"`
	Sales       uint   `json:"sales"`
	Paid        uint   `json:"paid"`
}

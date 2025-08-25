package params

// Performa Studio
type PerformaStudioResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	TotalDuration int64  `json:"total_duration"`
	TotalSales    uint   `json:"total_sales"`
	TotalPaid     uint   `json:"total_paid"`
}

type PerformaStudioDetailResponse struct {
	StudioName string                       `json:"studio_name"`
	Duration   int64                        `json:"duration"`
	Sales      uint                         `json:"sales"`
	Paid       uint                         `json:"paid"`
	AvgSales   uint                         `json:"avg_sales"`
	AvgPaid    uint                         `json:"avg_paid"`
	List       []PerformaStudioItemResponse `json:"list"`
}

type PerformaStudioItemResponse struct {
	AccountName string `json:"account_name"`
	Duration    int64  `json:"duration"`
	Sales       uint   `json:"sales"`
	Paid        uint   `json:"paid"`
}

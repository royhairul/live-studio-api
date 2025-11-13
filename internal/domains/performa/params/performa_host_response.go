package params

// === Common Struct ===
type TotalPerformaHost struct {
	Duration int64 `json:"duration"`
	Sales    int64 `json:"sales"`
	Paid     int64 `json:"paid"`
}

// === Summary Response (untuk list all hosts) ===
type PerformaHostSummaryResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	TotalDuration int64  `json:"total_duration"`
	TotalSales    int64  `json:"total_sales"`
	TotalPaid     int64  `json:"total_paid"`
}

// === Detail Response (untuk host by ID) ===
type PerformaHostDetailResponse struct {
	PerformaHostSummaryResponse
	AvgSales int64                      `json:"avg_sales,omitempty"`
	AvgPaid  int64                      `json:"avg_paid,omitempty"`
	Total    *TotalPerformaHost         `json:"total,omitempty"`
	List     []PerformaHostItemResponse `json:"list,omitempty"`
}

// === Item Response (per akun atau sesi) ===
type PerformaHostItemResponse struct {
	AccountName string `json:"account_name"`
	Duration    int64  `json:"duration"`
	Sales       int64  `json:"sales"`
	Paid        int64  `json:"paid"`
}

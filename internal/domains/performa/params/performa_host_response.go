package params

import "time"

// Performa Host
type PerformaHostResponse struct {
	ID            string                       `json:"id"`
	Name          string                       `json:"name"`
	Date          *time.Time                   `json:"date"`
	TotalDuration int64                        `json:"total_duration"`
	TotalSales    uint                         `json:"total_sales"`
	TotalPaid     uint                         `json:"total_paid"`
	AvgSales      uint                         `json:"avg_sales"`
	AvgPaid       uint                         `json:"avg_paid"`
	List          []PerformaHostDetailResponse `json:"list"`
}

type PerformaHostDetailResponse struct {
	AccountName string    `json:"account_name"`
	Duration    time.Time `json:"duration"`
	Sales       uint      `json:"sales"`
	Paid        uint      `json:"paid"`
}

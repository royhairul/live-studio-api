package params

import "time"

// Performa Host
type PerformaStudioResponse struct {
	ID            string                         `json:"id"`
	Name          string                         `json:"name"`
	Date          *time.Time                     `json:"date"`
	TotalDuration int64                          `json:"total_duration"`
	TotalSales    uint                           `json:"total_sales"`
	TotalPaid     uint                           `json:"total_paid"`
	AvgSales      uint                           `json:"avg_sales"`
	AvgPaid       uint                           `json:"avg_paid"`
	List          []PerformaStudioDetailResponse `json:"list"`
}

type PerformaStudioDetailResponse struct {
	StudioName string    `json:"studio_name"`
	Duration   time.Time `json:"duration"`
	Sales      uint      `json:"sales"`
	Paid       uint      `json:"paid"`
}

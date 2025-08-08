package params

import (
	"time"

	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/host/entity"
)

type HostResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	StudioID   uint      `json:"studio_id"`
	StudioName string    `json:"studio_name"`
}

func NewHostResponse(host *entity.Host) *HostResponse {
	return &HostResponse{
		ID:         *host.ID,
		Name:       host.Name,
		Phone:      host.Phone,
		StudioID:   host.Studio.ID,
		StudioName: host.Studio.Name,
	}
}

type HostGroupedByStudioResponse struct {
	StudioID   uint           `json:"studio_id"`
	StudioName string         `json:"studio_name"`
	Hosts      []HostResponse `json:"hosts"`
}

// Performa Host
type HostPerformaResponse struct {
	ID            string                       `json:"id"`
	Name          string                       `json:"name"`
	Date          *time.Time                   `json:"date"`
	TotalDuration int64                        `json:"total_duration"`
	TotalSales    uint                         `json:"total_sales"`
	TotalPaid     uint                         `json:"total_paid"`
	AvgSales      uint                         `json:"avg_sales"`
	AvgPaid       uint                         `json:"avg_paid"`
	List          []HostPerformaDetailResponse `json:"list"`
}

type HostPerformaDetailResponse struct {
	AccountName string    `json:"account_name"`
	Duration    time.Time `json:"duration"`
	Sales       uint      `json:"sales"`
	Paid        uint      `json:"paid"`
}

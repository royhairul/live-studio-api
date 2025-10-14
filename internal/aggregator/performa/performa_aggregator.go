package performa

import (
	"time"

	performaparam "github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

type PerformaAggregator interface {
	Calculate(startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error)
	CalculateByHosts(startDate, endDate *time.Time) ([]*performaparam.PerformaHostSummaryResponse, error)
	CalculateByHost(host_id string, startDate, endDate *time.Time) (performaparam.PerformaHostDetailResponse, error)
	CalculateByStudio(studio_id string, startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error)
}

type TotalPerformaHost struct {
	Duration int64
	GMVSales int64
	GMVPaid  int64
	AvgSales int64
	AvgPaid  int64
}

type TotalPerformaAccount struct {
	GMV               int64
	Ads               int64
	CommissionTotal   int64
	CommissionPaid    int64
	CommissionPending int64
	Income            int64
}

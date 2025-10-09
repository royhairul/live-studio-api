package performa

import (
	"time"

	performaparam "github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

type PerformaAggregator interface {
	Calculate(startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerforma, error)
	CalculateByStudio(studio_id string, startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerforma, error)
}

type TotalPerforma struct {
	GMV               int64
	Ads               int64
	CommissionTotal   int64
	CommissionPaid    int64
	CommissionPending int64
	Income            int64
}

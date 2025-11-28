package performa

import (
	"context"
	"time"

	performaparam "github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

type PerformaAggregator interface {
	Calculate(ctx context.Context, startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error)
	CalculateByHosts(ctx context.Context, startDate, endDate *time.Time) ([]*performaparam.PerformaHostSummaryResponse, error)
	CalculateByHost(ctx context.Context, host_id string, startDate, endDate *time.Time) (performaparam.PerformaHostDetailResponse, error)
	CalculateByStudio(ctx context.Context, studio_id string, startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error)
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

package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

type PerformaService interface {
	GetHosts(ctx context.Context, startDate string, endDate string) ([]*params.PerformaHostSummaryResponse, error)
	GetHostByID(ctx context.Context, id string, startDate string, endDate string) (*params.PerformaHostDetailResponse, error)

	GetAccounts(ctx context.Context, startDate string, endDate string) (*params.PerformaAccountResponse, error)

	GetStudios(ctx context.Context, startDate string, endDate string) (*params.PerformaStudioResponse, error)
	GetStudioByID(ctx context.Context, id string, startDate string, endDate string) (*params.PerformaStudioDetailResponse, error)
}

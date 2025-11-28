package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/dashboard/params"
)

type DashboardService interface {
	// TODO: define service methods
	DashboardAdmin(ctx context.Context, startDate, endDate string) (*params.DashboardResponse, error)
}

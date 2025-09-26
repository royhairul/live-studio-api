package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/params"
)

type DashboardService interface {
	// TODO: define service methods
	DashboardAdmin(startDate, endDate string) (*params.DashboardResponse, error)
}

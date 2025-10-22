package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/finance/params"
)

type FinanceService interface {
	FindAll(startDate, endDate *time.Time) (*params.FinanceResponse, error)

	WithAccountUniqueID(unique_id string) FinanceService
	WithStudioID(studio_id string) FinanceService
	WithStatus(status string) FinanceService
	WithPaymentMethod(method string) FinanceService
}

package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
)

type ShopeeFinanceService interface {
	GetPaymentCommission(cookie string, startDate, endDate *time.Time) (*params.ShopeeFinanceCommissionReportResponse, error)
}

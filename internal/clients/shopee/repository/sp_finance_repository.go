package repository

import "github.com/royhairul/live-studio-api/internal/clients/shopee/params"

type FinanceRepository interface {
	GetShopeeLiveSalesReport(req params.ShopeeLiveFinanceRequest, cookie string) (*params.ShopeeLiveFinanceResponse, error)
	GetShopeeLiveSalesReportRange(req params.ShopeeLiveFinanceRequest, cookie string) ([]params.ShopeeLiveReportItem, error)
	GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItem, error)
}

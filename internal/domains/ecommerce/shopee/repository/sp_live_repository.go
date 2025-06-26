package repository

import "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"

type ShopeeLiveRepository interface {
	GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItemRT, error)
}

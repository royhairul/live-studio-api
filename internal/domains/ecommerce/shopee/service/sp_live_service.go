package service

import "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"

type ShopeeLiveService interface {
	GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItemRT, error)
}

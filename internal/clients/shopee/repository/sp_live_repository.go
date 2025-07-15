package repository

import "github.com/royhairul/live-studio-api/internal/clients/shopee/params"

type ShopeeLiveRepository interface {
	GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItemRT, error)
}

package service

import (
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/repository"
)

type ShopeeLiveServiceImpl struct {
	repo repository.ShopeeLiveRepository
}

func NewShopeeLiveService(repo repository.ShopeeLiveRepository) ShopeeLiveService {
	return &ShopeeLiveServiceImpl{repo}
}

// GetShopeeLiveRealTime implements ShopeeLiveService.
func (s ShopeeLiveServiceImpl) GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItemRT, error) {
	return s.repo.GetShopeeLiveRealTime(cookie)
}

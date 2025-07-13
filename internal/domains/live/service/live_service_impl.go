package service

import (
	"log"
	"time"

	"github.com/royhairul/live-studio-api/helpers"
	"github.com/royhairul/live-studio-api/internal/domains/live/params"

	ShopeeParams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	AccountRepo "github.com/royhairul/live-studio-api/internal/domains/account/repository"
)

type LiveServiceImpl struct {
	accountRepo   AccountRepo.AccountRepository
	shopeeLiveSvc ShopeeService.ShopeeLiveService
}

func NewLiveService(accountRepo AccountRepo.AccountRepository, shopeeLiveSvc ShopeeService.ShopeeLiveService) LiveService {
	return &LiveServiceImpl{accountRepo, shopeeLiveSvc}
}

// GetLive implements LiveService.
func (l *LiveServiceImpl) GetLive() ([]*params.LiveResponse, error) {
	accounts, err := l.accountRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var allRealtimeData []*params.LiveResponse

	for _, account := range accounts {
		realtimeData, err := l.shopeeLiveSvc.GetShopeeLiveRealTime(account.Cookie)
		if err != nil {
			log.Printf("Failed to get data realtime for account %s: %v", account.Name, err)
			continue
		}

		// Filter berdasarkan tanggal hari ini
		var todayData []*ShopeeParams.ShopeeLiveReportItemRT
		for _, session := range realtimeData {
			if helpers.IsToday(session.StartTime) {
				//  duration
				session.Duration = time.Now().UnixMilli() - session.StartTime

				//  Change to hours
				durationHours := float64(session.Duration) / 3600000.0

				// if duration > 0, set OmsetPerHours
				if durationHours > 0 {
					session.OmsetPerHour = session.ConfirmedSales / durationHours
				} else {
					session.OmsetPerHour = 0
				}
				todayData = append(todayData, &session)
			}
		}

		allRealtimeData = append(allRealtimeData, &params.LiveResponse{
			AccountName: account.Name,
			Total:       len(realtimeData),
			Relive:      len(todayData),
			ReportLive:  todayData,
		})

	}
	return allRealtimeData, nil
}

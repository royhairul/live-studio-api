package service

import (
	"log"
	"time"

	"github.com/royhairul/live-studio-api/helpers"
	"github.com/royhairul/live-studio-api/internal/domains/live/params"

	shopeeparams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	shopeeservice "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
)

type LiveServiceImpl struct {
	accountSvc    accountservice.AccountService
	shopeeLiveSvc shopeeservice.ShopeeLiveService
}

func NewLiveService(accountSvc accountservice.AccountService, shopeeLiveSvc shopeeservice.ShopeeLiveService) LiveService {
	return &LiveServiceImpl{accountSvc, shopeeLiveSvc}
}

// GetLive implements LiveService.
func (l *LiveServiceImpl) GetLive() ([]*params.LiveResponse, error) {
	accounts, err := l.accountSvc.FindAll()
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
		var todayData []*shopeeparams.ShopeeLiveReportItemRT
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

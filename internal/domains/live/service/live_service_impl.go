package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/live/params"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"

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
func (l *LiveServiceImpl) GetLive(ctx context.Context) ([]*params.LiveResponse, error) {
	accounts, err := l.accountSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var allRealtimeData []*params.LiveResponse

	for _, account := range accounts {
		realtimeData, err := l.shopeeLiveSvc.GetLiveSessionRT(account.Cookie)
		if err != nil {
			log.Printf("Failed to get data realtime for account %s: %v", account.Name, err)
			continue
		}

		// Filter berdasarkan tanggal hari ini
		var todayData []shopeeparams.ShopeeLiveReportItemRT
		for _, session := range realtimeData {
			if utils.IsToday(session.StartTime) {
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
				todayData = append(todayData, session)
			}
		}

		allRealtimeData = append(allRealtimeData, &params.LiveResponse{
			AccountID:   fmt.Sprint(account.ID),
			AccountName: account.Name,
			Total:       len(realtimeData),
			Relive:      len(todayData),
			ReportLive:  todayData,
		})

	}
	return allRealtimeData, nil
}

// GetLiveDetail implements LiveService.
func (l *LiveServiceImpl) GetLiveDetail(ctx context.Context, accountID string, sessionID string, productPage string, productPageSize string) (*params.LiveDetailResponse, error) {
	account, err := l.accountSvc.WithID(accountID).FindOne(ctx)
	if err != nil {
		return nil, err
	}

	overview, err := l.shopeeLiveSvc.GetDashboardOverviewRT(account.Cookie, sessionID)
	if err != nil {
		return nil, err
	}

	viewerProfile, err := l.shopeeLiveSvc.GetDashboardViewerRT(account.Cookie, sessionID)
	if err != nil {
		return nil, err
	}

	viewerSource, err := l.shopeeLiveSvc.GetDashboardViewerSourceRT(account.Cookie, sessionID)
	if err != nil {
		return nil, err
	}
	viewerSource.CalculatePercentage()

	buyerProfile, err := l.shopeeLiveSvc.GetDashboardBuyerRT(account.Cookie, sessionID)
	if err != nil {
		return nil, err
	}

	productPageInt, _ := strconv.Atoi(productPage)
	productPageSizeInt, _ := strconv.Atoi(productPageSize)

	productList, err := l.shopeeLiveSvc.GetDashboardProductListRT(account.Cookie, sessionID, productPageInt, productPageSizeInt)
	if err != nil {
		return nil, err
	}

	return &params.LiveDetailResponse{
		AccountID:     fmt.Sprint(account.ID),
		AccountName:   account.Name,
		Overview:      overview,
		ViewerProfile: viewerProfile,
		ViewerSource:  viewerSource,
		BuyerProfile:  buyerProfile,
		Products:      productList,
	}, nil
}

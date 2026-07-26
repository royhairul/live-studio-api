package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/live/entity"
	"github.com/royhairul/live-studio-api/internal/domains/live/params"
	"github.com/royhairul/live-studio-api/internal/domains/live/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"

	shopeeparams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	shopeeservice "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
)

// maxSyncPages caps the pagination walk so a bad totalPage from Shopee cannot
// spin forever.
const maxSyncPages = 50

// windowDays is the widest span liveList/v2 will answer in a single call: a
// timeDim of "31d" comes back empty, so a longer range must be split into
// windows of at most this many days.
const windowDays = 30

type LiveServiceImpl struct {
	accountSvc    accountservice.AccountService
	shopeeLiveSvc shopeeservice.ShopeeLiveService
	repository    repository.LiveRepository
}

func NewLiveService(
	accountSvc accountservice.AccountService,
	shopeeLiveSvc shopeeservice.ShopeeLiveService,
	repository repository.LiveRepository,
) LiveService {
	return &LiveServiceImpl{accountSvc, shopeeLiveSvc, repository}
}

// GetLive implements LiveService. It powers the preview stream, so it reports
// only sessions still on air — Shopee's realtime list also returns sessions that
// have already finished.
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

		// Keep only sessions that are still running. Shopee's duration is left
		// untouched — overwriting it with elapsed-since-start would make every
		// finished session look like it is still on air.
		liveSessions := make([]shopeeparams.ShopeeLiveReportItemRT, 0, len(realtimeData))
		for _, session := range realtimeData {
			if !utils.IsLive(session.StartTime, session.Duration) {
				continue
			}

			session.IsLive = true

			durationHours := float64(session.Duration) / 3600000.0
			if durationHours > 0 {
				session.OmsetPerHour = session.ConfirmedSales / durationHours
			} else {
				session.OmsetPerHour = 0
			}

			liveSessions = append(liveSessions, session)
		}

		allRealtimeData = append(allRealtimeData, &params.LiveResponse{
			AccountID:   fmt.Sprint(account.ID),
			AccountName: account.Name,
			Total:       len(realtimeData),
			Relive:      len(liveSessions),
			ReportLive:  liveSessions,
		})

	}
	return allRealtimeData, nil
}

// GetStoredHistory implements LiveService. History is read from the lives table,
// so it keeps working when an account's Shopee cookie expires — only the sync
// endpoints talk to Shopee.
func (l *LiveServiceImpl) GetStoredHistory(ctx context.Context, filter params.LiveFilter) (*params.StoredLiveResponse, error) {
	total, err := l.repository.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	lives, err := l.repository.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	history := make([]*params.StoredLiveItem, 0, len(lives))
	for _, live := range lives {
		history = append(history, params.NewStoredLiveItem(live))
	}

	totalPage := 0
	if filter.PageSize > 0 {
		totalPage = int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))
	}

	return &params.StoredLiveResponse{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: totalPage,
		History:   history,
	}, nil
}

// SyncHistory implements LiveService. It walks every page of the account's live
// history and upserts each session, so re-running it refreshes figures that
// settle after a stream ends rather than creating duplicates.
func (l *LiveServiceImpl) SyncHistory(ctx context.Context, accountID string, req shopeeparams.ShopeeLiveHistoryRequest) (*params.LiveSyncResponse, error) {
	account, err := l.accountSvc.WithID(accountID).FindOne(ctx)
	if err != nil {
		return nil, err
	}

	result := &params.LiveSyncResponse{
		AccountID:   fmt.Sprint(account.ID),
		AccountName: account.Name,
	}

	page := req.WithDefaults().Page
	for i := 0; i < maxSyncPages; i++ {
		req.Page = page

		batch, err := l.shopeeLiveSvc.GetLiveHistory(account.Cookie, req)
		if err != nil {
			return nil, err
		}

		for _, session := range batch.List {
			live := entity.Live{
				SessionID:  session.SessionID,
				Title:      session.Title,
				CoverImage: session.CoverImage,
				Status:     session.Status,
				StartTime:  timehandler.ParseInt64MilliDate(session.StartTime),
				Duration:   session.Duration,

				Views:            session.Views,
				Viewers:          session.Viewers,
				PeakViews:        session.PeakViews,
				AvgViewsDuration: session.AvgViewsDuration,
				Comments:         session.Comments,
				Likes:            session.Likes,
				FollowersGrowth:  session.FollowersGrowth,
				EngagedUV:        session.EngagedUV,
				AvgEngagedCCU:    session.AvgEngagedCCU,
				ThirtyMinsCount:  session.ThirtyMinsCount,

				Atc:               session.Atc,
				ProductClicks:     session.ProductClicks,
				ConversionRate:    session.ConversionRate,
				PlacedOrders:      session.PlacedOrders,
				PlacedItemSold:    session.PlacedItemSold,
				PlacedSales:       session.PlacedSales,
				ConfirmedOrders:   session.ConfirmedOrders,
				ConfirmedItemSold: session.ConfirmedItemSold,
				ConfirmedSales:    session.ConfirmedSales,
				PaidOrders:        session.PaidOrders,
				PaidSales:         session.PaidSales,

				AccountID: account.ID,
			}

			created, err := l.repository.Upsert(ctx, &live)
			if err != nil {
				return nil, err
			}

			result.Fetched++
			if created {
				result.Created++
			} else {
				result.Updated++
			}
		}

		if len(batch.List) == 0 || page >= batch.TotalPage {
			break
		}
		page++
	}

	return result, nil
}

// SyncAllHistory implements LiveService. A failing account is recorded and
// skipped rather than aborting the run, so one expired cookie cannot block the
// rest — same tolerance GetLive applies.
func (l *LiveServiceImpl) SyncAllHistory(ctx context.Context, req shopeeparams.ShopeeLiveHistoryRequest) ([]*params.LiveSyncResponse, error) {
	accounts, err := l.accountSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*params.LiveSyncResponse, 0, len(accounts))

	for _, account := range accounts {
		accountID := fmt.Sprint(account.ID)

		result, err := l.SyncHistory(ctx, accountID, req)
		if err != nil {
			log.Printf("Failed to sync live history for account %s: %v", account.Name, err)
			results = append(results, &params.LiveSyncResponse{
				AccountID:   accountID,
				AccountName: account.Name,
				Error:       err.Error(),
			})
			continue
		}

		results = append(results, result)
	}

	return results, nil
}

// SyncHistoryRange implements LiveService. liveList/v2 answers at most a 30-day
// window, so a multi-month backfill is walked as a sequence of 30-day windows,
// newest first, each a full SyncHistory (all pages). Adjacent windows share their
// boundary day on purpose; Upsert keys on session_id, so that overlap refreshes a
// row instead of duplicating it, and guarantees no day falls between windows.
func (l *LiveServiceImpl) SyncHistoryRange(
	ctx context.Context,
	accountID string,
	base shopeeparams.ShopeeLiveHistoryRequest,
	start, end time.Time,
) (*params.LiveSyncRangeResponse, error) {
	account, err := l.accountSvc.WithID(accountID).FindOne(ctx)
	if err != nil {
		return nil, err
	}

	result := &params.LiveSyncRangeResponse{
		AccountID:   fmt.Sprint(account.ID),
		AccountName: account.Name,
		StartDate:   timehandler.FormatDate(&start),
		EndDate:     timehandler.FormatDate(&end),
	}

	for cur := end; ; cur = cur.AddDate(0, 0, -windowDays) {
		// Reuse the single-window sync, but force a full walk (page 1) of a
		// 30-day window ending at cur. name/orderBy/sort/pageSize from base are
		// preserved, so every query parameter is still sent to Shopee.
		req := base
		req.Page = 1
		req.TimeDim = "30d"
		req.EndDate = timehandler.FormatDate(&cur)

		window, err := l.SyncHistory(ctx, accountID, req)
		if err != nil {
			return nil, err
		}

		result.Windows++
		result.Fetched += window.Fetched
		result.Created += window.Created
		result.Updated += window.Updated
		result.Details = append(result.Details, params.LiveSyncWindow{
			EndDate: req.EndDate,
			TimeDim: req.TimeDim,
			Fetched: window.Fetched,
			Created: window.Created,
			Updated: window.Updated,
		})

		// Stop once this window already reaches back to (or past) start.
		if lower := cur.AddDate(0, 0, -windowDays); !lower.After(start) {
			break
		}
	}

	return result, nil
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

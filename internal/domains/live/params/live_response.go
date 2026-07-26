package params

import (
	"time"

	shopeeparam "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"github.com/royhairul/live-studio-api/internal/domains/live/entity"
)

type LiveResponse struct {
	AccountID   string                               `json:"id"`
	AccountName string                               `json:"name"`
	Relive      int                                  `json:"relive"`
	Total       int                                  `json:"total"`
	ReportLive  []shopeeparam.ShopeeLiveReportItemRT `json:"reportLive"`
}

// StoredLiveItem is one persisted session, served from the lives table rather
// than proxied from Shopee.
type StoredLiveItem struct {
	SessionID  int64      `json:"session_id"`
	Title      string     `json:"title"`
	CoverImage string     `json:"cover_image"`
	Status     int        `json:"status"`
	StartTime  *time.Time `json:"start_time"`
	Duration   int64      `json:"duration"`

	Views            int     `json:"views"`
	Viewers          int     `json:"viewers"`
	PeakViews        int     `json:"peak_views"`
	AvgViewsDuration int     `json:"avg_views_duration"`
	Comments         int     `json:"comments"`
	Likes            int     `json:"likes"`
	FollowersGrowth  int     `json:"followers_growth"`
	EngagedUV        int     `json:"engaged_uv"`
	AvgEngagedCCU    int     `json:"avg_engaged_ccu"`
	ThirtyMinsCount  int     `json:"thirty_mins_count"`
	Atc              int     `json:"atc"`
	ProductClicks    int     `json:"product_clicks"`
	ConversionRate   float64 `json:"conversion_rate"`

	PlacedOrders      int     `json:"placed_orders"`
	PlacedItemSold    int     `json:"placed_item_sold"`
	PlacedSales       float64 `json:"placed_sales"`
	ConfirmedOrders   int     `json:"confirmed_orders"`
	ConfirmedItemSold int     `json:"confirmed_item_sold"`
	ConfirmedSales    float64 `json:"confirmed_sales"`
	PaidOrders        int     `json:"paid_orders"`
	PaidSales         float64 `json:"paid_sales"`

	OmsetPerHour float64 `json:"omset_per_hour"`

	AccountID   uint   `json:"account_id"`
	AccountName string `json:"account_name"`
	StudioID    uint16 `json:"studio_id"`
	StudioName  string `json:"studio_name"`

	SyncedAt time.Time `json:"synced_at"`
}

// NewStoredLiveItem maps a persisted row to its API shape. omset_per_hour is
// derived here rather than stored, so it always reflects the current figures.
func NewStoredLiveItem(live *entity.Live) *StoredLiveItem {
	item := &StoredLiveItem{
		SessionID:  live.SessionID,
		Title:      live.Title,
		CoverImage: live.CoverImage,
		Status:     live.Status,
		StartTime:  live.StartTime,
		Duration:   live.Duration,

		Views:            live.Views,
		Viewers:          live.Viewers,
		PeakViews:        live.PeakViews,
		AvgViewsDuration: live.AvgViewsDuration,
		Comments:         live.Comments,
		Likes:            live.Likes,
		FollowersGrowth:  live.FollowersGrowth,
		EngagedUV:        live.EngagedUV,
		AvgEngagedCCU:    live.AvgEngagedCCU,
		ThirtyMinsCount:  live.ThirtyMinsCount,
		Atc:              live.Atc,
		ProductClicks:    live.ProductClicks,
		ConversionRate:   live.ConversionRate,

		PlacedOrders:      live.PlacedOrders,
		PlacedItemSold:    live.PlacedItemSold,
		PlacedSales:       live.PlacedSales,
		ConfirmedOrders:   live.ConfirmedOrders,
		ConfirmedItemSold: live.ConfirmedItemSold,
		ConfirmedSales:    live.ConfirmedSales,
		PaidOrders:        live.PaidOrders,
		PaidSales:         live.PaidSales,

		AccountID:   live.AccountID,
		AccountName: live.Account.Name,
		StudioID:    live.Account.StudioID,
		StudioName:  live.Account.Studio.Name,

		SyncedAt: live.UpdatedAt,
	}

	if durationHours := float64(live.Duration) / 3600000.0; durationHours > 0 {
		item.OmsetPerHour = live.ConfirmedSales / durationHours
	}

	return item
}

// StoredLiveResponse is a page of persisted sessions.
type StoredLiveResponse struct {
	Page      int               `json:"page"`
	PageSize  int               `json:"pageSize"`
	Total     int64             `json:"total"`
	TotalPage int               `json:"totalPage"`
	History   []*StoredLiveItem `json:"history"`
}

// LiveSyncResponse reports what one account's sync wrote to the database.
// Error is set when that account failed, so a bad cookie on one account still
// lets the rest of the run report its results.
type LiveSyncResponse struct {
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
	Fetched     int    `json:"fetched"`
	Created     int    `json:"created"`
	Updated     int    `json:"updated"`
	Error       string `json:"error,omitempty"`
}

// LiveSyncWindow reports one 30-day window inside a range sync. Shopee's
// liveList/v2 returns nothing beyond a 30-day timeDim, so a multi-month backfill
// is a sequence of these.
type LiveSyncWindow struct {
	EndDate string `json:"end_date"`
	TimeDim string `json:"time_dim"`
	Fetched int    `json:"fetched"`
	Created int    `json:"created"`
	Updated int    `json:"updated"`
}

// LiveSyncRangeResponse aggregates a multi-window backfill for one account.
// The Fetched/Created/Updated totals sum every window in Details.
type LiveSyncRangeResponse struct {
	AccountID   string           `json:"account_id"`
	AccountName string           `json:"account_name"`
	StartDate   string           `json:"start_date"`
	EndDate     string           `json:"end_date"`
	Windows     int              `json:"windows"`
	Fetched     int              `json:"fetched"`
	Created     int              `json:"created"`
	Updated     int              `json:"updated"`
	Details     []LiveSyncWindow `json:"details"`
	Error       string           `json:"error,omitempty"`
}

type LiveDetailResponse struct {
	AccountID     string                                                                       `json:"id"`
	AccountName   string                                                                       `json:"name"`
	Overview      shopeeparam.ShopeeLiveOverviewResponse                                       `json:"overview"`
	BuyerProfile  []shopeeparam.ShopeeLiveAudienceAnalyticsResponse                            `json:"buyer_profile"`
	ViewerProfile []shopeeparam.ShopeeLiveAudienceAnalyticsResponse                            `json:"viewer_profile"`
	ViewerSource  shopeeparam.ShopeeLiveAudienceAnalyticsResponse                              `json:"viewer_source"`
	Products      shopeeparam.ShopeeApiPaginationResult[shopeeparam.ShopeeLiveProductResponse] `json:"products"`
}

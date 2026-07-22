package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/pkg/snowflakeid"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
)

// Live is one ended Shopee live session, persisted from the liveList/v2 sync.
//
// Column set mirrors the fields Shopee actually returns — there is no end time,
// only StartTime plus Duration. Shopee sends JSON null for several metrics
// (views, peak_views, likes, followers_growth, product_clicks, conversion_rate,
// paid_orders, paid_sales), which land here as 0: a 0 means "not reported".
//
// Studio is intentionally not denormalised — reach it through Account, the same
// way Transaction does.
type Live struct {
	ID         int64  `gorm:"primaryKey"`
	SessionID  int64  `gorm:"uniqueIndex;not null"`
	Title      string `gorm:"type:varchar(255)"`
	CoverImage string `gorm:"type:text"`
	Status     int
	StartTime  *time.Time
	Duration   int64 // milliseconds, as reported by Shopee

	Views            int
	Viewers          int
	PeakViews        int
	AvgViewsDuration int
	Comments         int
	Likes            int
	FollowersGrowth  int
	EngagedUV        int
	AvgEngagedCCU    int
	ThirtyMinsCount  int

	Atc               int
	ProductClicks     int
	ConversionRate    float64
	PlacedOrders      int
	PlacedItemSold    int
	PlacedSales       float64
	ConfirmedOrders   int
	ConfirmedItemSold int
	ConfirmedSales    float64
	PaidOrders        int
	PaidSales         float64

	AccountID uint
	Account   accountentity.Account `gorm:"foreignKey:AccountID;references:ID;preload:Studio"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	tenantdb.TenantBase
}

func (l *Live) BeforeCreate(tx *gorm.DB) error {
	if snowflakeid.Node == nil {
		return fmt.Errorf("snowflake node is not initialized")
	}
	if l.ID == 0 {
		l.ID = snowflakeid.Node.Generate().Int64()
	}
	return nil
}

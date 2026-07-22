package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/royhairul/live-studio-api/internal/domains/live/entity"
	"github.com/royhairul/live-studio-api/internal/domains/live/params"
)

// upsertColumns are refreshed when a session is synced again. Deliberately
// excludes id, session_id, account_id, tenant_id, created_at and deleted_at so a
// re-sync never reassigns ownership or resurrects a deleted row.
var upsertColumns = []string{
	"title",
	"cover_image",
	"status",
	"start_time",
	"duration",
	"views",
	"viewers",
	"peak_views",
	"avg_views_duration",
	"comments",
	"likes",
	"followers_growth",
	"engaged_uv",
	"avg_engaged_ccu",
	"thirty_mins_count",
	"atc",
	"product_clicks",
	"conversion_rate",
	"placed_orders",
	"placed_item_sold",
	"placed_sales",
	"confirmed_orders",
	"confirmed_item_sold",
	"confirmed_sales",
	"paid_orders",
	"paid_sales",
	"updated_at",
}

type LiveRepositoryImpl struct {
	DB *gorm.DB
}

func NewLiveRepository(db *gorm.DB) LiveRepository {
	return &LiveRepositoryImpl{DB: db}
}

// BuildQuery implements LiveRepository.
func (r *LiveRepositoryImpl) BuildQuery(ctx context.Context, filter params.LiveFilter) *gorm.DB {
	query := r.DB.WithContext(ctx).Model(&entity.Live{}).
		Preload("Account").
		Preload("Account.Studio")

	if filter.ID != nil {
		query = query.Where("lives.id = ?", *filter.ID)
	}

	if filter.SessionID != nil {
		query = query.Where("lives.session_id = ?", *filter.SessionID)
	}

	if filter.AccountID != nil {
		query = query.Where("lives.account_id = ?", *filter.AccountID)
	}

	if filter.StartTime != nil && filter.EndTime != nil {
		query = query.Where("lives.start_time::date BETWEEN ? AND ?", filter.StartTime, filter.EndTime)
	}

	if filter.StudioID != nil {
		query = query.Joins("JOIN accounts ON accounts.id = lives.account_id").
			Where("accounts.studio_id = ?", *filter.StudioID)
	}

	return query
}

// Upsert implements LiveRepository.
func (r *LiveRepositoryImpl) Upsert(ctx context.Context, data *entity.Live) (bool, error) {
	db := r.DB.WithContext(ctx)

	// Postgres reports RowsAffected 1 for both branches of ON CONFLICT DO UPDATE,
	// so look first to tell a create apart from a refresh.
	var ids []int64
	if err := db.Model(&entity.Live{}).
		Where("session_id = ?", data.SessionID).
		Limit(1).
		Pluck("id", &ids).Error; err != nil {
		return false, err
	}
	created := len(ids) == 0

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_id"}},
		DoUpdates: clause.AssignmentColumns(upsertColumns),
	}).Create(data).Error; err != nil {
		return false, err
	}

	return created, nil
}

// FindAll implements LiveRepository.
func (r *LiveRepositoryImpl) FindAll(ctx context.Context, filter params.LiveFilter) ([]*entity.Live, error) {
	query := r.BuildQuery(ctx, filter).Order("lives.start_time DESC")

	if filter.PageSize > 0 {
		query = query.Limit(filter.PageSize).Offset(filter.Offset())
	}

	var items []*entity.Live
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Count implements LiveRepository.
func (r *LiveRepositoryImpl) Count(ctx context.Context, filter params.LiveFilter) (int64, error) {
	// Preloads are irrelevant to a count and Postgres rejects them alongside
	// one, so rebuild without them.
	countFilter := filter
	countFilter.Page, countFilter.PageSize = 0, 0

	var total int64
	if err := r.BuildQuery(ctx, countFilter).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// FindOne implements LiveRepository.
func (r *LiveRepositoryImpl) FindOne(ctx context.Context, filter params.LiveFilter) (*entity.Live, error) {
	var item *entity.Live
	if err := r.BuildQuery(ctx, filter).First(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

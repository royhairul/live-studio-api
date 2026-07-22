package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/live/entity"
	"github.com/royhairul/live-studio-api/internal/domains/live/params"
)

type LiveRepository interface {
	BuildQuery(ctx context.Context, filter params.LiveFilter) *gorm.DB

	// Upsert inserts the session or refreshes its metrics if session_id already
	// exists. Reports whether a new row was inserted.
	Upsert(ctx context.Context, data *entity.Live) (bool, error)

	FindAll(ctx context.Context, filter params.LiveFilter) ([]*entity.Live, error)
	FindOne(ctx context.Context, filter params.LiveFilter) (*entity.Live, error)

	// Count returns how many rows match the filter, ignoring its pagination.
	Count(ctx context.Context, filter params.LiveFilter) (int64, error)
}

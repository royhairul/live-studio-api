package repository

import (
	"context"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/target/entity"
)

type TargetRepository interface {
	// TODO: define repository methods
	FindAll(ctx context.Context) ([]*entity.Target, error)
	FindByID(ctx context.Context, id string) (*entity.Target, error)
	FindByDate(ctx context.Context, date time.Time) (*entity.Target, error)
	FindByStudioAndDate(ctx context.Context, studioID string, date time.Time) (*entity.Target, error)
	Create(ctx context.Context, data *entity.Target) (*entity.Target, error)
	Update(ctx context.Context, data *entity.Target) (*entity.Target, error)
	Delete(ctx context.Context, id string) error
}

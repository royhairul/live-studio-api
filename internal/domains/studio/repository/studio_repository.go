package repository

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/studio/entity"
)

type StudioRepository interface {
	FindAll(ctx context.Context) ([]*entity.Studio, error)
	FindByID(ctx context.Context, id string) (*entity.Studio, error)
	Create(ctx context.Context, studio *entity.Studio) error
	Save(ctx context.Context, studio *entity.Studio) error
	Delete(ctx context.Context, id string) error
}

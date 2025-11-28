package repository

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/host/entity"
)

type HostRepository interface {
	FindAll(ctx context.Context) ([]*entity.Host, error)
	FindByID(ctx context.Context, id string) (*entity.Host, error)
	Create(ctx context.Context, host *entity.Host) (*entity.Host, error)
	Update(ctx context.Context, host *entity.Host) (*entity.Host, error)
	Delete(ctx context.Context, id string) error
}

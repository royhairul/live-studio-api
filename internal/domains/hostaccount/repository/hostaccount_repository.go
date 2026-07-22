package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/entity"
	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/params"
)

type HostAccountRepository interface {
	BuildQuery(ctx context.Context, filter params.HostAccountFilter) *gorm.DB

	Create(ctx context.Context, data *entity.HostAccount) (*entity.HostAccount, error)
	Update(ctx context.Context, data *entity.HostAccount) (*entity.HostAccount, error)
	Delete(ctx context.Context, id string) error

	FindAll(ctx context.Context, filter params.HostAccountFilter) ([]*entity.HostAccount, error)
	FindOne(ctx context.Context, filter params.HostAccountFilter) (*entity.HostAccount, error)
}

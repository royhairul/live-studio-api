package repository

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
)

type AccountRepository interface {
	FindAll(ctx context.Context, filter params.AccountFilter) ([]*entity.Account, error)
	FindOne(ctx context.Context, filter params.AccountFilter) (*entity.Account, error)
	Create(ctx context.Context, account *entity.Account) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) (*entity.Account, error)
	Save(ctx context.Context, account *entity.Account) (*entity.Account, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
}

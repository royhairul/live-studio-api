package repository

import (
	"context"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/params"
)

type AccountadsRepository interface {
	// TODO: define repository methods
	FindAll(ctx context.Context, filter params.AccountadsFilter) ([]*entity.Accountads, error)
	FindOne(ctx context.Context, filter params.AccountadsFilter) (*entity.Accountads, error)
	Create(ctx context.Context, data *entity.Accountads) (*entity.Accountads, error)
	Update(ctx context.Context, data *entity.Accountads) (*entity.Accountads, error)
	Delete(ctx context.Context, id string) error

	FindByDateAndAccount(ctx context.Context, date *time.Time, AccountID string) (*entity.Accountads, error)
}

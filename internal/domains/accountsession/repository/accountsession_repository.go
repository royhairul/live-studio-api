package repository

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
)

type AccountsessionRepository interface {
	// TODO: define repository methods
	FindAll(ctx context.Context, filter params.AccountsessionFilter) ([]*entity.Accountsession, error)
	FindOne(ctx context.Context, filter params.AccountsessionFilter) (*entity.Accountsession, error)
	Create(ctx context.Context, data *entity.Accountsession) (*entity.Accountsession, error)
	Update(ctx context.Context, data *entity.Accountsession) (*entity.Accountsession, error)
	Delete(ctx context.Context, id string) error

	FindAllByAttendanceID(ctx context.Context, id string) ([]*entity.Accountsession, error)
	FindAllByAccountID(ctx context.Context, id string) ([]*entity.Accountsession, error)
	FindAllByStudioID(ctx context.Context, id string) ([]*entity.Accountsession, error)
}

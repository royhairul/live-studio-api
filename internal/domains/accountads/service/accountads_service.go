package service

import (
	"context"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/params"
)

type AccountadsService interface {
	FindAll(ctx context.Context) ([]*params.AccountadsResponse, error)
	FindOne(ctx context.Context, id string) (*params.AccountadsResponse, error)
	Create(ctx context.Context, req params.CreateAccountadsRequest) (*params.AccountadsResponse, error)
	CreateOrUpdate(ctx context.Context, req params.CreateAccountadsRequest) (*params.AccountadsResponse, error)
	Update(ctx context.Context, id string, req params.UpdateAccountadsRequest) (*params.AccountadsResponse, error)
	Delete(ctx context.Context, id string) error

	GetTotalAds(ctx context.Context) (*params.AccountadsTotalResponse, error)

	WithAccountID(accountID string) AccountadsService
	WithAccounts(accountIDs []string) AccountadsService
	WithDateRange(startDate, endDate time.Time) AccountadsService
}

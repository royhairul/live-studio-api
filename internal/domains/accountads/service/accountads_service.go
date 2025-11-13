package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/params"
)

type AccountadsService interface {
	FindAll() ([]*params.AccountadsResponse, error)
	FindOne(id string) (*params.AccountadsResponse, error)
	Create(req params.CreateAccountadsRequest) (*params.AccountadsResponse, error)
	CreateOrUpdate(req params.CreateAccountadsRequest) (*params.AccountadsResponse, error)
	Update(id string, req params.UpdateAccountadsRequest) (*params.AccountadsResponse, error)
	Delete(id string) error

	GetTotalAds() (*params.AccountadsTotalResponse, error)

	WithAccountID(accountID string) AccountadsService
	WithAccounts(accountIDs []string) AccountadsService
	WithDateRange(startDate, endDate time.Time) AccountadsService
}

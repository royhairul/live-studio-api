package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/params"
)

type AccountadsService interface {
	FindAll() ([]*params.AccountadsResponse, error)
	FindByID(id string) (*params.AccountadsResponse, error)
	Create(req params.CreateAccountadsRequest) (*params.AccountadsResponse, error)
	Update(id string, req params.UpdateAccountadsRequest) (*params.AccountadsResponse, error)
	Delete(id string) error

	FindByDateAndAccounts(startDate, endDate *time.Time, accountID string) ([]*params.AccountadsResponse, error)
}

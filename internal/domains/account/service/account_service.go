package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
)

type AccountService interface {
	FindAll() ([]*params.AccountResponse, error)
	FindOne() (*params.AccountResponse, error)
	CreateOrUpdate(req params.CreateAccountRequest) (*params.AccountResponse, error)
	Update(id string, req params.UpdateAccountRequest) (*params.AccountResponse, error)
	Delete(id string) error

	WithID(id string) AccountService
	WithUniqueID(uid string) AccountService
	WithStudioID(studioID string) AccountService
}

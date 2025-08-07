package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
)

type AccountService interface {
	FindAll() ([]*params.AccountResponse, error)
	FindById(id string) (*params.AccountResponse, error)
	FindByUniqueId(uid string) (*params.AccountResponse, error)
	CreateOrUpdate(params.CreateAccountRequest) (*params.AccountResponse, error)
	Delete(id string) error

	FindByStudio(studioId string) ([]*params.AccountResponse, error)
}

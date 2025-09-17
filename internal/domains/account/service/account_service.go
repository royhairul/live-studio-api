package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
)

type AccountService interface {
	FindAll() ([]*params.AccountResponse, error)
	FindById(id string) (*params.AccountResponse, error)
	FindByUniqueId(uid string) (*params.AccountResponse, error)
	CreateOrUpdate(req params.CreateAccountRequest) (*params.AccountResponse, error)
	Update(id string, req params.UpdateAccountRequest) (*params.AccountResponse, error)
	Delete(id string) error

	FindByStudio(studioId string) ([]*params.AccountResponse, error)
}

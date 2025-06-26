package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
)

type AccountService interface {
	FindAll() ([]*params.AccountResponse, error)
	FindById(id string) (*params.AccountDetailResponse, error)
	FindByUniqueId(uid string) (*params.AccountDetailResponse, error)
	CreateOrUpdate(params.CreateAccountRequest) (*entity.Account, error)
	Delete(id string) error
}

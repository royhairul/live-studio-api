package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
)

type AccountRepository interface {
	FindAll(filter params.AccountFilter) ([]*entity.Account, error)
	FindOne(filter params.AccountFilter) (*entity.Account, error)
	Create(account *entity.Account) (*entity.Account, error)
	Update(account *entity.Account) (*entity.Account, error)
	Save(account *entity.Account) (*entity.Account, error)
	Delete(id string) error
}

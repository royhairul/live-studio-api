package repository

import "github.com/royhairul/live-studio-api/internal/domains/account/entity"

type AccountRepository interface {
	FindAll() ([]*entity.Account, error)
	FindById(id string) (*entity.Account, error)
	FindByUniqueId(uid string) (*entity.Account, error)
	Create(account *entity.Account) (*entity.Account, error)
	Save(account *entity.Account) (*entity.Account, error)
	Delete(id string) error

	FindByStudio(studioId string) ([]*entity.Account, error)
}

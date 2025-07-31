package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
)

type TransactionRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Transaction, error)
	FindByID(id string) (*entity.Transaction, error)
	Create(data *entity.Transaction) (*entity.Transaction, error)
	Update(data *entity.Transaction) (*entity.Transaction, error)
	Delete(id string) error

	FindByUniqueID(uid string) (*entity.Transaction, error)
}

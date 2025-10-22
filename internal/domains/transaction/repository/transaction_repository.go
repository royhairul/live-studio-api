package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
)

type TransactionRepository interface {
	FindAll(filter params.TransactionFilter) ([]*entity.Transaction, error)
	FindOne(filter params.TransactionFilter) (*entity.Transaction, error)
	Create(data *entity.Transaction) (*entity.Transaction, error)
	Update(data *entity.Transaction) (*entity.Transaction, error)
	Delete(id string) error
}

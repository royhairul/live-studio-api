package repository

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
)

type TransactionRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Transaction, error)
	FindAllByStatus(status string) ([]*entity.Transaction, error)
	FindAllByAccount(accountID string) ([]*entity.Transaction, error)
	FindAllByAccountStatus(accountID string, status string) ([]*entity.Transaction, error)
	FindAllByAccountStatusDate(accountID string, status string, startDate *time.Time, endDate *time.Time) ([]*entity.Transaction, error)
	FindAllByDate(startDate *time.Time, endDate *time.Time) ([]*entity.Transaction, error)
	FindByID(id string) (*entity.Transaction, error)
	Create(data *entity.Transaction) (*entity.Transaction, error)
	Update(data *entity.Transaction) (*entity.Transaction, error)
	Delete(id string) error

	FindByUniqueID(uid string) (*entity.Transaction, error)
}

package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
)

type TransactionService interface {
	FindAll() (*params.TransactionResponse, error)
	FindAllGrouped() ([]*params.TransactionGroupedResponse, error)
	FindOne() (*params.TransactionList, error)
	Create(req params.CreateTransactionRequest) ([]*params.CreatedTransactionResponse, error)
	Update(id string, req params.UpdateTransactionRequest) (*params.TransactionList, error)
	Delete(id string) error

	GetTotalCommission() (*params.Commission, error)

	WithID(id string) TransactionService
	WithAccountID(accountID string) TransactionService
	WithStatus(status string) TransactionService
	WithDate(startTime time.Time, endTime time.Time) TransactionService
}

package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
)

type TransactionService interface {
	FindAll() ([]*params.TransactionResponse, error)
	FindAllByStatus(status string) ([]*params.TransactionResponse, error)
	FindAllByAccountAndStatus(accountID string, status string) (*params.TransactionResponse, error)
	FindAllByAccountStatusDate(accountID string, status string, startTime *time.Time, endDate *time.Time) (*params.TransactionResponse, error)
	FindAllByDate(accountID string, startDate *time.Time, endDate *time.Time) ([]*params.TransactionResponse, error)

	FindByAccount(accountID string) (*params.TransactionResponse, error)
	FindByID(id string) (*params.TransactionResponse, error)

	Create(req params.CreateTransactionRequest) ([]*params.CreatedTransactionResponse, error)
	Update(id string, req params.UpdateTransactionRequest) (*params.TransactionResponse, error)
	Delete(id string) error
}

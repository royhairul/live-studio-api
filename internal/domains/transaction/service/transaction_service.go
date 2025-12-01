package service

import (
	"context"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
)

type TransactionService interface {
	FindAll(ctx context.Context) (*params.TransactionResponse, error)
	FindAllGrouped(ctx context.Context) ([]*params.TransactionGroupedResponse, error)
	FindOne(ctx context.Context) (*params.TransactionList, error)
	Create(ctx context.Context, req params.CreateTransactionRequest) ([]*params.CreatedTransactionResponse, error)
	Update(ctx context.Context, id string, req params.UpdateTransactionRequest) (*params.TransactionList, error)
	Delete(ctx context.Context, id string) error

	GetTotalCommission(ctx context.Context) (*params.Commission, error)

	WithID(id string) TransactionService
	WithAccountID(accountID string) TransactionService
	WithStudioID(studioID string) TransactionService
	WithStatus(status string) TransactionService
	WithDate(startTime time.Time, endTime time.Time) TransactionService
}

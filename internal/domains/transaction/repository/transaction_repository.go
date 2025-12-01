package repository

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
)

type TransactionRepository interface {
	FindAll(ctx context.Context, filter params.TransactionFilter) ([]*entity.Transaction, error)
	FindOne(ctx context.Context, filter params.TransactionFilter) (*entity.Transaction, error)
	Create(ctx context.Context, data *entity.Transaction) (*entity.Transaction, error)
	Update(ctx context.Context, data *entity.Transaction) (*entity.Transaction, error)
	Delete(ctx context.Context, id string) error
}

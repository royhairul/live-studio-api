package service

import "github.com/royhairul/live-studio-api/internal/domains/transaction/params"

type TransactionService interface {
	FindAll() ([]*params.TransactionResponse, error)
	FindByID(id string) (*params.TransactionResponse, error)
	Create(req params.CreateTransactionRequest) ([]*params.CreatedTransactionResponse, error)
	Update(id string, req params.UpdateTransactionRequest) (*params.TransactionResponse, error)
	Delete(id string) error
}

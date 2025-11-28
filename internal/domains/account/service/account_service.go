package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/account/params"
)

type AccountService interface {
	FindAll(ctx context.Context) ([]*params.AccountResponse, error)
	FindOne(ctx context.Context) (*params.AccountResponse, error)
	CreateOrUpdate(ctx context.Context, req params.CreateAccountRequest) (*params.AccountResponse, error)
	Update(ctx context.Context, id string, req params.UpdateAccountRequest) (*params.AccountResponse, error)
	Delete(ctx context.Context, id string) error

	WithID(id string) AccountService
	WithUniqueID(uid string) AccountService
	WithStudioID(studioID string) AccountService
}

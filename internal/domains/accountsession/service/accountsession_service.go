package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
)

type AccountsessionService interface {
	FindAll(ctx context.Context) ([]*params.AccountsessionResponse, error)
	FindOne(ctx context.Context) (*params.AccountsessionResponse, error)
	Create(ctx context.Context, req params.CreateAccountsessionRequest) (*params.AccountsessionResponse, error)
	Update(ctx context.Context, id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error)
	Delete(ctx context.Context, id string) error

	UpdateEndSession(ctx context.Context, id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error)

	WithID(id string) AccountsessionService
	WithAttendanceID(attendanceID string) AccountsessionService
	WithAccountID(accountID string) AccountsessionService
	WithStudioID(studioID string) AccountsessionService
}

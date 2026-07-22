package service

import (
	"context"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/params"
)

type HostAccountService interface {
	FindAll(ctx context.Context) ([]*params.HostAccountResponse, error)
	FindOne(ctx context.Context) (*params.HostAccountResponse, error)

	Create(ctx context.Context, req params.CreateHostAccountRequest) (*params.HostAccountResponse, error)
	Update(ctx context.Context, id string, req params.UpdateHostAccountRequest) (*params.HostAccountResponse, error)
	Delete(ctx context.Context, id string) error

	// AccountIDsForHost returns the accounts a host held at any point inside the
	// window, each with the span it actually covered — the reporting entry point.
	AccountIDsForHost(ctx context.Context, hostID string, start, end time.Time) ([]params.HostAccountSpan, error)

	WithID(id string) HostAccountService
	WithHostID(hostID string) HostAccountService
	WithAccountID(accountID string) HostAccountService
	WithStudioID(studioID string) HostAccountService
	WithActiveOnly() HostAccountService
}

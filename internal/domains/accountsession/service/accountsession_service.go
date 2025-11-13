package service

import "github.com/royhairul/live-studio-api/internal/domains/accountsession/params"

type AccountsessionService interface {
	FindAll() ([]*params.AccountsessionResponse, error)
	FindOne() (*params.AccountsessionResponse, error)
	Create(req params.CreateAccountsessionRequest) (*params.AccountsessionResponse, error)
	Update(id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error)
	Delete(id string) error

	UpdateEndSession(id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error)

	WithID(id string) AccountsessionService
	WithAttendanceID(attendanceID string) AccountsessionService
	WithAccountID(accountID string) AccountsessionService
	WithStudioID(studioID string) AccountsessionService
}

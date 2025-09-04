package service

import "github.com/royhairul/live-studio-api/internal/domains/accountsession/params"

type AccountsessionService interface {
	FindAll() ([]*params.AccountsessionResponse, error)
	FindByID(id string) (*params.AccountsessionResponse, error)
	Create(req params.CreateAccountsessionRequest) (*params.AccountsessionResponse, error)
	Update(id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error)
	Delete(id string) error

	FindAllByAttendanceID(id string) ([]*params.AccountsessionResponse, error)
	FindAllByAccountID(id string) ([]*params.AccountsessionResponse, error)
	FindAllByStudioID(id string) ([]*params.AccountsessionResponse, error)
	UpdateEndSession(id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error)
}

package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/repository"
)

type AccountsessionServiceImpl struct {
	repository repository.AccountsessionRepository
}

func NewAccountsessionService(repository repository.AccountsessionRepository) AccountsessionService {
	return &AccountsessionServiceImpl{repository}
}

// Create implements AccountsessionService.
func (s *AccountsessionServiceImpl) Create(req params.CreateAccountsessionRequest) (*params.AccountsessionResponse, error) {
	accountsession := entity.Accountsession{
		AccountID:     req.AccountID,
		AttendanceID:  req.AttendanceID,
		GMVSalesStart: req.GMVSalesStart,
		GMVPaidStart:  req.GMVPaidStart,
	}

	created, err := s.repository.Create(&accountsession)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountsessionResponse(created)
	return result, nil
}

// Update implements AccountsessionService.
func (s *AccountsessionServiceImpl) Update(id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error) {
	panic("unimplemented")
}

// UpdateEndSession implements AccountsessionService.
func (s *AccountsessionServiceImpl) UpdateEndSession(id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error) {
	accountsession, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	accountsession.GMVSalesEnd = req.GMVSalesEnd
	accountsession.GMVPaidEnd = req.GMVPaidEnd
	updated, err := s.repository.Update(accountsession)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountsessionResponse(updated)
	return result, nil
}

// FindAll implements AccountsessionService.
func (s *AccountsessionServiceImpl) FindAll() ([]*params.AccountsessionResponse, error) {
	accountSessions, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*params.AccountsessionResponse
	for _, session := range accountSessions {
		result = append(result, params.NewAccountsessionResponse(session))
	}

	return result, nil
}

// FindByID implements AccountsessionService.
func (s *AccountsessionServiceImpl) FindByID(id string) (*params.AccountsessionResponse, error) {
	panic("unimplemented")
}

// Delete implements AccountsessionService.
func (s *AccountsessionServiceImpl) Delete(id string) error {
	panic("unimplemented")
}

// FindByAttendanceID implements AccountsessionService.
func (s *AccountsessionServiceImpl) FindAllByAttendanceID(id string) ([]*params.AccountsessionResponse, error) {
	accountsession, err := s.repository.FindAllByAttendanceID(id)
	if err != nil {
		return nil, err
	}

	var result []*params.AccountsessionResponse
	for _, session := range accountsession {
		accountsessionResponse := params.NewAccountsessionResponse(session)
		result = append(result, accountsessionResponse)
	}

	return result, nil
}

// FindAllByAccountID implements AccountsessionService.
func (s *AccountsessionServiceImpl) FindAllByAccountID(id string) ([]*params.AccountsessionResponse, error) {
	accountsession, err := s.repository.FindAllByAccountID(id)
	if err != nil {
		return nil, err
	}

	var result []*params.AccountsessionResponse
	for _, session := range accountsession {
		accountsessionResponse := params.NewAccountsessionResponse(session)
		result = append(result, accountsessionResponse)
	}

	return result, nil
}

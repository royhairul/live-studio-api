package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/repository"
	"gorm.io/gorm"
)

type AccountsessionServiceImpl struct {
	repository repository.AccountsessionRepository
	options    params.AccountsessionFilter
}

// WithAccountID implements AccountsessionService.
func (s *AccountsessionServiceImpl) WithID(id string) AccountsessionService {
	s.options.ID = &id
	return s
}

// WithAccountID implements AccountsessionService.
func (s *AccountsessionServiceImpl) WithAccountID(accountID string) AccountsessionService {
	s.options.AccountID = &accountID
	return s
}

// WithAttendanceID implements AccountsessionService.
func (s *AccountsessionServiceImpl) WithAttendanceID(attendanceID string) AccountsessionService {
	s.options.AttendanceID = &attendanceID
	return s
}

// WithStudioID implements AccountsessionService.
func (s *AccountsessionServiceImpl) WithStudioID(studioID string) AccountsessionService {
	s.options.StudioID = &studioID
	return s
}

func NewAccountsessionService(repository repository.AccountsessionRepository) AccountsessionService {
	return &AccountsessionServiceImpl{
		repository: repository,
		options:    params.AccountsessionFilter{},
	}
}

// Create implements AccountsessionService.
func (s *AccountsessionServiceImpl) Create(ctx context.Context, req params.CreateAccountsessionRequest) (*params.AccountsessionResponse, error) {
	accountsession := entity.Accountsession{
		AccountID:     req.AccountID,
		AttendanceID:  req.AttendanceID,
		GMVSalesStart: req.GMVSalesStart,
		GMVPaidStart:  req.GMVPaidStart,
	}

	created, err := s.repository.Create(ctx, &accountsession)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountsessionResponse(created)
	return result, nil
}

func (s *AccountsessionServiceImpl) Update(ctx context.Context, id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error) {
	accountsession, err := s.repository.FindOne(ctx, params.AccountsessionFilter{ID: &id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("accountsession with ID %s not found", id)
		}
		return nil, err
	}

	accountsession.GMVSalesEnd = req.GMVSalesEnd
	accountsession.GMVPaidEnd = req.GMVPaidEnd

	updated, err := s.repository.Update(ctx, accountsession)
	if err != nil {
		return nil, err
	}

	return params.NewAccountsessionResponse(updated), nil
}

// UpdateEndSession implements AccountsessionService.
func (s *AccountsessionServiceImpl) UpdateEndSession(ctx context.Context, id string, req params.UpdateEndSessionRequest) (*params.AccountsessionResponse, error) {
	accountsession, err := s.repository.FindOne(ctx, params.AccountsessionFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	accountsession.GMVSalesEnd = req.GMVSalesEnd
	accountsession.GMVPaidEnd = req.GMVPaidEnd
	updated, err := s.repository.Update(ctx, accountsession)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountsessionResponse(updated)
	return result, nil
}

// FindAll implements AccountsessionService.
func (s *AccountsessionServiceImpl) FindAll(ctx context.Context) ([]*params.AccountsessionResponse, error) {
	accountSessions, err := s.repository.FindAll(ctx, s.options)
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
func (s *AccountsessionServiceImpl) FindOne(ctx context.Context) (*params.AccountsessionResponse, error) {
	accountSession, err := s.repository.FindOne(ctx, s.options)
	if err != nil {
		return nil, err
	}
	return params.NewAccountsessionResponse(accountSession), nil
}

// Delete implements AccountsessionService.
func (s *AccountsessionServiceImpl) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

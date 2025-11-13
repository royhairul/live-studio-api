package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/params"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
	"gorm.io/gorm"
)

type AccountadsServiceImpl struct {
	repository repository.AccountadsRepository
	options    params.AccountadsFilter
}

func NewAccountadsService(repository repository.AccountadsRepository) AccountadsService {
	return &AccountadsServiceImpl{
		repository: repository,
		options:    params.AccountadsFilter{},
	}
}

// GetTotalAds implements AccountadsService.
func (s *AccountadsServiceImpl) GetTotalAds() (*params.AccountadsTotalResponse, error) {
	ads, err := s.FindAll()
	if err != nil {
		return nil, err
	}

	var total int64
	for _, ad := range ads {
		total += int64(ad.Ads)
	}

	return &params.AccountadsTotalResponse{
		TotalAds: total,
	}, nil
}

// Create implements AccountadsService.
func (s *AccountadsServiceImpl) Create(req params.CreateAccountadsRequest) (*params.AccountadsResponse, error) {
	date, err := timehandler.ParseDate(req.Date)
	if err != nil {
		return nil, err
	}

	accountAds := entity.Accountads{
		AccountID: req.AccountID,
		Date:      date,
		Spend:     req.Ads,
	}

	created, err := s.repository.Create(&accountAds)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountadsResponse(created)
	return result, nil
}

// CreateOrUpdate implements AccountadsService.
func (s *AccountadsServiceImpl) CreateOrUpdate(req params.CreateAccountadsRequest) (*params.AccountadsResponse, error) {
	// Parse date dari request
	parsedDate, err := timehandler.ParseDate(req.Date)
	if err != nil {
		return nil, err
	}

	accountIDStr := fmt.Sprint(req.AccountID)
	exist, err := s.repository.FindOne(params.AccountadsFilter{
		AccountID: &accountIDStr,
		StartDate: parsedDate,
		EndDate:   parsedDate,
	})
	if err != nil {
		// if error record not found, create a new
		if errors.Is(err, gorm.ErrRecordNotFound) {
			accountAds := entity.Accountads{
				AccountID: req.AccountID,
				Date:      parsedDate,
				Spend:     req.Ads,
			}
			created, err := s.repository.Create(&accountAds)
			if err != nil {
				return nil, err
			}
			return params.NewAccountadsResponse(created), nil
		}
		// other error
		return nil, err
	}

	// if already exist, update record
	exist.AccountID = req.AccountID
	exist.Date = parsedDate
	exist.Spend = req.Ads

	updated, err := s.repository.Update(exist)
	if err != nil {
		return nil, err
	}

	return params.NewAccountadsResponse(updated), nil
}

// Update implements AccountadsService.
func (s *AccountadsServiceImpl) Update(id string, req params.UpdateAccountadsRequest) (*params.AccountadsResponse, error) {
	item, err := s.repository.FindOne(params.AccountadsFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	// update field sesuai request
	if req.AccountID != nil {
		item.AccountID = *req.AccountID
	}
	if req.Date != nil && *req.Date != "" {
		parsedDate, err := timehandler.ParseDate(*req.Date)
		if err != nil {
			return nil, err
		}
		item.Date = parsedDate
	}
	if req.Ads != nil {
		item.Spend = *req.Ads
	}

	// simpan ke repository
	updated, err := s.repository.Update(item)
	if err != nil {
		return nil, err
	}

	// mapping ke response
	result := params.NewAccountadsResponse(updated)
	return result, nil
}

// FindAll implements AccountadsService.
func (s *AccountadsServiceImpl) FindAll() ([]*params.AccountadsResponse, error) {
	items, err := s.repository.FindAll(s.options)
	if err != nil {
		return nil, err
	}

	var result []*params.AccountadsResponse
	for _, item := range items {
		result = append(result, params.NewAccountadsResponse(item))
	}
	return result, nil
}

// FindByID implements AccountadsService.
func (s *AccountadsServiceImpl) FindOne(id string) (*params.AccountadsResponse, error) {
	item, err := s.repository.FindOne(params.AccountadsFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	result := params.NewAccountadsResponse(item)
	return result, nil
}

// Delete implements AccountadsService.
func (s *AccountadsServiceImpl) Delete(id string) error {
	if err := s.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

// WithAccountID implements AccountadsService.
func (s *AccountadsServiceImpl) WithAccountID(accountID string) AccountadsService {
	s.options.AccountID = &accountID
	return s
}

// WithDateRange implements AccountadsService.
func (s *AccountadsServiceImpl) WithDateRange(startDate time.Time, endDate time.Time) AccountadsService {
	s.options.StartDate = &startDate
	s.options.EndDate = &endDate
	return s
}

// WithAccounts implements AccountadsService.
func (s *AccountadsServiceImpl) WithAccounts(accountIDs []string) AccountadsService {
	s.options.AccountIDs = accountIDs
	return s
}

package service

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/helpers/timehandler"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/params"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/repository"
)

type AccountadsServiceImpl struct {
	repository repository.AccountadsRepository
}

func NewAccountadsService(repository repository.AccountadsRepository) AccountadsService {
	return &AccountadsServiceImpl{repository}
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

// Update implements AccountadsService.
func (s *AccountadsServiceImpl) Update(id string, req params.UpdateAccountadsRequest) (*params.AccountadsResponse, error) {
	panic("unimplemented")
}

// FindAll implements AccountadsService.
func (s *AccountadsServiceImpl) FindAll() ([]*params.AccountadsResponse, error) {
	items, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*params.AccountadsResponse
	for _, item := range items {
		result = append(result, params.NewAccountadsResponse(item))
	}
	return result, nil
}

// FindByDateAndAccounts implements AccountadsService.
func (s *AccountadsServiceImpl) FindByDateAndAccounts(startDate, endDate *time.Time, accountID string) ([]*params.AccountadsResponse, error) {
	// item, err := s.repository.FindByDateAndAccount(date, accountID)
	// if err != nil {
	// 	return nil, err
	// }

	var results []*params.AccountadsResponse
	for d := *startDate; !d.After(*endDate); d = d.AddDate(0, 0, 1) {
		item, err := s.repository.FindByDateAndAccount(&d, accountID)
		if err != nil {
			// kalau error query, bisa langsung return atau skip
			return nil, fmt.Errorf("failed on date %s: %w", d.Format("2006-01-02"), err)
		}
		if item != nil {
			results = append(results, params.NewAccountadsResponse(item))
		}
	}

	return results, nil
}

// FindByID implements AccountadsService.
func (s *AccountadsServiceImpl) FindByID(id string) (*params.AccountadsResponse, error) {
	panic("unimplemented")
}

// Delete implements AccountadsService.
func (s *AccountadsServiceImpl) Delete(id string) error {
	panic("unimplemented")
}

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
	item, err := s.repository.FindByID(id)
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
	item, err := s.repository.FindByID(id)
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

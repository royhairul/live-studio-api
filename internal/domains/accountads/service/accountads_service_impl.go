package service

import (
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
	panic("unimplemented")
}

// FindByID implements AccountadsService.
func (s *AccountadsServiceImpl) FindByID(id string) (*params.AccountadsResponse, error) {
	panic("unimplemented")
}

// Delete implements AccountadsService.
func (s *AccountadsServiceImpl) Delete(id string) error {
	panic("unimplemented")
}

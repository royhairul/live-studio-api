package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"github.com/royhairul/live-studio-api/internal/domains/studio/params"
	"github.com/royhairul/live-studio-api/internal/domains/studio/repository"
)

type StudioServiceImpl struct {
	// TODO: add repository dependency
	repository repository.StudioRepository
}

func NewStudioService(repository repository.StudioRepository) StudioService {
	return &StudioServiceImpl{repository}
}

// Create implements StudioService.
func (s *StudioServiceImpl) Create(studioReq params.CreateStudioRequest) (*params.StudioResponse, error) {
	studio := entity.Studio{
		Name:    studioReq.Name,
		Address: studioReq.Address,
	}

	if err := s.repository.Create(&studio); err != nil {
		return nil, err
	}

	result := params.NewStudioResponse(&studio)
	return result, nil
}

// Delete implements StudioService.
func (s *StudioServiceImpl) Delete(id string) error {
	if err := s.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

// FindAll implements StudioService.
func (s *StudioServiceImpl) FindAll() ([]*params.StudioResponse, error) {
	studios, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.StudioResponse
	for _, studio := range studios {
		results = append(results, params.NewStudioResponse(studio))
	}

	return results, nil
}

// FindByID implements StudioService.
func (s *StudioServiceImpl) FindByID(id string) (*params.StudioResponse, error) {
	studio, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewStudioResponse(studio)

	return result, nil
}

// Update implements StudioService.
func (s *StudioServiceImpl) Update(id string, studioReq params.UpdateStudioRequest) (*params.StudioResponse, error) {
	studio, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	studio.Name = *studioReq.Name
	studio.Address = *studioReq.Address

	if err := s.repository.Save(studio); err != nil {
		return nil, err
	}

	result := params.NewStudioResponse(studio)

	return result, nil
}

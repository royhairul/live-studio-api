package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"github.com/royhairul/live-studio-api/internal/domains/studio/params"
	"github.com/royhairul/live-studio-api/internal/domains/studio/repository"
)

type StudioServiceImpl struct {
	repository repository.StudioRepository
}

func NewStudioService(repository repository.StudioRepository) StudioService {
	return &StudioServiceImpl{repository}
}

// Create implements StudioService.
func (s *StudioServiceImpl) Create(ctx context.Context, studioReq params.CreateStudioRequest) (*params.StudioResponse, error) {
	studio := entity.Studio{
		Name:    studioReq.Name,
		Address: studioReq.Address,
	}

	if err := s.repository.Create(ctx, &studio); err != nil {
		return nil, err
	}

	result := params.NewStudioResponse(&studio)
	return result, nil
}

// Delete implements StudioService.
func (s *StudioServiceImpl) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// FindAll implements StudioService.
func (s *StudioServiceImpl) FindAll(ctx context.Context) ([]*params.StudioResponse, error) {
	studios, err := s.repository.FindAll(ctx)
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
func (s *StudioServiceImpl) FindByID(ctx context.Context, id string) (*params.StudioResponse, error) {
	studio, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := params.NewStudioResponse(studio)

	return result, nil
}

// Update implements StudioService.
func (s *StudioServiceImpl) Update(ctx context.Context, id string, studioReq params.UpdateStudioRequest) (*params.StudioResponse, error) {
	studio, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	studio.Name = *studioReq.Name
	studio.Address = *studioReq.Address

	if err := s.repository.Save(ctx, studio); err != nil {
		return nil, err
	}

	result := params.NewStudioResponse(studio)

	return result, nil
}

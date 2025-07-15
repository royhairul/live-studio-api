package service

import "github.com/royhairul/live-studio-api/internal/domains/studio/params"

type StudioService interface {
	// TODO: define service methods
	FindAll() ([]*params.StudioResponse, error)
	FindByID(id string) (*params.StudioResponse, error)
	Create(studioReq params.CreateStudioRequest) (*params.StudioResponse, error)
	Update(id string, studioReq params.UpdateStudioRequest) (*params.StudioResponse, error)
	Delete(id string) error
}

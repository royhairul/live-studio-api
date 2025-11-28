package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/studio/params"
)

type StudioService interface {
	FindAll(ctx context.Context) ([]*params.StudioResponse, error)
	FindByID(ctx context.Context, id string) (*params.StudioResponse, error)
	Create(ctx context.Context, studioReq params.CreateStudioRequest) (*params.StudioResponse, error)
	Update(ctx context.Context, id string, studioReq params.UpdateStudioRequest) (*params.StudioResponse, error)
	Delete(ctx context.Context, id string) error
}

package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/target/params"
)

type TargetService interface {
	FindAll(ctx context.Context, req params.TargetRequest) ([]*params.TargetResponse, error)
	FindByID(ctx context.Context, id string) (*params.TargetResponse, error)
	Create(ctx context.Context, req params.CreateTargetRequest) (*params.CreatedTargetResponse, error)
	CreateOrUpdate(ctx context.Context, req params.CreateTargetRequest) (*params.CreatedTargetResponse, error)
	Update(ctx context.Context, id string, req params.UpdateTargetRequest) (*params.UpdatedTargetResponse, error)
	Delete(ctx context.Context, id string) error
}

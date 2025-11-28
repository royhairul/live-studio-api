package service

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/host/params"
)

type HostService interface {
	FindAll(ctx context.Context) ([]*params.HostResponse, error)
	FindByID(ctx context.Context, id string) (*params.HostResponse, error)
	Create(ctx context.Context, hostReq params.CreateHostRequest) (*params.HostResponse, error)
	Update(ctx context.Context, id string, hostReq params.UpdateHostRequest) (*params.HostResponse, error)
	Delete(ctx context.Context, id string) error

	FindAllGroupedByStudio(ctx context.Context) ([]*params.HostGroupedByStudioResponse, error)
}

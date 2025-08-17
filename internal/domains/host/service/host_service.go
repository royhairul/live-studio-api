package service

import "github.com/royhairul/live-studio-api/internal/domains/host/params"

type HostService interface {
	FindAll() ([]*params.HostResponse, error)
	FindByID(id string) (*params.HostResponse, error)
	Create(hostReq params.CreateHostRequest) (*params.HostResponse, error)
	Update(id string, hostReq params.UpdateHostRequest) (*params.HostResponse, error)
	Delete(id string) error

	FindAllGroupedByStudio() ([]*params.HostGroupedByStudioResponse, error)
}

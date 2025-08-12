package service

import "github.com/royhairul/live-studio-api/internal/domains/performa/params"

type PerformaService interface {
	GetHosts(startTime string, endTime string) ([]*params.PerformaHostResponse, error)
	GetHostByID(id string, startTime string, endTime string) (*params.PerformaHostResponse, error)

	GetStudios(startTime string, endTime string) ([]*params.PerformaStudioResponse, error)
	GetStudioByID(id string, startTime string, endTime string) (*params.PerformaStudioResponse, error)
}

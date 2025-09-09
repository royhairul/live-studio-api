package service

import "github.com/royhairul/live-studio-api/internal/domains/performa/params"

type PerformaService interface {
	GetHosts(startDate string, endDate string) ([]*params.PerformaHostResponse, error)
	GetHostByID(id string, startDate string, endDate string) (*params.PerformaHostDetailResponse, error)

	GetAccounts(startDate string, endDate string) ([]*params.PerformaStudioDetailItemResponse, error)
	GetAccountByID()

	GetStudios(startDate string, endDate string) (*params.PerformaStudioResponse, error)
	GetStudioByID(id string, startDate string, endDate string) (*params.PerformaStudioDetailResponse, error)
}

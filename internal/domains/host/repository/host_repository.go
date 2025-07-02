package repository

import "github.com/royhairul/live-studio-api/internal/domains/host/entity"

type HostRepository interface {
	FindAll() ([]*entity.Host, error)
	FindByID(id string) (*entity.Host, error)
	Create(host *entity.Host) (*entity.Host, error)
	Update(host *entity.Host) (*entity.Host, error)
	Delete(id string) error
}

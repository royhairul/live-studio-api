package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/performa/entity"
)

type PerformaRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Performa, error)
	FindByID(id string) (*entity.Performa, error)
	Create(data *entity.Performa) (*entity.Performa, error)
	Update(data *entity.Performa) (*entity.Performa, error)
	Delete(id string) error
}
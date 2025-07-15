package repository

import "github.com/royhairul/live-studio-api/internal/domains/shift/entity"

type ShiftRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Shift, error)
	FindByID(id uint) (*entity.Shift, error)
	Create(shift *entity.Shift) (*entity.Shift, error)
	Save(shift *entity.Shift) (*entity.Shift, error)
	Delete(id uint) error
}

package repository

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/target/entity"
)

type TargetRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Target, error)
	FindByID(id string) (*entity.Target, error)
	FindByDate(date time.Time) (*entity.Target, error)
	Create(data *entity.Target) (*entity.Target, error)
	Update(data *entity.Target) (*entity.Target, error)
	Delete(id string) error
}

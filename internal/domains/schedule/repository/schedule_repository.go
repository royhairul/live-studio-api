package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
)

type ScheduleRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Schedule, error)
	FindByID(id uint) (*entity.Schedule, error)
	Create(schedule *entity.Schedule) (*entity.Schedule, error)
	Save(schedule *entity.Schedule) (*entity.Schedule, error)
	Delete(id uint) error

	FindByShiftAndDate(shiftID string, date time.Time) ([]*entity.Schedule, error)
	FindByHostShiftAndDate(hostID uuid.UUID, shiftID string, date time.Time) (*entity.Schedule, error)
}

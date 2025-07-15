package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	FindAll() ([]*entity.Attendance, error)
	Create(attendance *entity.Attendance) (*entity.Attendance, error)
	Save(attendance *entity.Attendance) error
	Delete(id string) error

	WithTx(tx *gorm.DB) AttendanceRepository
	BeginTransaction() *gorm.DB

	FindUncheckedOutByHost() ([]*entity.Attendance, error)
	FindByID(id uint) (*entity.Attendance, error)
	FindByScheduleID(id uint) (*entity.Attendance, error)
	FindByHostShiftAndDate(hostID uuid.UUID, shiftID string, date time.Time) (*entity.Attendance, error)
}

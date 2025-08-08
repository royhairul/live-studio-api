package repository

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
)

type AttendanceRepository interface {
	FindAll() ([]*entity.Attendance, error)
	Create(attendance *entity.Attendance) (*entity.Attendance, error)
	Save(attendance *entity.Attendance) error
	Delete(id string) error

	FindUncheckedOutByHost() ([]*entity.Attendance, error)
	FindByID(id uint) (*entity.Attendance, error)
	FindByScheduleID(id uint) (*entity.Attendance, error)
	FindByHostShiftAndDate(hostID string, shiftID uint, date time.Time) (*entity.Attendance, error)
	FindAllByDateRange(startTime time.Time, endTime time.Time) ([]*entity.Attendance, error)
}

package repository

import (
	"context"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
)

type AttendanceRepository interface {
	FindAll(ctx context.Context, filter params.AttendanceFilter) ([]*entity.Attendance, error)
	FindOne(ctx context.Context, filter params.AttendanceFilter) (*entity.Attendance, error)
	Create(ctx context.Context, attendance *entity.Attendance) (*entity.Attendance, error)
	Save(ctx context.Context, attendance *entity.Attendance) error
	Delete(ctx context.Context, id string) error

	FindUncheckedOutByHost(ctx context.Context) ([]*entity.Attendance, error)
	FindUncheckedOutByStudio(ctx context.Context, studioID string, date *time.Time) (*entity.Attendance, error)

	FindByID(ctx context.Context, id uint) (*entity.Attendance, error)
	FindByScheduleID(ctx context.Context, id uint) (*entity.Attendance, error)
	FindByHostShiftAndDate(ctx context.Context, hostID string, shiftID uint, date time.Time) (*entity.Attendance, error)
	FindAllByDateRange(ctx context.Context, startTime *time.Time, endTime *time.Time) ([]*entity.Attendance, error)
	FindAllByHostID(ctx context.Context, id string) ([]*entity.Attendance, error)
	FindByAccountID(ctx context.Context, id uint) (*entity.Attendance, error)
}

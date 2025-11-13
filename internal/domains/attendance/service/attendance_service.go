package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
)

type AttendanceService interface {
	FindAll() ([]*params.AttendanceResponse, error)
	FindUncheckedOut() ([]*params.AttendanceResponse, error)

	CheckIn(req params.AttendanceCheckInRequest) (*params.AttendanceResponse, error)
	CheckOut(req params.AttendanceCheckOutRequest) ([]*params.AttendanceResponse, error)
	GenerateNote(schedule *scheduleentity.Schedule, attendanceDate time.Time, shiftID uint) string

	WithDateRange(startTime, endTime time.Time) AttendanceService
	WithHostID(hostID string) AttendanceService
	WithAccountID(accountID string) AttendanceService
	WithStudioID(studioID string) AttendanceService
}

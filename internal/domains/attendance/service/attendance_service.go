package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
)

type AttendanceService interface {
	FindAll() ([]*params.AttendanceResponse, error)
	FindUncheckedOut() ([]*params.AttendanceResponse, error)
	CheckIn(req params.AttendanceCheckInRequest) (*params.AttendanceCheckInSummary, error)
	CheckOut(req params.AttendanceCheckOutRequest) error
	GenerateNote(schedule *scheduleentity.Schedule, attendanceDate time.Time, shiftID uint) string
}

package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/models"
)

type AttendanceService interface {
	FindAll() ([]*params.AttendanceResponse, error)
	FindUncheckedOut() ([]*params.AttendanceResponse, error)
	CheckIn(req params.AttendanceCheckInRequest) (*params.AttendanceCheckInSummary, error)
	CheckOut(req params.AttendanceCheckOutRequest) error
	GenerateNote(schedule *models.Schedule, attendanceDate time.Time, shiftID uint) string
}

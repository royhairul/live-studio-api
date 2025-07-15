package params

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceCheckInRequest struct {
	HostIDs    []uuid.UUID `json:"host_ids" validate:"required"`
	Date       time.Time   `json:"date" vaidate:"required"`
	ShiftID    string      `json:"shift_id" validate:"required"`
	Attendance string      `json:"attendance" validate:"required"`
}

type AttendanceCheckOutRequest struct {
	AttendanceIDs []uint `json:"attendance_ids" validate:"required"`
}

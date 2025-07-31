package params

import (
	"time"
)

type AttendanceCheckInRequest struct {
	Date     time.Time `json:"date" vaidate:"required"`
	HostID   string    `json:"host_id" validate:"required"`
	ShiftID  uint      `json:"shift_id" validate:"required"`
	StudioID uint      `json:"studio_id" validate:"required"`
}

type AttendanceCheckOutRequest struct {
	AttendanceIDs []uint `json:"attendance_ids" validate:"required"`
}

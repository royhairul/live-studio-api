package params

import (
	"time"
)

type AttendanceCheckInRequest struct {
	HostIDs    []uint    `json:"host_ids" validate:"required"`
	Date       time.Time `json:"date" vaidate:"required"`
	ShiftID    uint      `json:"shift_id" validate:"required"`
	Attendance string    `json:"attendance" validate:"required"`
}

type AttendanceCheckOutRequest struct {
	AttendanceIDs []uint `json:"attendance_ids" validate:"required"`
}

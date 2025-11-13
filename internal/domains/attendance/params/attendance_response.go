package params

import (
	"time"

	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
)

type AttendanceResponse struct {
	ID uint `json:"id"`

	HostID uuid.UUID `json:"host_id"`
	Name   string    `json:"host_name"`

	Date     *time.Time `json:"date"`
	CheckIn  *time.Time `json:"check_in"`
	CheckOut *time.Time `json:"check_out"`
	Duration int64      `json:"duration"`

	ShiftID   uint   `json:"shift_id"`
	ShiftName string `json:"shift_name"`

	StudioID   uint   `json:"studio_id"`
	StudioName string `json:"studio_name"`

	Note string `json:"note"`
}

func NewAttendanceResponse(attendance *entity.Attendance) *AttendanceResponse {
	return &AttendanceResponse{
		ID:         attendance.ID,
		Name:       attendance.Host.Name,
		HostID:     *attendance.Host.ID,
		Date:       attendance.Date,
		CheckIn:    attendance.CheckedInAt,
		CheckOut:   attendance.CheckedOutAt,
		Duration:   attendance.Duration(),
		ShiftID:    attendance.Shift.ID,
		ShiftName:  attendance.Shift.Name,
		StudioID:   attendance.Studio.ID,
		StudioName: attendance.Studio.Name,
		Note:       attendance.Note,
	}
}

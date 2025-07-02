package params

import (
	"time"

	"github.com/google/uuid"
)

type ScheduleRequest struct {
	// TODO: add request fields
}

type CreateScheduleRequest struct {
	HostID   uuid.UUID `json:"host_id" validate:"required"`
	ShiftID  uint      `json:"shift_id" validate:"required"`
	StudioID uint      `json:"studio_id" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
}

type UpdateScheduleRequest struct {
	HostID   *uuid.UUID `json:"host_id" validate:"omitempty"`
	ShiftID  *uint      `json:"shift_id" validate:"omitempty"`
	StudioID *uint      `json:"studio_id" validate:"omitempty"`
	Date     *time.Time `json:"date" validate:"omitempty"`
}

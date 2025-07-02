package params

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/shift/entity"
)

type ShiftResponse struct {
	// TODO: add response fields
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func NewShiftResponse(shift *entity.Shift) *ShiftResponse {
	return &ShiftResponse{
		ID:        shift.ID,
		Name:      shift.Name,
		StartTime: shift.StartTime,
		EndTime:   shift.EndTime,
	}
}

package params

import (
	"time"

	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
)

type ScheduleResponse struct {
	// TODO: add response fields
	HostID    *uuid.UUID `json:"host_id"`
	HostName  string     `json:"host_name"`
	ShiftName string     `json:"shift_name"`
	Date      time.Time  `json:"date"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
}

func NewScheduleResponse(schedule *entity.Schedule) *ScheduleResponse {
	return &ScheduleResponse{
		HostID:    schedule.Host.ID,
		HostName:  schedule.Host.Name,
		ShiftName: schedule.Shift.Name,
		Date:      schedule.Date,
		StartTime: schedule.Shift.StartTime,
		EndTime:   schedule.Shift.EndTime,
	}
}

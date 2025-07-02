package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/schedule/params"
)

type ScheduleService interface {
	// TODO: define service methods
	FindAll() ([]*params.ScheduleResponse, error)
	FindByID(id uint) (*params.ScheduleResponse, error)
	Create(scheduleReq *params.CreateScheduleRequest) (*params.ScheduleResponse, error)
	Update(id uint, scheduleReq *params.UpdateScheduleRequest) (*params.ScheduleResponse, error)
	Delete(id uint) error

	FindByShiftAndDate(shiftID string, dateStr string) ([]*params.ScheduleResponse, error)
}

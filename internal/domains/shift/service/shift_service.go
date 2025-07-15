package service

import "github.com/royhairul/live-studio-api/internal/domains/shift/params"

type ShiftService interface {
	// TODO: define service methods
	FindAll() ([]*params.ShiftResponse, error)
	FindByID(id uint) (*params.ShiftResponse, error)
	Create(shift *params.CreateShiftRequest) (*params.ShiftResponse, error)
	Update(id uint, shift *params.UpdateShiftRequest) (*params.ShiftResponse, error)
	Delete(id uint) error
}

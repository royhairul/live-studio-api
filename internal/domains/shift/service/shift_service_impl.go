package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/shift/entity"
	"github.com/royhairul/live-studio-api/internal/domains/shift/params"
	"github.com/royhairul/live-studio-api/internal/domains/shift/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"
)

type ShiftServiceImpl struct {
	// TODO: add repository dependency
	repository repository.ShiftRepository
}

func NewShiftService(repository repository.ShiftRepository) ShiftService {
	return &ShiftServiceImpl{repository}
}

// Create implements ShiftService.
func (s *ShiftServiceImpl) Create(shiftReq *params.CreateShiftRequest) (*params.ShiftResponse, error) {
	parseStartTime, err := utils.ParseTime(shiftReq.StartTime)
	if err != nil {
		return nil, err
	}

	parseEndTime, err := utils.ParseTime(shiftReq.EndTime)
	if err != nil {
		return nil, err
	}

	shift := entity.Shift{
		Name:      shiftReq.Name,
		StartTime: parseStartTime,
		EndTime:   parseEndTime,
	}
	created, err := s.repository.Create(&shift)
	if err != nil {
		return nil, err
	}

	result := params.NewShiftResponse(created)
	return result, nil
}

// Delete implements ShiftService.
func (s *ShiftServiceImpl) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return err
	}
	return nil
}

// FindAll implements ShiftService.
func (s *ShiftServiceImpl) FindAll() ([]*params.ShiftResponse, error) {
	shifts, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.ShiftResponse
	for _, shift := range shifts {
		results = append(results, params.NewShiftResponse(shift))
	}

	return results, nil
}

// FindByID implements ShiftService.
func (s *ShiftServiceImpl) FindByID(id uint) (*params.ShiftResponse, error) {
	shift, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewShiftResponse(shift)
	return result, nil
}

// Update implements ShiftService.
func (s *ShiftServiceImpl) Update(id uint, shiftReq *params.UpdateShiftRequest) (*params.ShiftResponse, error) {
	parseStartTime, err := utils.ParseTime(*shiftReq.StartTime)
	if err != nil {
		return nil, err
	}

	parseEndTime, err := utils.ParseTime(*shiftReq.EndTime)
	if err != nil {
		return nil, err
	}

	existingShift, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	existingShift.Name = *shiftReq.Name
	existingShift.StartTime = parseStartTime
	existingShift.EndTime = parseEndTime

	updated, err := s.repository.Save(existingShift)
	if err != nil {
		return nil, err
	}
	result := params.NewShiftResponse(updated)
	return result, nil
}

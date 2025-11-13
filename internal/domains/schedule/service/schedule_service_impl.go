package service

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/params"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/repository"
	shiftrepository "github.com/royhairul/live-studio-api/internal/domains/shift/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/constants"
)

type ScheduleServiceImpl struct {
	repository      repository.ScheduleRepository
	shiftRepository shiftrepository.ShiftRepository
}

func NewScheduleService(repository repository.ScheduleRepository, shiftRepository shiftrepository.ShiftRepository) ScheduleService {
	return &ScheduleServiceImpl{repository, shiftRepository}
}

// Create implements ScheduleService.
func (s *ScheduleServiceImpl) Create(scheduleReq *params.CreateScheduleRequest) (*params.ScheduleResponse, error) {
	// take the shift data first
	shift, err := s.shiftRepository.FindByID(scheduleReq.ShiftID)
	if err != nil {
		return nil, err
	}

	schedule := entity.Schedule{
		HostID:  scheduleReq.HostID,
		ShiftID: scheduleReq.ShiftID,
		Date:    scheduleReq.Date,

		StartTime: shift.StartTime,
		EndTime:   shift.EndTime,
	}

	created, err := s.repository.Create(&schedule)
	if err != nil {
		return nil, err
	}

	result := params.NewScheduleResponse(created)

	return result, nil
}

// Delete implements ScheduleService.
func (s *ScheduleServiceImpl) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

// FindAll implements ScheduleService.
func (s *ScheduleServiceImpl) FindAll() ([]*params.ScheduleResponse, error) {
	schedules, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.ScheduleResponse
	for _, schedule := range schedules {
		results = append(results, params.NewScheduleResponse(schedule))
	}

	return results, nil
}

// FindByID implements ScheduleService.
func (s *ScheduleServiceImpl) FindByID(id uint) (*params.ScheduleResponse, error) {
	schedule, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewScheduleResponse(schedule)
	return result, nil
}

// Update implements ScheduleService.
func (s *ScheduleServiceImpl) Update(id uint, scheduleReq *params.UpdateScheduleRequest) (*params.ScheduleResponse, error) {
	// Ambil schedule lama dulu
	schedule, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Jika ShiftID diupdate, ambil shift baru
	if scheduleReq.ShiftID != nil && *scheduleReq.ShiftID != schedule.ShiftID {
		shift, err := s.shiftRepository.FindByID(*scheduleReq.ShiftID)
		if err != nil {
			return nil, err
		}
		schedule.ShiftID = shift.ID
		schedule.StartTime = shift.StartTime
		schedule.EndTime = shift.EndTime
	}

	// Update HostID jika disediakan
	if scheduleReq.HostID != nil {
		schedule.HostID = *scheduleReq.HostID
	}

	// Update Date jika disediakan
	if !scheduleReq.Date.IsZero() {
		schedule.Date = *scheduleReq.Date
	}

	updated, err := s.repository.Save(schedule)
	if err != nil {
		return nil, err
	}

	result := params.NewScheduleResponse(updated)
	return result, nil
}

// FindByHostShiftAndDate implements ScheduleService.
func (s *ScheduleServiceImpl) FindByShiftAndDate(shiftID string, dateStr string) ([]*params.ScheduleResponse, error) {
	// Convert date string to time.Time if necessary
	date, err := time.Parse(constants.LayoutYYMMDD, dateStr)
	if err != nil {
		return nil, err
	}

	schedules, err := s.repository.FindByShiftAndDate(shiftID, date)
	if err != nil {
		return nil, err
	}

	var results []*params.ScheduleResponse
	for _, schedule := range schedules {
		results = append(results, params.NewScheduleResponse(schedule))
	}

	return results, nil
}

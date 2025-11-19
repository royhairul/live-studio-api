package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScheduleRepositoryImpl struct {
	// TODO: add database instance
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &ScheduleRepositoryImpl{DB: db}
}

// Create implements ScheduleRepository.
func (s *ScheduleRepositoryImpl) Create(schedule *entity.Schedule) (*entity.Schedule, error) {
	if err := s.DB.Create(&schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

// Delete implements ScheduleRepository.
func (s *ScheduleRepositoryImpl) Delete(id uint) error {
	if err := s.DB.Delete(&entity.Schedule{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// FindAll implements ScheduleRepository.
func (s *ScheduleRepositoryImpl) FindAll() ([]*entity.Schedule, error) {
	var schedules []*entity.Schedule
	if err := s.DB.Preload(clause.Associations).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// FindByID implements ScheduleRepository.
func (s *ScheduleRepositoryImpl) FindByID(id uint) (*entity.Schedule, error) {
	var schedule *entity.Schedule
	if err := s.DB.Preload(clause.Associations).Where("id = ?", id).First(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

// Save implements ScheduleRepository.
func (s *ScheduleRepositoryImpl) Save(schedule *entity.Schedule) (*entity.Schedule, error) {
	if err := s.DB.Save(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *ScheduleRepositoryImpl) FindByShiftAndDate(shiftID string, date time.Time) ([]*entity.Schedule, error) {
	var schedules []*entity.Schedule

	err := s.DB.
		Preload("Host").
		Preload("Host.Studio").
		Preload("Shift").
		Where("shift_id = ? OR DATE(date) = ?", shiftID, date.Format(constants.LayoutYYMMDD)).
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *ScheduleRepositoryImpl) FindByHostShiftAndDate(hostID uuid.UUID, shiftID string, date time.Time) (*entity.Schedule, error) {
	var schedule *entity.Schedule

	err := s.DB.
		Preload("Host").
		Preload("Host.Studio").
		Preload("Shift").
		Where("host_id = ? AND shift_id = ? OR DATE(date) = ?", hostID, shiftID, date.Format(constants.LayoutYYMMDD)).
		First(&schedule).Error
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

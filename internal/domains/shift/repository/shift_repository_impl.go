package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/shift/entity"
	"gorm.io/gorm"
)

type ShiftRepositoryImpl struct {
	// TODO: add database instance
	DB *gorm.DB
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &ShiftRepositoryImpl{DB: db}
}

// Create implements ShiftRepository.
func (s *ShiftRepositoryImpl) Create(shift *entity.Shift) (*entity.Shift, error) {
	if err := s.DB.Create(shift).Error; err != nil {
		return nil, err
	}

	return shift, nil
}

// Delete implements ShiftRepository.
func (s *ShiftRepositoryImpl) Delete(id uint) error {
	if err := s.DB.Delete(&entity.Shift{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindAll implements ShiftRepository.
func (s *ShiftRepositoryImpl) FindAll() ([]*entity.Shift, error) {
	var shifts []*entity.Shift
	if err := s.DB.Find(&shifts).Error; err != nil {
		return nil, err
	}

	return shifts, nil
}

// FindByID implements ShiftRepository.
func (s *ShiftRepositoryImpl) FindByID(id uint) (*entity.Shift, error) {
	var shift entity.Shift
	if err := s.DB.Where("id = ?", id).First(&shift).Error; err != nil {
		return nil, err
	}
	return &shift, nil
}

// Save implements ShiftRepository.
func (s *ShiftRepositoryImpl) Save(shift *entity.Shift) (*entity.Shift, error) {
	if err := s.DB.Save(&shift).Error; err != nil {
		return nil, err
	}
	return shift, nil
}

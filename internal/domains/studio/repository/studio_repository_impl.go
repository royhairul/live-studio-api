package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"gorm.io/gorm"
)

type StudioRepositoryImpl struct {
	DB *gorm.DB
}

func NewStudioRepository(db *gorm.DB) StudioRepository {
	return &StudioRepositoryImpl{
		DB: db,
	}
}

// Create implements StudioRepository.
func (s *StudioRepositoryImpl) Create(studio *entity.Studio) error {
	if err := s.DB.Create(studio).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements StudioRepository.
func (s *StudioRepositoryImpl) Delete(id string) error {
	if err := s.DB.Delete(&entity.Studio{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindAll implements StudioRepository.
func (s *StudioRepositoryImpl) FindAll() ([]*entity.Studio, error) {
	var studios []*entity.Studio
	if err := s.DB.Find(&studios).Error; err != nil {
		return nil, err
	}

	return studios, nil
}

// FindByID implements StudioRepository.
func (s *StudioRepositoryImpl) FindByID(id string) (*entity.Studio, error) {
	var studio entity.Studio
	if err := s.DB.Where("id = ?", id).First(&studio).Error; err != nil {
		return nil, err
	}

	return &studio, nil
}

// Save implements StudioRepository.
func (s *StudioRepositoryImpl) Save(studio *entity.Studio) error {
	if err := s.DB.Save(&studio).Error; err != nil {
		return err
	}

	return nil
}

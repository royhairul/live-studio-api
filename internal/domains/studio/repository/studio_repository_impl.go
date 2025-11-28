package repository

import (
	"context"

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
func (s *StudioRepositoryImpl) Create(ctx context.Context, studio *entity.Studio) error {
	if err := s.DB.WithContext(ctx).Create(studio).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements StudioRepository.
func (s *StudioRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := s.DB.WithContext(ctx).Delete(&entity.Studio{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindAll implements StudioRepository.
func (s *StudioRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Studio, error) {
	var studios []*entity.Studio
	if err := s.DB.WithContext(ctx).Find(&studios).Error; err != nil {
		return nil, err
	}

	return studios, nil
}

// FindByID implements StudioRepository.
func (s *StudioRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Studio, error) {
	var studio entity.Studio
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&studio).Error; err != nil {
		return nil, err
	}

	return &studio, nil
}

// Save implements StudioRepository.
func (s *StudioRepositoryImpl) Save(ctx context.Context, studio *entity.Studio) error {
	if err := s.DB.WithContext(ctx).Save(&studio).Error; err != nil {
		return err
	}

	return nil
}

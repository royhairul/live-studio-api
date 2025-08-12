package repository

import (
	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/performa/entity"
)

type PerformaRepositoryImpl struct {
	DB *gorm.DB
}

func NewPerformaRepository(db *gorm.DB) PerformaRepository {
	return &PerformaRepositoryImpl{DB: db}
}

// Create implements PerformaRepository.
func (r *PerformaRepositoryImpl) Create(data *entity.Performa) (*entity.Performa, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements PerformaRepository.
func (r *PerformaRepositoryImpl) FindAll() ([]*entity.Performa, error) {
	var items []*entity.Performa
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements PerformaRepository.
func (r *PerformaRepositoryImpl) FindByID(id string) (*entity.Performa, error) {
	var item entity.Performa
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements PerformaRepository.
func (r *PerformaRepositoryImpl) Update(data *entity.Performa) (*entity.Performa, error) {
	if err := r.DB.Model(&entity.Performa{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements PerformaRepository.
func (r *PerformaRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Performa{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

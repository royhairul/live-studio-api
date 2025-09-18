package repository

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/royhairul/live-studio-api/internal/domains/target/entity"
)

type TargetRepositoryImpl struct {
	DB *gorm.DB
}

func NewTargetRepository(db *gorm.DB) TargetRepository {
	return &TargetRepositoryImpl{DB: db}
}

// Create implements TargetRepository.
func (r *TargetRepositoryImpl) Create(data *entity.Target) (*entity.Target, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Preload(clause.Associations).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements TargetRepository.
func (r *TargetRepositoryImpl) FindAll() ([]*entity.Target, error) {
	var items []*entity.Target
	if err := r.DB.Preload(clause.Associations).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements TargetRepository.
func (r *TargetRepositoryImpl) FindByID(id string) (*entity.Target, error) {
	var item entity.Target
	if err := r.DB.Preload(clause.Associations).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByDate implements TargetRepository.
func (r *TargetRepositoryImpl) FindByDate(date time.Time) (*entity.Target, error) {
	var item *entity.Target
	err := r.DB.Preload(clause.Associations).Where("date::date = ?", date).First(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

// FindByStudioAndDate implements TargetRepository.
func (r *TargetRepositoryImpl) FindByStudioAndDate(studioID string, date time.Time) (*entity.Target, error) {
	var item *entity.Target
	err := r.DB.
		Preload(clause.Associations).
		Where("studio_id = ?", studioID).
		Where("date::date = ?", date).
		First(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

// Update implements TargetRepository.
func (r *TargetRepositoryImpl) Update(data *entity.Target) (*entity.Target, error) {
	if err := r.DB.Model(&entity.Target{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Preload(clause.Associations).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements TargetRepository.
func (r *TargetRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Target{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

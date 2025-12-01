package repository

import (
	"context"
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
func (r *TargetRepositoryImpl) Create(ctx context.Context, data *entity.Target) (*entity.Target, error) {
	if err := r.DB.WithContext(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.WithContext(ctx).Preload(clause.Associations).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements TargetRepository.
func (r *TargetRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Target, error) {
	var items []*entity.Target
	if err := r.DB.WithContext(ctx).Preload(clause.Associations).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements TargetRepository.
func (r *TargetRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Target, error) {
	var item entity.Target
	if err := r.DB.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByDate implements TargetRepository.
func (r *TargetRepositoryImpl) FindByDate(ctx context.Context, date time.Time) (*entity.Target, error) {
	var item *entity.Target
	err := r.DB.WithContext(ctx).Preload(clause.Associations).Where("date::date = ?", date).First(&item).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

// FindByStudioAndDate implements TargetRepository.
func (r *TargetRepositoryImpl) FindByStudioAndDate(ctx context.Context, studioID string, date time.Time) (*entity.Target, error) {
	var item *entity.Target
	err := r.DB.WithContext(ctx).
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
func (r *TargetRepositoryImpl) Update(ctx context.Context, data *entity.Target) (*entity.Target, error) {
	if err := r.DB.WithContext(ctx).Model(&entity.Target{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.WithContext(ctx).Preload(clause.Associations).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements TargetRepository.
func (r *TargetRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := r.DB.WithContext(ctx).Delete(&entity.Target{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

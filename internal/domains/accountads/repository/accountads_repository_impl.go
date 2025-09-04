package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
)

type AccountadsRepositoryImpl struct {
	DB *gorm.DB
}

func NewAccountadsRepository(db *gorm.DB) AccountadsRepository {
	return &AccountadsRepositoryImpl{DB: db}
}

// Create implements AccountadsRepository.
func (r *AccountadsRepositoryImpl) Create(data *entity.Accountads) (*entity.Accountads, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements AccountadsRepository.
func (r *AccountadsRepositoryImpl) FindAll() ([]*entity.Accountads, error) {
	var items []*entity.Accountads
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements AccountadsRepository.
func (r *AccountadsRepositoryImpl) FindByID(id string) (*entity.Accountads, error) {
	var item entity.Accountads
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements AccountadsRepository.
func (r *AccountadsRepositoryImpl) Update(data *entity.Accountads) (*entity.Accountads, error) {
	if err := r.DB.Model(&entity.Accountads{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements AccountadsRepository.
func (r *AccountadsRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Accountads{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// FindByDateAndAccount implements AccountadsRepository.
func (r *AccountadsRepositoryImpl) FindByDateAndAccount(date *time.Time, AccountID string) (*entity.Accountads, error) {
	var item entity.Accountads
	err := r.DB.
		Where("date::date = ?", date.Format("2006-01-02")).
		Where("account_id = ?", AccountID).
		First(&item).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
)

type TransactionRepositoryImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{DB: db}
}

// Create implements TransactionRepository.
func (r *TransactionRepositoryImpl) Create(data *entity.Transaction) (*entity.Transaction, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindAll() ([]*entity.Transaction, error) {
	var items []*entity.Transaction
	if err := r.DB.Preload("Account").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindByID(id string) (*entity.Transaction, error) {
	var item entity.Transaction
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByUniqueID implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindByUniqueID(uid string) (*entity.Transaction, error) {
	var item entity.Transaction
	err := r.DB.Where("unique_id = ?", uid).First(&item).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Update implements TransactionRepository.
func (r *TransactionRepositoryImpl) Update(data *entity.Transaction) (*entity.Transaction, error) {
	if err := r.DB.Model(&entity.Transaction{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements TransactionRepository.
func (r *TransactionRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Transaction{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

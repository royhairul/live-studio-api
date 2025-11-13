package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
)

type TransactionRepositoryImpl struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{DB: db}
}

// Query implements TransactionRepository.
func (r *TransactionRepositoryImpl) BuildQuery(filter params.TransactionFilter) *gorm.DB {
	query := r.DB.Model(&entity.Transaction{}).
		Preload(clause.Associations).
		Preload("Account.Studio")

	if filter.ID != nil {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.UniqueID != nil {
		query = query.Where("unique_id = ?", filter.UniqueID)
	}

	if filter.StartTime != nil && filter.EndTime != nil {
		query = query.Where("purchase_time::date BETWEEN ? AND ?", filter.StartTime, filter.EndTime)
	}

	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.StudioID != nil {
		query = query.Joins("JOIN accounts ON accounts.id = transactions.account_id").
			Where("accounts.studio_id = ?", *filter.StudioID)
	}

	return query
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
func (r *TransactionRepositoryImpl) FindAll(filter params.TransactionFilter) ([]*entity.Transaction, error) {
	var items []*entity.Transaction
	if err := r.BuildQuery(filter).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindOne implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindOne(filter params.TransactionFilter) (*entity.Transaction, error) {
	var item *entity.Transaction
	if err := r.BuildQuery(filter).First(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
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

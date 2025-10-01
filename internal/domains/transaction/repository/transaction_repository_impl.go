package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

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
	query := r.DB.Model(&entity.Transaction{}).Preload("Account")

	if filter.StartTime != nil && filter.EndTime != nil {
		query = query.Where("purchase_time::date BETWEEN ? AND ?", filter.StartTime, filter.EndTime)
	}
	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
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

// WithStatus implements TransactionRepository.
func (r *TransactionRepositoryImpl) WithStatus(status string) {
	panic("unimplemented")
}

// FindByID implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindByID(id string) (*entity.Transaction, error) {
	var item entity.Transaction
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindByStatus implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindAllByStatus(status string) ([]*entity.Transaction, error) {
	var items []*entity.Transaction
	if err := r.DB.Preload("Account").Find(&items, "status = ?", status).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAllByAccount implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindAllByAccount(accountID string) ([]*entity.Transaction, error) {
	var items []*entity.Transaction
	if err := r.DB.Preload("Account").Find(&items, "account_id = ?", accountID).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAllByAccountStatus implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindAllByAccountStatus(accountID string, status string) ([]*entity.Transaction, error) {
	var items []*entity.Transaction
	if err := r.DB.Preload("Account").Find(&items, "account_id = ? AND status = ? ", accountID, status).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAllByAccountStatusDate implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindAllByAccountStatusDate(accountID string, status string, startDate *time.Time, endDate *time.Time) ([]*entity.Transaction, error) {
	var items []*entity.Transaction
	err := r.DB.Preload("Account").
		Where("account_id = ?", accountID).
		Where("status = ?", status).
		Where("date::date BETWEEN ? AND ?", startDate, endDate).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FindAllByDate implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindAllByDate(startDate *time.Time, endDate *time.Time) ([]*entity.Transaction, error) {
	var items []*entity.Transaction
	err := r.DB.Preload("Account").
		Where("purchase_time::date BETWEEN ? AND ?", startDate, endDate).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FindByUniqueID implements TransactionRepository.
func (r *TransactionRepositoryImpl) FindByUniqueID(uid string) (*entity.Transaction, error) {
	var item entity.Transaction
	if err := r.DB.Where("unique_id = ?", uid).Take(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
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

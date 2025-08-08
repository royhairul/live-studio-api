package repository

import (
	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
)

type AccountsessionRepositoryImpl struct {
	DB *gorm.DB
}

func NewAccountsessionRepository(db *gorm.DB) AccountsessionRepository {
	return &AccountsessionRepositoryImpl{DB: db}
}

// Create implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) Create(data *entity.Accountsession) (*entity.Accountsession, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) FindAll() ([]*entity.Accountsession, error) {
	var items []*entity.Accountsession
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) FindByID(id string) (*entity.Accountsession, error) {
	var item entity.Accountsession
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) Update(data *entity.Accountsession) (*entity.Accountsession, error) {
	if err := r.DB.Model(&entity.Accountsession{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Accountsession{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// FindByAttendanceID implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) FindAllByAttendanceID(id string) ([]*entity.Accountsession, error) {
	var items []*entity.Accountsession
	if err := r.DB.Preload("Account").Preload("Attendance").Find(&items, "attendance_id = ?", id).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAllByAccountID implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) FindAllByAccountID(id string) ([]*entity.Accountsession, error) {
	var items []*entity.Accountsession
	if err := r.DB.Preload("Account").Preload("Attendance").Find(&items, "account_id = ?", id).Error; err != nil {
		return nil, err
	}
	return items, nil
}

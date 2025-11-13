package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
)

type AccountsessionRepositoryImpl struct {
	DB *gorm.DB
}

func NewAccountsessionRepository(db *gorm.DB) AccountsessionRepository {
	return &AccountsessionRepositoryImpl{DB: db}
}

func (r *AccountsessionRepositoryImpl) BuildQuery(filter params.AccountsessionFilter) *gorm.DB {
	query := r.DB.Model(entity.Accountsession{}).Preload(clause.Associations)

	if filter.AccountID != nil {
		query = query.Where("account_id = ?", filter.AccountID)
	}
	if filter.AttendanceID != nil {
		query = query.Where("attendance_id = ?", filter.AttendanceID)
	}
	if filter.StudioID != nil {
		query = query.Where("studio_id = ?", filter.StudioID)
	}

	return query
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
func (r *AccountsessionRepositoryImpl) FindAll(filter params.AccountsessionFilter) ([]*entity.Accountsession, error) {
	var items []*entity.Accountsession
	if err := r.BuildQuery(filter).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) FindOne(filter params.AccountsessionFilter) (*entity.Accountsession, error) {
	var item entity.Accountsession
	if err := r.BuildQuery(filter).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAllByStudioID implements AccountsessionRepository.
func (r *AccountsessionRepositoryImpl) FindAllByStudioID(id string) ([]*entity.Accountsession, error) {
	var items []*entity.Accountsession
	if err := r.DB.Where("studio_id = ?", id).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
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
	if err := r.DB.Preload("Studio").Preload("Account").Preload("Attendance").Find(&items, "attendance_id = ?", id).Error; err != nil {
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

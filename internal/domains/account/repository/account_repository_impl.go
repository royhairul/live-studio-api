package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{DB: db}
}

func (a *AccountRepositoryImpl) FindAll() ([]*entity.Account, error) {
	var accounts []*entity.Account
	if err := a.DB.Preload("Studio").Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *AccountRepositoryImpl) FindById(id string) (*entity.Account, error) {
	var account entity.Account
	if err := a.DB.Preload("Studio").First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) FindByUniqueId(uid string) (*entity.Account, error) {
	var account entity.Account
	if err := a.DB.Preload("Studio").First(&account, "unique_id = ?", uid).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) Create(account *entity.Account) (*entity.Account, error) {
	if err := a.DB.Preload("Studio").Create(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountRepositoryImpl) Save(account *entity.Account) (*entity.Account, error) {
	if err := a.DB.Preload("Studio").Save(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountRepositoryImpl) Delete(id string) error {
	if err := a.DB.Delete(&entity.Account{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

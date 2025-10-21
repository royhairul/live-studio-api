package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepositoryImpl struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &AccountRepositoryImpl{DB: db}
}

func (a *AccountRepositoryImpl) BuildQuery(filter params.AccountFilter) *gorm.DB {
	query := a.DB.Debug().Model(entity.Account{}).Preload(clause.Associations)

	if filter.ID != nil {
		query = query.Where("id = ?", filter.ID)
	}
	if filter.StudioID != nil {
		query = query.Where("studio_id = ?", filter.StudioID)
	}
	if filter.UniqueID != nil {
		query = query.Where("unique_id = ?", filter.UniqueID)
	}

	return query
}

func (a *AccountRepositoryImpl) FindAll(filter params.AccountFilter) ([]*entity.Account, error) {
	var accounts []*entity.Account
	if err := a.BuildQuery(filter).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *AccountRepositoryImpl) FindOne(filter params.AccountFilter) (*entity.Account, error) {
	var account entity.Account
	if err := a.BuildQuery(filter).First(&account).Error; err != nil {
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
	if err := a.DB.Create(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

// Update implements AccountRepository.
func (a *AccountRepositoryImpl) Update(account *entity.Account) (*entity.Account, error) {
	if err := a.DB.Model(&entity.Account{}).Where("id = ?", account.ID).Updates(account).Error; err != nil {
		return nil, err
	}
	if err := a.DB.Preload("Studio").First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (a *AccountRepositoryImpl) Save(account *entity.Account) (*entity.Account, error) {
	if err := a.DB.Save(account).Error; err != nil {
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

// FindByStudio implements AccountRepository.
func (a *AccountRepositoryImpl) FindByStudio(studioId string) ([]*entity.Account, error) {
	var accounts []*entity.Account
	if err := a.DB.Preload("Studio").Find(&accounts, "studio_id = ?", studioId).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

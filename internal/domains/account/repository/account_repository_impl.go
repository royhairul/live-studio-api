package repository

import (
	"context"

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

func (a *AccountRepositoryImpl) BuildQuery(ctx context.Context, filter params.AccountFilter) *gorm.DB {
	query := a.DB.WithContext(ctx).Debug().Model(entity.Account{}).Preload(clause.Associations)

	if filter.IncludeDeleted {
		query = query.Unscoped()
	}

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

func (a *AccountRepositoryImpl) FindAll(ctx context.Context, filter params.AccountFilter) ([]*entity.Account, error) {
	var accounts []*entity.Account
	if err := a.BuildQuery(ctx, filter).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *AccountRepositoryImpl) FindOne(ctx context.Context, filter params.AccountFilter) (*entity.Account, error) {
	var account entity.Account
	if err := a.BuildQuery(ctx, filter).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) FindByUniqueId(ctx context.Context, uid string) (*entity.Account, error) {
	var account entity.Account
	if err := a.DB.WithContext(ctx).Preload("Studio").First(&account, "unique_id = ?", uid).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepositoryImpl) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	if err := a.DB.WithContext(ctx).Create(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

// Update implements AccountRepository.
func (a *AccountRepositoryImpl) Update(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	if err := a.DB.WithContext(ctx).Model(&entity.Account{}).Where("id = ?", account.ID).Updates(account).Error; err != nil {
		return nil, err
	}
	if err := a.DB.WithContext(ctx).Preload("Studio").First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (a *AccountRepositoryImpl) Save(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	if err := a.DB.WithContext(ctx).Save(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (a *AccountRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := a.DB.WithContext(ctx).Delete(&entity.Account{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// Restore implements AccountRepository.
func (a *AccountRepositoryImpl) Restore(ctx context.Context, id string) error {
	if err := a.DB.
		WithContext(ctx).
		Unscoped().Model(&entity.Account{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}

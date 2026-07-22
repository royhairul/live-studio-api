package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/entity"
	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/params"
)

type HostAccountRepositoryImpl struct {
	DB *gorm.DB
}

func NewHostAccountRepository(db *gorm.DB) HostAccountRepository {
	return &HostAccountRepositoryImpl{DB: db}
}

func (r *HostAccountRepositoryImpl) BuildQuery(ctx context.Context, filter params.HostAccountFilter) *gorm.DB {
	query := r.DB.WithContext(ctx).Model(&entity.HostAccount{}).
		Preload("Host").
		Preload("Account").
		Preload("Account.Studio")

	if filter.ID != nil {
		query = query.Where("host_accounts.id = ?", *filter.ID)
	}

	if filter.ExcludeID != nil {
		query = query.Where("host_accounts.id <> ?", *filter.ExcludeID)
	}

	if filter.HostID != nil {
		query = query.Where("host_accounts.host_id = ?", *filter.HostID)
	}

	if filter.AccountID != nil {
		query = query.Where("host_accounts.account_id = ?", *filter.AccountID)
	}

	if filter.StudioID != nil {
		// Table-qualified: a bare studio_id would be ambiguous once accounts is joined.
		query = query.Joins("JOIN accounts ON accounts.id = host_accounts.account_id").
			Where("accounts.studio_id = ?", *filter.StudioID)
	}

	if filter.OnlyActive {
		query = query.Where("host_accounts.valid_to IS NULL")
	}

	if filter.ActiveAt != nil {
		query = query.Where("host_accounts.valid_from <= ?", *filter.ActiveAt).
			Where("host_accounts.valid_to IS NULL OR host_accounts.valid_to > ?", *filter.ActiveAt)
	}

	// Overlap test: the assignment started before the window ended, and had not
	// already ended when the window began.
	if filter.StartDate != nil && filter.EndDate != nil {
		query = query.Where("host_accounts.valid_from <= ?", *filter.EndDate).
			Where("host_accounts.valid_to IS NULL OR host_accounts.valid_to > ?", *filter.StartDate)
	}

	return query.Order("host_accounts.valid_from DESC")
}

func (r *HostAccountRepositoryImpl) Create(ctx context.Context, data *entity.HostAccount) (*entity.HostAccount, error) {
	db := r.DB.WithContext(ctx)
	if err := db.Create(data).Error; err != nil {
		return nil, err
	}
	return r.FindOne(ctx, params.HostAccountFilter{ID: idPtr(data.ID)})
}

func (r *HostAccountRepositoryImpl) Update(ctx context.Context, data *entity.HostAccount) (*entity.HostAccount, error) {
	db := r.DB.WithContext(ctx)

	// Select on the nullable column so clearing valid_to (reopening an
	// assignment) is not skipped as a zero value.
	if err := db.Model(&entity.HostAccount{}).
		Where("id = ?", data.ID).
		Select("host_id", "account_id", "valid_from", "valid_to", "note").
		Updates(data).Error; err != nil {
		return nil, err
	}

	return r.FindOne(ctx, params.HostAccountFilter{ID: idPtr(data.ID)})
}

func (r *HostAccountRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.HostAccount{}).Error
}

func (r *HostAccountRepositoryImpl) FindAll(ctx context.Context, filter params.HostAccountFilter) ([]*entity.HostAccount, error) {
	var items []*entity.HostAccount
	if err := r.BuildQuery(ctx, filter).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *HostAccountRepositoryImpl) FindOne(ctx context.Context, filter params.HostAccountFilter) (*entity.HostAccount, error) {
	var item entity.HostAccount
	if err := r.BuildQuery(ctx, filter).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func idPtr(id uint) *string {
	value := fmt.Sprint(id)
	return &value
}

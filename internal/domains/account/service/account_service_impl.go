package service

import (
	"context"
	"errors"
	"fmt"

	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
	"github.com/royhairul/live-studio-api/internal/domains/account/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"
	"gorm.io/gorm"
)

type AccountServiceImpl struct {
	repository repository.AccountRepository
	shopeeSvc  ShopeeService.ShopeeAccountService
	options    params.AccountFilter
}

func NewAccountService(
	repository repository.AccountRepository,
	shopeeSvc ShopeeService.ShopeeAccountService,
) AccountService {
	return &AccountServiceImpl{
		repository: repository,
		shopeeSvc:  shopeeSvc,
		options:    params.AccountFilter{},
	}
}

// WithID implements AccountService.
func (a *AccountServiceImpl) WithID(id string) AccountService {
	instance := *a
	instance.options.ID = &id
	return &instance
}

// WithStudioID implements AccountService.
func (a *AccountServiceImpl) WithStudioID(studioID string) AccountService {
	instance := *a
	instance.options.StudioID = &studioID
	return &instance
}

// WithUniqueID implements AccountService.
func (a *AccountServiceImpl) WithUniqueID(uid string) AccountService {
	instance := *a
	instance.options.UniqueID = &uid
	return &instance
}

func (a *AccountServiceImpl) FindAll(ctx context.Context) ([]*params.AccountResponse, error) {
	accounts, err := a.repository.FindAll(ctx, a.options)
	if err != nil {
		return nil, err
	}

	var result []*params.AccountResponse
	for _, account := range accounts {
		result = append(result, params.NewAccountResponse(account))
	}

	return result, nil
}

func (a *AccountServiceImpl) FindOne(ctx context.Context) (*params.AccountResponse, error) {
	account, err := a.repository.FindOne(ctx, a.options)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountResponse(account)

	return result, err
}

func (a *AccountServiceImpl) CreateOrUpdate(ctx context.Context, req params.CreateAccountRequest) (*params.AccountResponse, error) {
	accountShopee, err := a.shopeeSvc.GetShopeeAccount(req.Cookie)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired cookie")
	}

	account := entity.Account{
		Name:     utils.GetDisplayName(accountShopee.Nickname, accountShopee.Username),
		Username: accountShopee.Username,
		Email:    accountShopee.Email,
		UniqueID: fmt.Sprintf("%d", accountShopee.ShopId),
		Platform: "Shopee",
		Cookie:   req.Cookie,
		StudioID: req.StudioID,
		Device:   req.Device,
	}

	// Cari existing berdasarkan UniqueID + TenantID
	existing, err := a.repository.FindOne(ctx, params.AccountFilter{
		UniqueID:       &account.UniqueID,
		IncludeDeleted: true,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existing != nil {
		// Jika soft deleted, reset DeletedAt
		if existing.DeletedAt.Valid {
			_ = a.repository.Restore(ctx, fmt.Sprint(existing.ID))
			existing.DeletedAt = gorm.DeletedAt{}
		}

		// Update fields
		existing.Name = account.Name
		existing.Username = account.Username
		existing.Email = account.Email
		existing.Cookie = account.Cookie
		existing.Device = account.Device
		existing.StudioID = account.StudioID

		updated, err := a.repository.Update(ctx, existing)
		if err != nil {
			return nil, err
		}
		return params.NewAccountResponse(updated), nil
	}

	// Jika tidak ada record, create baru
	created, err := a.repository.Create(ctx, &account)
	if err != nil {
		return nil, err
	}

	return params.NewAccountResponse(created), nil
}

// Update implements AccountService.
func (a *AccountServiceImpl) Update(ctx context.Context, id string, req params.UpdateAccountRequest) (*params.AccountResponse, error) {
	// Check in database
	existing, err := a.repository.FindOne(ctx, params.AccountFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	updatedAccount, err := a.repository.Save(ctx, existing)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountResponse(updatedAccount)
	return result, nil
}

func (a *AccountServiceImpl) Delete(ctx context.Context, id string) error {
	if err := a.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

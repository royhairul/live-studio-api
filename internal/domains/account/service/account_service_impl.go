package service

import (
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

func (a *AccountServiceImpl) FindAll() ([]*params.AccountResponse, error) {
	accounts, err := a.repository.FindAll(a.options)
	if err != nil {
		return nil, err
	}

	var result []*params.AccountResponse
	for _, account := range accounts {
		result = append(result, params.NewAccountResponse(account))
	}

	return result, nil
}

func (a *AccountServiceImpl) FindOne() (*params.AccountResponse, error) {
	account, err := a.repository.FindOne(a.options)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountResponse(account)

	return result, err
}

func (a *AccountServiceImpl) CreateOrUpdate(req params.CreateAccountRequest) (*params.AccountResponse, error) {
	// Get Shopee Account info
	accountShopee, err := a.shopeeSvc.GetShopeeAccount(req.Cookie)
	if err != nil {
		return nil, err
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

	existing, err := a.repository.FindOne(params.AccountFilter{UniqueID: &account.UniqueID, IncludeDeleted: true})
	if err != nil {
		// Create new if not found
		if existing == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			created, err := a.repository.Create(&account)
			if err != nil {
				return nil, err
			}
			return params.NewAccountResponse(created), nil
		}
		return nil, err
	}

	if existing.DeletedAt.Valid {
		existing.DeletedAt = gorm.DeletedAt{}
	}

	// Update fields if exists
	existing.Name = account.Name
	existing.Username = account.Username
	existing.Email = account.Email
	existing.Cookie = account.Cookie
	existing.Device = account.Device
	existing.StudioID = account.StudioID

	updated, err := a.repository.Save(existing)
	if err != nil {
		return nil, err
	}

	return params.NewAccountResponse(updated), nil
}

// Update implements AccountService.
func (a *AccountServiceImpl) Update(id string, req params.UpdateAccountRequest) (*params.AccountResponse, error) {
	// Check in database
	existing, err := a.repository.FindOne(params.AccountFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	updatedAccount, err := a.repository.Save(existing)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountResponse(updatedAccount)
	return result, nil
}

func (a *AccountServiceImpl) Delete(id string) error {
	if err := a.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

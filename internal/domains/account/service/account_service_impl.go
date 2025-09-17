package service

import (
	"strconv"

	"github.com/royhairul/live-studio-api/helpers"
	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
	"github.com/royhairul/live-studio-api/internal/domains/account/repository"
)

type AccountServiceImpl struct {
	repository repository.AccountRepository
	shopeeSvc  ShopeeService.ShopeeAccountService
}

func NewAccountService(repository repository.AccountRepository, shopeeSvc ShopeeService.ShopeeAccountService) AccountService {
	return &AccountServiceImpl{repository, shopeeSvc}
}

func (a *AccountServiceImpl) FindAll() ([]*params.AccountResponse, error) {
	accounts, err := a.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*params.AccountResponse
	for _, account := range accounts {
		result = append(result, params.NewAccountResponse(account))
	}

	return result, nil
}

func (a *AccountServiceImpl) FindById(id string) (*params.AccountResponse, error) {
	account, err := a.repository.FindById(id)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountResponse(account)

	return result, err
}

func (a *AccountServiceImpl) FindByUniqueId(uid string) (*params.AccountResponse, error) {
	account, err := a.repository.FindByUniqueId(uid)
	if err != nil {
		return nil, err
	}

	result := params.NewAccountResponse(account)

	return result, err
}

func (a *AccountServiceImpl) CreateOrUpdate(req params.CreateAccountRequest) (*params.AccountResponse, error) {
	// Get Shopee Account
	accountShopee, err := a.shopeeSvc.GetShopeeAccount(req.Cookie)
	if err != nil {
		return nil, err
	}

	account := entity.Account{
		Name:     helpers.GetDisplayName(accountShopee.Nickname, accountShopee.Username),
		Username: accountShopee.Username,
		Email:    accountShopee.Email,
		UniqueID: strconv.FormatInt(int64(accountShopee.ShopId), 10),
		Platform: "Shopee",
		Cookie:   req.Cookie,
		StudioID: req.StudioID,
	}

	var result *params.AccountResponse

	// Check in database
	existing, err := a.repository.FindByUniqueId(account.UniqueID)
	if err != nil {
		account, err := a.repository.Save(&account)
		if err != nil {
			return nil, err
		}

		result = params.NewAccountResponse(account)
	}

	updatedAccount, err := a.repository.Save(existing)
	if err != nil {
		return nil, err
	}

	result = params.NewAccountResponse(updatedAccount)

	return result, nil
}

// Update implements AccountService.
func (a *AccountServiceImpl) Update(id string, req params.UpdateAccountRequest) (*params.AccountResponse, error) {
	// Check in database
	existing, err := a.repository.FindById(id)
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

// FindByStudio implements AccountService.
func (a *AccountServiceImpl) FindByStudio(studioId string) ([]*params.AccountResponse, error) {
	accounts, err := a.repository.FindByStudio(studioId)
	if err != nil {
		return nil, err
	}

	var result []*params.AccountResponse
	for _, account := range accounts {
		result = append(result, params.NewAccountResponse(account))
	}

	return result, nil
}

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

	accountsResp := params.NewAccountResponse(accounts)
	return accountsResp, nil
}

func (a *AccountServiceImpl) FindById(id string) (*params.AccountDetailResponse, error) {
	account, err := a.repository.FindById(id)
	if err != nil {
		return nil, err
	}

	accountDetailResp := params.AccountDetailResponse{
		ID:         account.ID,
		Name:       account.Name,
		Platform:   account.Platform,
		UniqueID:   account.UniqueID,
		Username:   account.Username,
		Email:      account.Email,
		StudioName: account.Studio.Name,
		Cookie:     account.Cookie,
	}

	return &accountDetailResp, err
}

func (a *AccountServiceImpl) FindByUniqueId(uid string) (*params.AccountDetailResponse, error) {
	account, err := a.repository.FindByUniqueId(uid)
	if err != nil {
		return nil, err
	}

	accountDetailResp := params.AccountDetailResponse{
		ID:         account.ID,
		Name:       account.Name,
		Platform:   account.Platform,
		UniqueID:   account.UniqueID,
		Username:   account.Username,
		Email:      account.Email,
		StudioName: account.Studio.Name,
		Cookie:     account.Cookie,
	}

	return &accountDetailResp, err
}

func (a *AccountServiceImpl) CreateOrUpdate(request params.CreateAccountRequest) (*entity.Account, error) {
	// Get Shopee Account
	accountShopee, err := a.shopeeSvc.GetShopeeAccount(request.Cookie)
	if err != nil {
		return nil, err
	}

	account := entity.Account{
		Name:     helpers.GetDisplayName(accountShopee.Nickname, accountShopee.Username),
		Username: accountShopee.Username,
		Email:    accountShopee.Email,
		UniqueID: strconv.FormatInt(int64(accountShopee.ShopId), 10),
		Platform: "Shopee",
		Cookie:   request.Cookie,
		StudioID: request.StudioID,
	}

	// Check in database
	existing, err := a.repository.FindByUniqueId(account.UniqueID)
	if err != nil {
		account, err := a.repository.Save(&account)
		if err != nil {
			return nil, err
		}

		return account, nil
	}

	updatedAccount, err := a.repository.Save(existing)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

func (a *AccountServiceImpl) Delete(id string) error {
	if err := a.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

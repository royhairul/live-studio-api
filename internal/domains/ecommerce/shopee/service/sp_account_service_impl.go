package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"
	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/repository"
)

type shopeeAccountServiceImpl struct {
	repo repository.AccountRepository
}

func NewAccountShopeeService(repo repository.AccountRepository) ShopeeAccountService {
	return &shopeeAccountServiceImpl{repo: repo}
}

func (s *shopeeAccountServiceImpl) GetShopeeAccount(cookie string) (*params.ShopeeAccountResponse, error) {
	return s.repo.GetShopeeAccount(cookie)
}

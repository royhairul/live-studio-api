package service

import (
	"fmt"

	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"go.uber.org/fx"
)

type ShopeeAccountServiceDeps struct {
	fx.In
	ShopeeClient *shopee.ShopeeClient `name:"defaultShopeeClient"`
}

type ShopeeAccountServiceImpl struct {
	ShopeeClient *shopee.ShopeeClient
}

func NewAccountShopeeService(deps ShopeeAccountServiceDeps) ShopeeAccountService {
	return &ShopeeAccountServiceImpl{ShopeeClient: deps.ShopeeClient}
}

func (s *ShopeeAccountServiceImpl) GetShopeeAccount(cookie string) (*params.ShopeeAccountResponse, error) {
	endpoint := "/api/v4/account/basic/get_account_info"

	req, err := s.ShopeeClient.NewShopeeRequest("GET", endpoint, nil, nil, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeAccountResponse]
	if err := s.ShopeeClient.DoShopeeRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return &result.Data, nil
}

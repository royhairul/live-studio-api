package repository

import (
	"fmt"

	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee"
	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"
)

type shopeeAccountRepository struct {
	client *shopee.ShopeeClient
}

func NewShopeeAccountRepository(client *shopee.ShopeeClient) AccountRepository {
	return &shopeeAccountRepository{client: client}
}

func (r *shopeeAccountRepository) GetShopeeAccount(cookie string) (*params.ShopeeAccountResponse, error) {
	endpoint := "/api/v4/account/basic/get_account_info"

	req, err := r.client.NewRequest("GET", endpoint, nil, nil, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
	}

	var result params.ShopeeApiResponse[params.ShopeeAccountResponse]
	if err := r.client.DoRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return &result.Data, nil
}

package repository

import "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"

type AccountRepository interface {
	GetShopeeAccount(cookie string) (*params.ShopeeAccountResponse, error)
}

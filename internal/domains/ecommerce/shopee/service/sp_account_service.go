package service

import "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"

type ShopeeAccountService interface {
	GetShopeeAccount(cookie string) (*params.ShopeeAccountResponse, error)
}

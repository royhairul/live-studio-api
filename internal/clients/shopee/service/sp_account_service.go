package service

import "github.com/royhairul/live-studio-api/internal/clients/shopee/params"

type ShopeeAccountService interface {
	GetShopeeAccount(cookie string) (*params.ShopeeAccountResponse, error)
}

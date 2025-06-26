package service

import (
	ShopeeParams "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/params"

	"github.com/royhairul/live-studio-api/internal/domains/finance/params"
)

type FinanceService interface {
	FindAll(financeReq ShopeeParams.ShopeeLiveFinanceRequest) ([]*params.FinanceResponse, error)
}

package service

import (
	ShopeeParams "github.com/royhairul/live-studio-api/internal/clients/shopee/params"

	"github.com/royhairul/live-studio-api/internal/domains/finance/params"
)

type FinanceService interface {
	FindAll(financeReq ShopeeParams.ShopeeLiveFinanceRequest) ([]*params.FinanceResponse, error)

	FindAllCommission() (*params.CommissionTotalResponse, error)
	FindByStudioCommission() (*params.CommissionStudioDetailResponse, error)
}

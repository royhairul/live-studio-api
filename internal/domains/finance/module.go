package finance

import (
	"go.uber.org/fx"

	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	"github.com/royhairul/live-studio-api/internal/domains/finance/controller"
	"github.com/royhairul/live-studio-api/internal/domains/finance/service"
)

var Module = fx.Module(
	"finance",
	fx.Provide(
		ShopeeService.NewShopeeFinanceService,

		// Finance
		service.NewFinanceService,
		controller.NewFinanceController,
	),
	fx.Invoke(RegisterRoutes),
)

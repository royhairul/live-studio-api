package account

import (
	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	"github.com/royhairul/live-studio-api/internal/domains/account/controller"
	"github.com/royhairul/live-studio-api/internal/domains/account/repository"
	"github.com/royhairul/live-studio-api/internal/domains/account/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"account",
	fx.Provide(
		ShopeeService.NewAccountShopeeService,

		// account
		repository.NewAccountRepository,
		service.NewAccountService,
		controller.NewAccountController,
	),
	fx.Invoke(RegisterRoutes),
)

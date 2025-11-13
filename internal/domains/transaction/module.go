package transaction

import (
	"github.com/royhairul/live-studio-api/internal/domains/transaction/controller"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/repository"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/service"
	"go.uber.org/fx"

	shopeeservice "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
)

var Module = fx.Module(
	"transaction",
	fx.Provide(
		shopeeservice.NewShopeeCheckoutService,
		repository.NewTransactionRepository,
		service.NewTransactionService,
		controller.NewTransactionController,
	),
	fx.Invoke(RegisterRoutes),
)

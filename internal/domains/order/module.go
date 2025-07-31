package order

import (
	"github.com/royhairul/live-studio-api/internal/domains/order/controller"
	"github.com/royhairul/live-studio-api/internal/domains/order/repository"
	"github.com/royhairul/live-studio-api/internal/domains/order/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"order",
	fx.Provide(
		repository.NewOrderRepository,
		service.NewOrderService,
		controller.NewOrderController,
	),
	fx.Invoke(RegisterRoutes),
)

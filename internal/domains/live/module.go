package live

import (
	"go.uber.org/fx"

	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	"github.com/royhairul/live-studio-api/internal/domains/live/controller"
	"github.com/royhairul/live-studio-api/internal/domains/live/service"
)

var Module = fx.Module(
	"live",
	fx.Provide(
		ShopeeService.NewShopeeLiveService,
		service.NewLiveService,
		controller.NewLiveController,
	),
	fx.Invoke(RegisterRoutes),
)

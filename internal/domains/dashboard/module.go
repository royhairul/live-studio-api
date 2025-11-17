package dashboard

import (
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/controller"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"dashboard",
	fx.Provide(
		service.NewDashboardService,
		controller.NewDashboardController,
	),
	fx.Invoke(RegisterRoutes),
)

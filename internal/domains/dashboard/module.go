package dashboard

import (
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/controller"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/repository"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"dashboard",
	fx.Provide(
		repository.NewDashboardRepository,
		service.NewDashboardService,
		controller.NewDashboardController,
	),
	fx.Invoke(RegisterRoutes),
)

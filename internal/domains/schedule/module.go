package schedule

import (
	"github.com/royhairul/live-studio-api/internal/domains/schedule/controller"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/repository"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"schedule",
	fx.Provide(
		repository.NewScheduleRepository,
		service.NewScheduleService,
		controller.NewScheduleController,
	),
	fx.Invoke(RegisterRoutes),
)

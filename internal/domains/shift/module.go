package shift

import (
	"github.com/royhairul/live-studio-api/internal/domains/shift/controller"
	"github.com/royhairul/live-studio-api/internal/domains/shift/repository"
	"github.com/royhairul/live-studio-api/internal/domains/shift/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"shift",
	fx.Provide(
		repository.NewShiftRepository,
		service.NewShiftService,
		controller.NewShiftController,
	),
	fx.Invoke(RegisterRoutes),
)

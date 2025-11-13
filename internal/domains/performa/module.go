package performa

import (
	"github.com/royhairul/live-studio-api/internal/domains/performa/controller"
	"github.com/royhairul/live-studio-api/internal/domains/performa/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"performa",
	fx.Provide(
		service.NewPerformaService,
		controller.NewPerformaController,
	),
	fx.Invoke(RegisterRoutes),
)

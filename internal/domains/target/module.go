package target

import (
	"github.com/royhairul/live-studio-api/internal/domains/target/controller"
	"github.com/royhairul/live-studio-api/internal/domains/target/repository"
	"github.com/royhairul/live-studio-api/internal/domains/target/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"target",
	fx.Provide(
		repository.NewTargetRepository,
		service.NewTargetService,
		controller.NewTargetController,
	),
	fx.Invoke(RegisterRoutes),
)

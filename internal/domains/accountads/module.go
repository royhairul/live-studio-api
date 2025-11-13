package accountads

import (
	"github.com/royhairul/live-studio-api/internal/domains/accountads/controller"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/repository"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"accountads",
	fx.Provide(
		repository.NewAccountadsRepository,
		service.NewAccountadsService,
		controller.NewAccountadsController,
	),
	fx.Invoke(RegisterRoutes),
)

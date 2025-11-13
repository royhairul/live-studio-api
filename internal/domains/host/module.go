package host

import (
	"github.com/royhairul/live-studio-api/internal/domains/host/controller"
	"github.com/royhairul/live-studio-api/internal/domains/host/repository"
	"github.com/royhairul/live-studio-api/internal/domains/host/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"host",
	fx.Provide(
		repository.NewHostRepository,
		service.NewHostService,
		controller.NewHostController,
	),
	fx.Invoke(RegisterRoutes),
)

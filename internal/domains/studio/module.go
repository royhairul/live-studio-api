package studio

import (
	"github.com/royhairul/live-studio-api/internal/domains/studio/controller"
	"github.com/royhairul/live-studio-api/internal/domains/studio/repository"
	"github.com/royhairul/live-studio-api/internal/domains/studio/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"studio",
	fx.Provide(
		repository.NewStudioRepository,
		service.NewStudioService,
		controller.NewStudioController,
	),
	fx.Invoke(RegisterRoutes),
)

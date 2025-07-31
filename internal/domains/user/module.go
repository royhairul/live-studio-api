package user

import (
	"github.com/royhairul/live-studio-api/internal/domains/user/controller"
	"github.com/royhairul/live-studio-api/internal/domains/user/repository"
	"github.com/royhairul/live-studio-api/internal/domains/user/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		repository.NewUserRepository,
		service.NewUserService,
		controller.NewUserController,
	),
	fx.Invoke(RegisterRoutes),
)

package auth

import (
	"github.com/royhairul/live-studio-api/internal/domains/auth/controller"
	"github.com/royhairul/live-studio-api/internal/domains/auth/repository"
	"github.com/royhairul/live-studio-api/internal/domains/auth/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth",
	fx.Provide(
		repository.NewAuthRepository,
		service.NewAuthService,
		controller.NewAuthController,
	),
	fx.Invoke(RegisterRoutes),
)

package role

import (
	"github.com/royhairul/live-studio-api/internal/domains/role/controller"
	"github.com/royhairul/live-studio-api/internal/domains/role/repository"
	"github.com/royhairul/live-studio-api/internal/domains/role/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"role",
	fx.Provide(
		repository.NewRoleRepository,
		service.NewRoleService,
		controller.NewRoleController,
	),
	fx.Invoke(RegisterRoutes),
)

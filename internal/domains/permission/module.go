package permission

import (
	"github.com/royhairul/live-studio-api/internal/domains/permission/controller"
	"github.com/royhairul/live-studio-api/internal/domains/permission/repository"
	"github.com/royhairul/live-studio-api/internal/domains/permission/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"permission",
	fx.Provide(
		repository.NewPermissionRepository,
		service.NewPermissionService,
		controller.NewPermissionController,
	),
	fx.Invoke(RegisterRoutes),
)

package hostaccount

import (
	"go.uber.org/fx"

	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/controller"
	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/repository"
	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/service"
)

var Module = fx.Module(
	"hostaccount",
	fx.Provide(
		repository.NewHostAccountRepository,
		service.NewHostAccountService,
		controller.NewHostAccountController,
	),
	fx.Invoke(RegisterRoutes),
)

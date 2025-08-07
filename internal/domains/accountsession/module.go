package accountsession

import (
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/controller"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/repository"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"accountsession",
	fx.Provide(
		repository.NewAccountsessionRepository,
		service.NewAccountsessionService,
		controller.NewAccountsessionController,
	),
	fx.Invoke(RegisterRoutes),
)

package main

import (
	"github.com/royhairul/live-studio-api/config"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/domains/account"
	"github.com/royhairul/live-studio-api/internal/domains/attendance"
	"github.com/royhairul/live-studio-api/internal/domains/finance"
	"github.com/royhairul/live-studio-api/internal/domains/host"
	"github.com/royhairul/live-studio-api/internal/domains/live"
	"github.com/royhairul/live-studio-api/internal/domains/schedule"
	"github.com/royhairul/live-studio-api/internal/domains/shift"
	"github.com/royhairul/live-studio-api/internal/pkg/httpclient"
	"github.com/royhairul/live-studio-api/internal/pkg/server"
	"github.com/royhairul/live-studio-api/routes"
	"github.com/royhairul/live-studio-api/validators"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			database.ConnectDatabase,
			validators.InitValidator,
			routes.SetupRouter,
			routes.GroupAPI,
			httpclient.NewHttpClient,
			httpclient.NewClient,
			fx.Annotate(shopee.NewShopeeCreatorClient, fx.ResultTags(`name:"creatorShopeeClient"`)),
			fx.Annotate(shopee.NewShopeeSellerClient, fx.ResultTags(`name:"sellerShopeeClient"`)),
			fx.Annotate(shopee.NewShopeeDefaultClient, fx.ResultTags(`name:"defaultShopeeClient"`)),
		),

		fx.Invoke(
			database.MigrateDatabase,
			server.Start,
		),

		// Module Domains
		host.Module,
		shift.Module,
		schedule.Module,
		attendance.Module,
		account.Module,
		live.Module,
		finance.Module,
	)

	app.Run()
}

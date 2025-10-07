package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/config"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/aggregator"
	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/domains/account"
	"github.com/royhairul/live-studio-api/internal/domains/accountads"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession"
	"github.com/royhairul/live-studio-api/internal/domains/attendance"
	"github.com/royhairul/live-studio-api/internal/domains/auth"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard"
	"github.com/royhairul/live-studio-api/internal/domains/finance"
	"github.com/royhairul/live-studio-api/internal/domains/host"
	"github.com/royhairul/live-studio-api/internal/domains/live"
	"github.com/royhairul/live-studio-api/internal/domains/order"
	"github.com/royhairul/live-studio-api/internal/domains/performa"
	"github.com/royhairul/live-studio-api/internal/domains/permission"
	"github.com/royhairul/live-studio-api/internal/domains/product"
	"github.com/royhairul/live-studio-api/internal/domains/role"
	"github.com/royhairul/live-studio-api/internal/domains/schedule"
	"github.com/royhairul/live-studio-api/internal/domains/shift"
	"github.com/royhairul/live-studio-api/internal/domains/studio"
	"github.com/royhairul/live-studio-api/internal/domains/target"
	"github.com/royhairul/live-studio-api/internal/domains/transaction"
	"github.com/royhairul/live-studio-api/internal/domains/user"
	"github.com/royhairul/live-studio-api/internal/pkg/httpclient"
	"github.com/royhairul/live-studio-api/internal/pkg/snowflakeid"
	"github.com/royhairul/live-studio-api/internal/pkg/validators"
	"github.com/royhairul/live-studio-api/routes"
	"go.uber.org/fx"
)

func Start(router *gin.Engine, cfg *config.Config) {
	router.Run(fmt.Sprintf(":%s", cfg.ServerPort))
}

func NewApp() *fx.App {
	return fx.New(
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
			fx.Annotate(shopee.NewShopeeAffiliateClient, fx.ResultTags(`name:"affiliateShopeeClient"`)),
			fx.Annotate(shopee.NewShopeeDefaultClient, fx.ResultTags(`name:"defaultShopeeClient"`)),
		),
		fx.Invoke(
			database.MigrateDatabase,
			snowflakeid.InitSnowflake,
			Start, // server.Start ini tetap untuk start HTTP server
		),

		// Aggregator
		aggregator.Module,

		// Module Domains
		permission.Module,
		role.Module,
		auth.Module,
		user.Module,
		host.Module,
		studio.Module,
		shift.Module,
		schedule.Module,
		attendance.Module,
		account.Module,
		accountsession.Module,
		accountads.Module,
		performa.Module,
		live.Module,
		finance.Module,
		product.Module,
		order.Module,
		transaction.Module,
		target.Module,
		dashboard.Module,
	)
}

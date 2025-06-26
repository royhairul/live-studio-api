package finance

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee"
	"github.com/royhairul/live-studio-api/internal/domains/finance/controller"
	"github.com/royhairul/live-studio-api/internal/domains/finance/service"

	AccountRepo "github.com/royhairul/live-studio-api/internal/domains/account/repository"
	ShopeeRepo "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/repository"
	ShopeeService "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/service"
)

func RegisterRouter(router *gin.RouterGroup) {
	shopeeClient := shopee.NewShopeeClient("https://creator.shopee.co.id", 10*time.Second)
	shopeeFinanceRepo := ShopeeRepo.NewShopeeFinanceRepository(shopeeClient)
	shopeeFinanceService := ShopeeService.NewShopeeFinanceService(shopeeFinanceRepo)

	accountRepo := AccountRepo.NewAccountRepository(database.DB)

	financeService := service.NewFinanceService(shopeeFinanceService, accountRepo)
	financeController := controller.NewFinanceController(financeService)

	financeRouter := router.Group("/finance")
	{
		financeRouter.POST("/shopee", financeController.GetLiveFinance)
	}
}

package account

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/account/controller"
	"github.com/royhairul/live-studio-api/internal/domains/account/repository"
	"github.com/royhairul/live-studio-api/internal/domains/account/service"
	"github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee"
	ShopeeRepo "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/repository"
	ShopeeService "github.com/royhairul/live-studio-api/internal/domains/ecommerce/shopee/service"
)

func RegisterRouter(router *gin.RouterGroup) {
	// Shopee
	shopeeClient := shopee.NewShopeeClient("https://shopee.co.id", 10*time.Second)
	shopeeAccountRepo := ShopeeRepo.NewShopeeAccountRepository(shopeeClient)
	shopeeService := ShopeeService.NewAccountShopeeService(shopeeAccountRepo)

	accountRepo := repository.NewAccountRepository(database.DB)
	accountService := service.NewAccountService(accountRepo, shopeeService)

	// 4. Init Account Controller → inject AccountService
	accountController := controller.NewAccountController(accountService)

	// 5. Setup Route
	accountRoutes := router.Group("/account")
	{
		accountRoutes.GET("/shopee", accountController.FindAll)
		accountRoutes.POST("/shopee", accountController.CreateOrUpdate)
		accountRoutes.GET("/shopee/:id", accountController.FindById)
		accountRoutes.DELETE("/shopee/:id", accountController.Delete)
	}
}

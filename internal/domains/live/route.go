package live

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/domains/live/controller"
	"github.com/royhairul/live-studio-api/internal/domains/live/service"

	ShopeeRepo "github.com/royhairul/live-studio-api/internal/clients/shopee/repository"
	ShopeeService "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	AccountRepo "github.com/royhairul/live-studio-api/internal/domains/account/repository"
)

func RegisterRouter(router *gin.RouterGroup) {
	shopeeClient := shopee.NewShopeeClient("https://creator.shopee.co.id", 10*time.Second)
	shopeeLiveRepo := ShopeeRepo.NewShopeeLiveRepository(shopeeClient)
	shopeeLiveService := ShopeeService.NewShopeeLiveService(shopeeLiveRepo)

	accountRepo := AccountRepo.NewAccountRepository(database.DB)

	liveService := service.NewLiveService(accountRepo, shopeeLiveService)
	liveController := controller.NewLiveController(liveService)

	liveRouter := router.Group("/live")
	{
		liveRouter.GET("/shopee", liveController.GetLive)
	}
}

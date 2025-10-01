package finance

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/finance/controller"
)

func RegisterRouter(router *gin.RouterGroup, controller controller.FinanceController) {
	routes := router.Group("/finance")
	{
		routes.POST("/shopee", controller.GetLiveFinance)

		routes.GET("/commission", controller.GetLiveFinance)
	}
}

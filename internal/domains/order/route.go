package order

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/order/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.OrderController) {
	// TODO: define routes
	orderRouter := router.Group("/order")
	{
		orderRouter.GET("", controller.FindAll)
		orderRouter.POST("", controller.Create)
		orderRouter.GET("/:id", controller.FindByID)
		orderRouter.PUT("/:id", controller.Update)
		orderRouter.DELETE("/:id", controller.Delete)
	}
}

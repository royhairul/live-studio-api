package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.TransactionController) {
	// TODO: define routes
	routes := router.Group("/transaction")
	{
		routes.GET("", controller.FindAll)
		routes.POST("", controller.Create)
		routes.GET("/grouped", controller.FindAllGrouped)
		routes.GET("/:id", controller.FindByID)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)
	}
}

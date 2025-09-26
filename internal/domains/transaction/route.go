package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.TransactionController) {
	// TODO: define routes
	transactionRouter := router.Group("/transaction")
	{
		transactionRouter.GET("", controller.FindAll)
		transactionRouter.POST("", controller.Create)
		// transactionRouter.GET("/:id", controller.FindByID)
		transactionRouter.PUT("/:id", controller.Update)
		transactionRouter.DELETE("/:id", controller.Delete)
	}
}

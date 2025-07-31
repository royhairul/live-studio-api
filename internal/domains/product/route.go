package product

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/product/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.ProductController) {
	// TODO: define routes
	productRouter := router.Group("/product")
	{
		productRouter.GET("", controller.FindAll)
		productRouter.POST("", controller.Create)
		productRouter.GET("/:id", controller.FindByID)
		productRouter.PUT("/:id", controller.Update)
		productRouter.DELETE("/:id", controller.Delete)
	}
}

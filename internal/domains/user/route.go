package user

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/user/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.UserController) {
	routes := router.Group("/users")
	{
		routes.GET("", controller.GetAll)
		routes.POST("", controller.Create)
		routes.GET("/:id", controller.GetByID)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)

	}
}

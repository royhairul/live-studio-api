package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/permission/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.PermissionController) {
	// TODO: define routes
	routes := router.Group("/permission")
	{
		routes.GET("", controller.FindAll)
		routes.POST("", controller.Create)
		routes.GET("/:id", controller.FindByID)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)
	}
}

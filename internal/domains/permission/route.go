package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/permission/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.PermissionController) {
	// TODO: define routes
	permissionRouter := router.Group("/permission")
	{
		permissionRouter.GET("", controller.FindAll)
		permissionRouter.POST("", controller.Create)
		permissionRouter.GET("/:id", controller.FindByID)
		permissionRouter.PUT("/:id", controller.Update)
		permissionRouter.DELETE("/:id", controller.Delete)
	}
}

package role

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/role/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.RoleController) {
	// TODO: define routes
	roleRouter := router.Group("/role")
	{
		roleRouter.GET("", controller.FindAll)
		roleRouter.POST("", controller.Create)
		roleRouter.GET("/:id", controller.FindByID)
		roleRouter.PUT("/:id", controller.Update)
		roleRouter.DELETE("/:id", controller.Delete)
	}
}

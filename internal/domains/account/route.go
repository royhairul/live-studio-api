package account

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/account/controller"
)

func RegisterRouter(router *gin.RouterGroup, controller controller.AccountController) {
	routes := router.Group("/account")
	{
		routes.GET("", controller.FindAll)
		routes.POST("/", controller.CreateOrUpdate)
		routes.GET("/:id", controller.FindById)
		routes.DELETE("/:id", controller.Delete)
	}
}

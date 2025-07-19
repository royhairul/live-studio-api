package account

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/account/controller"
)

func RegisterRouter(router *gin.RouterGroup, controller controller.AccountController) {
	routes := router.Group("/account")
	{
		routes.GET("/shopee", controller.FindAll)
		routes.POST("/shopee", controller.CreateOrUpdate)
		routes.GET("/shopee/:id", controller.FindById)
		routes.DELETE("/shopee/:id", controller.Delete)
	}
}

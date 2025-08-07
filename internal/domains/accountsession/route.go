package accountsession

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.AccountsessionController) {
	// TODO: define routes
	accountsessionRouter := router.Group("/account-session")
	{
		accountsessionRouter.GET("", controller.FindAll)
		accountsessionRouter.POST("", controller.Create)
		accountsessionRouter.GET("/:id", controller.FindByID)
		accountsessionRouter.PUT("/:id", controller.Update)
		accountsessionRouter.DELETE("/:id", controller.Delete)
	}
}

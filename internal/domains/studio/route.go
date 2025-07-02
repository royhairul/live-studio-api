package studio

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/studio/controller"
	"github.com/royhairul/live-studio-api/internal/domains/studio/repository"
	"github.com/royhairul/live-studio-api/internal/domains/studio/service"
)

func RegisterRoutes(router *gin.RouterGroup) {
	// TODO: define routes
	validate := validator.New()
	studioRepo := repository.NewStudioRepository(database.DB)
	studioService := service.NewStudioService(studioRepo)
	studioController := controller.NewStudioController(studioService, validate)

	studioRouter := router.Group("/studio")
	{
		studioRouter.GET("", studioController.FindAll)
		studioRouter.POST("", studioController.Create)
		studioRouter.GET("/:id", studioController.FindByID)
		studioRouter.PUT("/:id", studioController.Update)
		studioRouter.DELETE("/:id", studioController.Delete)
	}
}

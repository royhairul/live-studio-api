package host

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/host/controller"
	"github.com/royhairul/live-studio-api/internal/domains/host/repository"
	"github.com/royhairul/live-studio-api/internal/domains/host/service"
)

func RegisterRouter(router *gin.RouterGroup) {
	validate := validator.New()
	hostRepo := repository.NewHostRepository(database.DB)
	hostService := service.NewHostService(hostRepo)
	hostController := controller.NewHostController(hostService, validate)

	hostRoutes := router.Group("/host")
	{
		hostRoutes.GET("", hostController.FindAll)
		hostRoutes.POST("", hostController.Create)
		hostRoutes.GET("/:id", hostController.FindByID)
		hostRoutes.PUT("/:id", hostController.Update)
		hostRoutes.DELETE("/:id", hostController.Delete)

		hostRoutes.GET("/group-by-studio", hostController.FindAllGroupedByStudio)
	}
}

package shift

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/shift/controller"
	"github.com/royhairul/live-studio-api/internal/domains/shift/repository"
	"github.com/royhairul/live-studio-api/internal/domains/shift/service"
)

func RegisterRoutes(router *gin.RouterGroup) {
	shiftRepo := repository.NewShiftRepository(database.DB)
	shiftService := service.NewShiftService(shiftRepo)
	shiftController := controller.NewShiftController(shiftService)

	// TODO: define routes
	shiftRouter := router.Group("/shift")
	{
		shiftRouter.GET("", shiftController.FindAll)
		shiftRouter.POST("", shiftController.Create)
		shiftRouter.GET("/:id", shiftController.FindByID)
		shiftRouter.PUT("/:id", shiftController.Update)
		shiftRouter.DELETE("/:id", shiftController.Delete)
	}
}

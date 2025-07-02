package attendance

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/controller"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/repository"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/service"
	schedulerepository "github.com/royhairul/live-studio-api/internal/domains/schedule/repository"
)

func RegisterRouter(router *gin.RouterGroup) {

	attendanceRepo := repository.NewAttendanceRepository(database.DB)
	scheduleRepo := schedulerepository.NewScheduleRepository(database.DB)
	attendanceService := service.NewAttendanceService(attendanceRepo, scheduleRepo)
	attendanceController := controller.NewAttendanceController(attendanceService)

	// 5. Setup Route
	attendanceRoutes := router.Group("/attendance")
	{
		attendanceRoutes.GET("", attendanceController.FindAll)
		attendanceRoutes.GET("/unchecked-out", attendanceController.FindUncheckedOut)
		attendanceRoutes.POST("/check-in", attendanceController.CheckIn)
		attendanceRoutes.POST("/check-out", attendanceController.CheckOut)
	}
}

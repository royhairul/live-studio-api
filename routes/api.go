package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/controllers"
	"github.com/royhairul/live-studio-api/internal/domains/account"
	"github.com/royhairul/live-studio-api/internal/domains/attendance"
	"github.com/royhairul/live-studio-api/internal/domains/finance"
	"github.com/royhairul/live-studio-api/internal/domains/host"
	"github.com/royhairul/live-studio-api/internal/domains/live"
	"github.com/royhairul/live-studio-api/internal/domains/schedule"
	"github.com/royhairul/live-studio-api/internal/domains/shift"
	"github.com/royhairul/live-studio-api/internal/domains/studio"
	"github.com/royhairul/live-studio-api/middleware"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// User Routes
	route.POST("/api/login", controllers.Login)
	route.POST("/api/register", controllers.Register)
	route.POST("/api/forgot-password", controllers.ForgotPassword)
	route.POST("/api/verify-otp", controllers.VerifyOtp)
	route.POST("/api/reset-password", controllers.ResetPassword)

	route.GET("/api/me", middleware.RequireRoles("superadmin", "admin", "host"), controllers.Me)

	// User Management Routes
	RegisterSuperAdminRoutes(route)

	// Product Routes
	route.GET("/api/products/all", controllers.ProductIndex)
	route.POST("/api/products/create", controllers.ProductCreate)

	// Akun Routes
	route.GET("/api/account", controllers.AccountIndex)
	route.POST("/api/account", controllers.AccountCreate)

	// Host Routes
	// route.GET("/api/host", controllers.HostIndex)
	// route.POST("/api/host", controllers.HostCreate)
	// route.GET("/api/host/:id", controllers.HostShow)
	// route.PUT("/api/host/:id", controllers.HostUpdate)
	// route.DELETE("/api/host/:id", controllers.HostDelete)
	route.GET("/api/host/:id/schedule", controllers.HostScheduleByHostID)

	// Host by Studio Routes
	// route.GET("/api/host/group-by-studio", controllers.HostGroupedByStudio)

	// Host Schedule Routes
	// route.GET("/api/host-schedule", controllers.HostScheduleIndex)
	// route.POST("/api/host-schedule", controllers.HostScheduleCreate)
	// route.GET("/api/host-schedule/:id", controllers.HostScheduleShow)
	// route.PUT("/api/host-schedule/:id", controllers.HostScheduleUpdate)
	// route.DELETE("/api/host-schedule/:id", controllers.HostScheduleDelete)
	// route.GET("/api/host-schedule/scheduled", controllers.HostScheduleScheduled)
	// route.POST("/api/host-schedule/switch", controllers.HostScheduleSwitch)

	// Shift Route
	// route.GET("/api/shift", controllers.ScheduleShiftIndex)
	// route.POST("/api/shift", controllers.ScheduleShiftCreate)
	// route.GET("/api/shift/:id", controllers.ScheduleShiftShow)
	// route.PUT("/api/shift/:id", controllers.ScheduleShiftUpdate)
	// route.DELETE("/api/shift/:id", controllers.ScheduleShiftDelete)

	// Studio Routes
	// route.GET("/api/studio", controllers.StudioIndex)
	// route.POST("/api/studio", controllers.StudioCreate)
	// route.GET("/api/studio/:id", controllers.StudioShow)
	// route.PUT("/api/studio/:id", controllers.StudioUpdate)
	// route.DELETE("/api/studio/:id", controllers.StudioDelete)

	// Roles
	route.GET("api/role", controllers.RoleIndex)
	route.POST("api/role", controllers.RoleCreate)
	route.GET("api/role/:id", controllers.RoleShow)
	route.PUT("api/role/:id", controllers.RoleUpdate)
	route.DELETE("api/role/:id", controllers.RoleDelete)

	// Permissions
	route.GET("api/permission", controllers.PermissionIndex)
	route.GET("api/permission/grouped", controllers.PermissionIndexGrouped)

	// Finance
	route.GET("api/finance/report", controllers.FinanceGetLiveReport)

	shift.RegisterRoutes(route.Group("/api"))
	studio.RegisterRoutes(route.Group("/api"))
	host.RegisterRouter(route.Group("/api"))
	schedule.RegisterRoutes(route.Group("/api"))
	account.RegisterRouter(route.Group("/api"))
	finance.RegisterRouter(route.Group("/api"))
	attendance.RegisterRouter(route.Group("/api"))
	live.RegisterRouter(route.Group("/api"))

	return route
}

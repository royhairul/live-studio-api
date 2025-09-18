package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/royhairul/live-studio-api/controllers"
// 	"github.com/royhairul/live-studio-api/internal/middleware"
// )

// func RegisterSuperAdminRoutes(r *gin.Engine) {
// 	superadmin := r.Group("api/superadmin")
// 	superadmin.Use(middleware.RequireRoles("superadmin"))
// 	{
// 		superadmin.GET("/user", controllers.UserIndex)
// 		superadmin.POST("/user", controllers.UserCreate)
// 		superadmin.GET("/user/:id", controllers.UserShow)
// 		superadmin.PUT("user/:id", controllers.UserUpdate)
// 		superadmin.DELETE("/user/:id", controllers.UserDelete)
// 	}
// }

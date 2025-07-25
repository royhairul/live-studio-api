package {{feature}}

import (
	"github.com/gin-gonic/gin"
	"{{Module}}/internal/domains/{{feature}}/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.{{Feature}}Controller) {
	// TODO: define routes
	{{feature}}Router := router.Group("/{{feature}}")
	{
		{{feature}}Router.GET("", controller.FindAll)
		{{feature}}Router.POST("", controller.Create)
		{{feature}}Router.GET("/:id", controller.FindByID)
		{{feature}}Router.PUT("/:id", controller.Update)
		{{feature}}Router.DELETE("/:id", controller.Delete)
	}
}

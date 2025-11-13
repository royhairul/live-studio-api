package controller

import "github.com/gin-gonic/gin"

type AccountController interface {
	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	CreateOrUpdate(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)

	FindByStudio(ctx *gin.Context)
}

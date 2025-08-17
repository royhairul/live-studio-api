package controller

import "github.com/gin-gonic/gin"

type HostController interface {
	FindAll(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)

	FindAllGroupedByStudio(ctx *gin.Context)
}

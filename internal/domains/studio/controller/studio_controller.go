package controller

import "github.com/gin-gonic/gin"

type StudioController interface {
	// TODO: define controller methods
	FindAll(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

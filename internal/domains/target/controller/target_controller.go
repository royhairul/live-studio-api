package controller

import "github.com/gin-gonic/gin"

type TargetController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
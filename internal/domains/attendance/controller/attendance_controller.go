package controller

import "github.com/gin-gonic/gin"

type AttendanceController interface {
	FindAll(ctx *gin.Context)
	FindUncheckedOut(ctx *gin.Context)
	CheckIn(ctx *gin.Context)
	CheckOut(ctx *gin.Context)
}

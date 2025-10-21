package controller

import "github.com/gin-gonic/gin"

type FinanceController interface {
	FindAll(ctx *gin.Context)
}

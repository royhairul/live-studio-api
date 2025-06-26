package controller

import "github.com/gin-gonic/gin"

type FinanceController interface {
	GetLiveFinance(ctx *gin.Context)
}

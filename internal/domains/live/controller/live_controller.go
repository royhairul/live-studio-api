package controller

import "github.com/gin-gonic/gin"

type LiveController interface {
	GetLive(ctx *gin.Context)
	GetLiveDetail(ctx *gin.Context)
}

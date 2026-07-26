package controller

import "github.com/gin-gonic/gin"

type LiveController interface {
	GetLive(ctx *gin.Context)
	GetLiveDetail(ctx *gin.Context)
	GetStoredHistory(ctx *gin.Context)
	GetStoredHistoryDetail(ctx *gin.Context)
	SyncHistory(ctx *gin.Context)
	SyncAllHistory(ctx *gin.Context)
	SyncHistoryRange(ctx *gin.Context)
}

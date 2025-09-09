package controller

import "github.com/gin-gonic/gin"

type PerformaController interface {
	GetHosts(ctx *gin.Context)
	GetHostByID(ctx *gin.Context)
	GetAccounts(ctx *gin.Context)
	GetStudios(ctx *gin.Context)
	GetStudioByID(ctx *gin.Context)
}

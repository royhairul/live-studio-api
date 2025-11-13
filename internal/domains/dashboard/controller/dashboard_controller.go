package controller

import "github.com/gin-gonic/gin"

type DashboardController interface {
	// TODO: define controller methods
	Dashboard(ctx *gin.Context)
}

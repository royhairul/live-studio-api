package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type DashboardControllerImpl struct {
	// TODO: add dependencies
	service service.DashboardService
}

func NewDashboardController(service service.DashboardService) DashboardController {
	return &DashboardControllerImpl{service}
}

// Dashboard implements DashboardController.
func (d *DashboardControllerImpl) Dashboard(ctx *gin.Context) {
	startDate, _ := ctx.GetQuery("startDate")
	endDate, _ := ctx.GetQuery("endDate")

	result, err := d.service.DashboardAdmin(ctx.Request.Context(), startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("Dashboard", result)
	ctx.JSON(200, resp)
}

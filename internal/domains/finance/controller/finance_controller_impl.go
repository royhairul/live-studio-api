package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/finance/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
)

type FinanceControllerImpl struct {
	service service.FinanceService
}

func NewFinanceController(financeSvc service.FinanceService) FinanceController {
	return &FinanceControllerImpl{service: financeSvc}
}

func (f *FinanceControllerImpl) FindAll(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	start, end, _ := timehandler.ParseDateRange(startDate, endDate)

	account := ctx.Query("account")
	studio := ctx.Query("studio")
	status := ctx.Query("status")
	payment := ctx.Query("payment")

	service := f.service

	if studio != "" {
		service = service.WithStudioID(studio)
	}

	if account != "" {
		service = service.WithAccountUniqueID(account)
	}

	if status != "" {
		service = service.WithStatus(status)
	}

	if payment != "" {
		service = service.WithPaymentMethod(payment)
	}

	finance, err := service.FindAll(ctx.Request.Context(), start, end)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved data successfully", finance)
	ctx.JSON(http.StatusOK, resp)
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/helpers/response"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/service"
)

type AttendanceControllerImpl struct {
	service service.AttendanceService
}

func NewAttendanceController(service service.AttendanceService) AttendanceController {
	return &AttendanceControllerImpl{service}
}

// FindAll implements AttendanceController.
func (c *AttendanceControllerImpl) FindAll(ctx *gin.Context) {
	attendances, err := c.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	if len(attendances) <= 0 {
		resp := response.NewBaseResponse("empty data attendance", attendances)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp := response.NewBaseResponse("retrieved all attendance successfully", attendances)
	ctx.JSON(http.StatusOK, resp)
}

// FindUncheckedOut implements AttendanceController.
func (c *AttendanceControllerImpl) FindUncheckedOut(ctx *gin.Context) {
	attendances, err := c.service.FindUncheckedOut()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	if len(attendances) <= 0 {
		resp := response.NewBaseResponse("empty data attendance", attendances)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp := response.NewBaseResponse("retrieved all attendance successfully", attendances)
	ctx.JSON(http.StatusOK, resp)
}

// CheckIn implements AttendanceController.
func (c *AttendanceControllerImpl) CheckIn(ctx *gin.Context) {
	var req params.AttendanceCheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	results, err := c.service.CheckIn(req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("attendance checkin successfully", results)
	ctx.JSON(http.StatusOK, resp)
}

// CheckOut implements AttendanceController.
func (c *AttendanceControllerImpl) CheckOut(ctx *gin.Context) {
	var req params.AttendanceCheckOutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	err := c.service.CheckOut(req)
	if err != nil {
		errorhandler.NewBadRequestError("error for checkout", err)
		return
	}

	resp := response.NewBaseResponse("attendance checkout successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

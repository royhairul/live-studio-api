package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type AttendanceControllerImpl struct {
	service  service.AttendanceService
	validate *validator.Validate
}

func NewAttendanceController(service service.AttendanceService, validate *validator.Validate) AttendanceController {
	return &AttendanceControllerImpl{service, validate}
}

// FindAll implements AttendanceController.
func (c *AttendanceControllerImpl) FindAll(ctx *gin.Context) {
	attendances, err := c.service.FindAll(ctx.Request.Context())
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
	attendances, err := c.service.FindUncheckedOut(ctx.Request.Context())
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

	results, err := c.service.CheckIn(ctx.Request.Context(), req)
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

	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.CheckOut(ctx.Request.Context(), req)
	if err != nil {
		errorhandler.NewBadRequestError("error for checkout", err)
		return
	}

	resp := response.NewBaseResponse("attendance checkout successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/shift/params"
	"github.com/royhairul/live-studio-api/internal/domains/shift/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/paramhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type ShiftControllerImpl struct {
	// TODO: add dependencies
	service service.ShiftService
}

func NewShiftController(service service.ShiftService) ShiftController {
	return &ShiftControllerImpl{service}
}

// Create implements ShiftController.
func (s *ShiftControllerImpl) Create(ctx *gin.Context) {
	var shiftReq params.CreateShiftRequest
	if err := ctx.ShouldBindJSON(&shiftReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	createdShift, err := s.service.Create(&shiftReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("shift created successfully", createdShift)

	ctx.JSON(http.StatusCreated, resp)
}

// Delete implements ShiftController.
func (s *ShiftControllerImpl) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := paramhandler.ParseUintParam(idParam)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	if err := s.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("shift deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// FindAll implements ShiftController.
func (s *ShiftControllerImpl) FindAll(ctx *gin.Context) {
	shifts, err := s.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("shifts retrieved successfully", shifts)
	ctx.JSON(http.StatusOK, resp)
}

// FindByID implements ShiftController.
func (s *ShiftControllerImpl) FindByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := paramhandler.ParseUintParam(idParam)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	shift, err := s.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("shift retrieved successfully", shift)
	ctx.JSON(http.StatusOK, resp)
}

// Update implements ShiftController.
func (s *ShiftControllerImpl) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := paramhandler.ParseUintParam(idParam)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	var shiftReq params.UpdateShiftRequest
	if err := ctx.ShouldBindJSON(&shiftReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	updatedShift, err := s.service.Update(id, &shiftReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("shift updated successfully", updatedShift)
	ctx.JSON(http.StatusOK, resp)
}

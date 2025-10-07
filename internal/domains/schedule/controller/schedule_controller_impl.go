package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/params"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/paramhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type ScheduleControllerImpl struct {
	service  service.ScheduleService
	validate *validator.Validate
}

func NewScheduleController(service service.ScheduleService, validate *validator.Validate) ScheduleController {
	return &ScheduleControllerImpl{service, validate}
}

// Create implements ScheduleController.
func (s *ScheduleControllerImpl) Create(ctx *gin.Context) {
	var scheduleReq params.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&scheduleReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	if err := s.validate.Struct(&scheduleReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	schedule, err := s.service.Create(&scheduleReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("schedule created successfully", schedule)
	ctx.JSON(http.StatusCreated, resp)
}

// Delete implements ScheduleController.
func (s *ScheduleControllerImpl) Delete(ctx *gin.Context) {
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

	resp := response.NewBaseResponse("schedule deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// FindAll implements ScheduleController.
func (s *ScheduleControllerImpl) FindAll(ctx *gin.Context) {
	schedules, err := s.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all data schedule successfully", schedules)
	ctx.JSON(http.StatusOK, resp)
}

// FindByID implements ScheduleController.
func (s *ScheduleControllerImpl) FindByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := paramhandler.ParseUintParam(idParam)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	schedule, err := s.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved schedule successfully", schedule)
	ctx.JSON(http.StatusOK, resp)
}

// Update implements ScheduleController.
func (s *ScheduleControllerImpl) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := paramhandler.ParseUintParam(idParam)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	var scheduleReq params.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&scheduleReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	if err := s.validate.Struct(&scheduleReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	schedule, err := s.service.Update(id, &scheduleReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("schedule updated successfully", schedule)
	ctx.JSON(http.StatusOK, resp)
}

// FindByHostShiftAndDate implements ScheduleController.
func (s *ScheduleControllerImpl) FindByShiftAndDate(ctx *gin.Context) {
	date := ctx.Query("date")
	shiftID := ctx.Query("shift_id")

	schedules, err := s.service.FindByShiftAndDate(shiftID, date)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}
	resp := response.NewBaseResponse("retrieved schedules by host, shift, and date successfully", schedules)
	ctx.JSON(http.StatusOK, resp)
}

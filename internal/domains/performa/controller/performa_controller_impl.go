package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/performa/service"
	"github.com/royhairul/live-studio-api/internal/pkg/constants"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type PerformaControllerImpl struct {
	service  service.PerformaService
	validate *validator.Validate
}

func NewPerformaController(service service.PerformaService, validate *validator.Validate) PerformaController {
	return &PerformaControllerImpl{service, validate}
}

// GetHosts implements PerformaController.
func (p *PerformaControllerImpl) GetHosts(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	result, err := p.service.GetHosts(ctx.Request.Context(), startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all performa host successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetHostByID implements PerformaController.
func (p *PerformaControllerImpl) GetHostByID(ctx *gin.Context) {
	id := ctx.Param("id")

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	// Jika salah satu atau keduanya kosong, isi dengan hari ini
	if startDate == "" || endDate == "" {
		today := time.Now().Format(constants.LayoutYYMMDD)
		startDate = today
		endDate = today
	}

	result, err := p.service.GetHostByID(ctx.Request.Context(), id, startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("performa host found", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetAccounts implements PerformaController.
func (p *PerformaControllerImpl) GetAccounts(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	result, err := p.service.GetAccounts(ctx.Request.Context(), startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all performa account successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetStudios implements PerformaController.
func (p *PerformaControllerImpl) GetStudios(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	result, err := p.service.GetStudios(ctx.Request.Context(), startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all performa studio successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetStudioByID implements PerformaController.
func (p *PerformaControllerImpl) GetStudioByID(ctx *gin.Context) {
	id := ctx.Param("id")

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	result, err := p.service.GetStudioByID(ctx.Request.Context(), id, startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("performa studio found", result)
	ctx.JSON(http.StatusOK, resp)
}

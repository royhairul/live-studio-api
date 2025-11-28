package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/internal/domains/studio/params"
	"github.com/royhairul/live-studio-api/internal/domains/studio/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type StudioControllerImpl struct {
	// TODO: add dependencies
	service  service.StudioService
	validate *validator.Validate
}

func NewStudioController(service service.StudioService, validate *validator.Validate) StudioController {
	return &StudioControllerImpl{service, validate}
}

// Create implements StudioController.
func (s *StudioControllerImpl) Create(ctx *gin.Context) {
	var studioReq params.CreateStudioRequest
	if err := ctx.ShouldBindJSON(&studioReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	if err := s.validate.Struct(studioReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	studio, err := s.service.Create(ctx.Request.Context(), studioReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("studio created successfully", studio)
	ctx.JSON(http.StatusOK, resp)
}

// Delete implements StudioController.
func (s *StudioControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := s.service.Delete(ctx.Request.Context(), id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("studio deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// FindAll implements StudioController.
func (s *StudioControllerImpl) FindAll(ctx *gin.Context) {
	studios, err := s.service.FindAll(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all studio successfully", studios)
	ctx.JSON(http.StatusOK, resp)
}

// FindByID implements StudioController.
func (s *StudioControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	studio, err := s.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("studio has found", studio)
	ctx.JSON(http.StatusOK, resp)
}

// Update implements StudioController.
func (s *StudioControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var studioReq params.UpdateStudioRequest
	if err := ctx.ShouldBindJSON(&studioReq); err != nil {
		errorhandler.HandleError(ctx, err)
	}

	studio, err := s.service.Update(ctx.Request.Context(), id, studioReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("studio updated successfully", studio)
	ctx.JSON(http.StatusOK, resp)
}

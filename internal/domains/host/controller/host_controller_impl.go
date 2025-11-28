package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/internal/domains/host/params"
	"github.com/royhairul/live-studio-api/internal/domains/host/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type HostControllerImpl struct {
	service  service.HostService
	validate *validator.Validate
}

func NewHostController(service service.HostService, validate *validator.Validate) HostController {
	return &HostControllerImpl{service, validate}
}

// Create implements HostController.
func (h *HostControllerImpl) Create(ctx *gin.Context) {
	var hostReq params.CreateHostRequest
	if err := ctx.ShouldBindJSON(&hostReq); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := h.validate.Struct(hostReq); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	createdHost, err := h.service.Create(ctx.Request.Context(), hostReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("created host successfully", createdHost)
	ctx.JSON(http.StatusOK, resp)
}

// Update implements HostController.
func (h *HostControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var hostReq params.UpdateHostRequest
	if err := ctx.ShouldBindJSON(&hostReq); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	createdHost, err := h.service.Update(ctx.Request.Context(), id, hostReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("updated host successfully", createdHost)
	ctx.JSON(http.StatusOK, resp)
}

// FindAll implements HostController.
func (h *HostControllerImpl) FindAll(ctx *gin.Context) {
	hosts, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all host successfully", hosts)
	ctx.JSON(http.StatusOK, resp)
}

// FindByID implements HostController.
func (h *HostControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	host, err := h.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("host found", host)
	ctx.JSON(http.StatusOK, resp)
}

// Delete implements HostController.
func (h *HostControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.service.Delete(ctx.Request.Context(), id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("host deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// FindAllGroupedByStudio implements HostController.
func (h *HostControllerImpl) FindAllGroupedByStudio(ctx *gin.Context) {
	hosts, err := h.service.FindAllGroupedByStudio(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all hosts grouped by studio successfully", hosts)
	ctx.JSON(http.StatusOK, resp)
}

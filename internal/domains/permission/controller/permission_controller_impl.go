package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/permission/params"
	"github.com/royhairul/live-studio-api/internal/domains/permission/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type PermissionControllerImpl struct {
	service  service.PermissionService
	validate *validator.Validate
}

func NewPermissionController(service service.PermissionService, validate *validator.Validate) PermissionController {
	return &PermissionControllerImpl{service, validate}
}

func (c *PermissionControllerImpl) Create(ctx *gin.Context) {
	var req params.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.Create(req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("created permission successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *PermissionControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdatePermissionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	// Validasi hanya field yang tidak nil (optional)
	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.Update(id, req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("updated permission successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *PermissionControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all permission successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// FindAllGrouped implements PermissionController.
func (c *PermissionControllerImpl) FindAllGrouped(ctx *gin.Context) {
	result, err := c.service.FindAllGrouped()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all grouped permission successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *PermissionControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("permission found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *PermissionControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted permission successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

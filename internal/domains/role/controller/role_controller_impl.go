package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/role/params"
	"github.com/royhairul/live-studio-api/internal/domains/role/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type RoleControllerImpl struct {
	service  service.RoleService
	validate *validator.Validate
}

func NewRoleController(service service.RoleService, validate *validator.Validate) RoleController {
	return &RoleControllerImpl{service, validate}
}

func (c *RoleControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateRoleRequest
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

	resp := response.NewBaseResponse("created role successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *RoleControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdateRoleRequest

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

	resp := response.NewBaseResponse("updated role successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *RoleControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all role successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *RoleControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("role found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *RoleControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted role successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

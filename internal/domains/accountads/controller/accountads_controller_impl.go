package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/helpers/response"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/params"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/service"
)

type AccountadsControllerImpl struct {
	service  service.AccountadsService
	validate *validator.Validate
}

func NewAccountadsController(service service.AccountadsService, validate *validator.Validate) AccountadsController {
	return &AccountadsControllerImpl{service, validate}
}

func (c *AccountadsControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateAccountadsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.CreateOrUpdate(req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("created accountads successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *AccountadsControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdateAccountadsRequest

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

	resp := response.NewBaseResponse("updated accountads successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountadsControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all accountads successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountadsControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindOne(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("accountads found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountadsControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted accountads successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type AccountsessionControllerImpl struct {
	service  service.AccountsessionService
	validate *validator.Validate
}

func NewAccountsessionController(service service.AccountsessionService, validate *validator.Validate) AccountsessionController {
	return &AccountsessionControllerImpl{service, validate}
}

func (c *AccountsessionControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateAccountsessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("created accountsession successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *AccountsessionControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdateEndSessionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	// Validasi hanya field yang tidak nil (optional)
	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("updated accountsession successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountsessionControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.service.FindAll(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all accountsession successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountsessionControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.WithID(id).FindOne(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("accountsession found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountsessionControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted accountsession successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

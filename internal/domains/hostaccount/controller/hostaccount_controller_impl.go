package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/params"
	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type HostAccountControllerImpl struct {
	service  service.HostAccountService
	validate *validator.Validate
}

func NewHostAccountController(service service.HostAccountService, validate *validator.Validate) HostAccountController {
	return &HostAccountControllerImpl{service, validate}
}

// FindAll lists assignments, optionally narrowed by host, account, or studio.
// ?active=true keeps only assignments that have not been ended.
func (c *HostAccountControllerImpl) FindAll(ctx *gin.Context) {
	svc := c.service

	if hostID := ctx.Query("host"); hostID != "" {
		svc = svc.WithHostID(hostID)
	}
	if accountID := ctx.Query("account"); accountID != "" {
		svc = svc.WithAccountID(accountID)
	}
	if studioID := ctx.Query("studio"); studioID != "" {
		svc = svc.WithStudioID(studioID)
	}
	if ctx.Query("active") == "true" {
		svc = svc.WithActiveOnly()
	}

	result, err := svc.FindAll(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewBaseResponse("retrieved host accounts successfully", result))
}

func (c *HostAccountControllerImpl) FindByID(ctx *gin.Context) {
	result, err := c.service.WithID(ctx.Param("id")).FindOne(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewBaseResponse("retrieved host account successfully", result))
}

func (c *HostAccountControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateHostAccountRequest
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

	ctx.JSON(http.StatusCreated, response.NewBaseResponse("created host account successfully", result))
}

func (c *HostAccountControllerImpl) Update(ctx *gin.Context) {
	var req params.UpdateHostAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), ctx.Param("id"), req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewBaseResponse("updated host account successfully", result))
}

func (c *HostAccountControllerImpl) Delete(ctx *gin.Context) {
	if err := c.service.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response.NewBaseResponse("deleted host account successfully", nil))
}

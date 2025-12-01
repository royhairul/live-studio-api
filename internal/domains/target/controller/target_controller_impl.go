package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/target/params"
	"github.com/royhairul/live-studio-api/internal/domains/target/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type TargetControllerImpl struct {
	service  service.TargetService
	validate *validator.Validate
}

func NewTargetController(service service.TargetService, validate *validator.Validate) TargetController {
	return &TargetControllerImpl{service, validate}
}

func (c *TargetControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.CreateOrUpdate(ctx.Request.Context(), req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("created target successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *TargetControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req params.UpdateTargetRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := c.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("updated target successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TargetControllerImpl) FindAll(ctx *gin.Context) {
	month := ctx.Query("month")
	year := ctx.Query("year")
	studio := ctx.Query("studio")

	req := params.TargetRequest{
		Month:  month,
		Year:   year,
		Studio: studio,
	}

	result, err := c.service.FindAll(ctx.Request.Context(), req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved target data successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TargetControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(ctx.Request.Context(), id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("target found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TargetControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted target successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

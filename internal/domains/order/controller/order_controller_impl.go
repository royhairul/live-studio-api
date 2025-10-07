package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/order/params"
	"github.com/royhairul/live-studio-api/internal/domains/order/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type OrderControllerImpl struct {
	service  service.OrderService
	validate *validator.Validate
}

func NewOrderController(service service.OrderService, validate *validator.Validate) OrderController {
	return &OrderControllerImpl{service, validate}
}

func (c *OrderControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateOrderRequest
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

	resp := response.NewBaseResponse("created Order successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *OrderControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdateOrderRequest

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

	resp := response.NewBaseResponse("updated Order successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *OrderControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all Order successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *OrderControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("Order found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *OrderControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted Order successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

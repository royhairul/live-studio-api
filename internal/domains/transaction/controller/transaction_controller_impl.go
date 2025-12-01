package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
)

type TransactionControllerImpl struct {
	service  service.TransactionService
	validate *validator.Validate
}

func NewTransactionController(service service.TransactionService, validate *validator.Validate) TransactionController {
	return &TransactionControllerImpl{service, validate}
}

func (c *TransactionControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateTransactionRequest
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

	resp := response.NewBaseResponse("created transaction successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *TransactionControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdateTransactionRequest

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

	resp := response.NewBaseResponse("updated Transaction successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TransactionControllerImpl) FindAll(ctx *gin.Context) {
	status := ctx.Query("status")
	account := ctx.Query("account")
	studio := ctx.Query("studio")

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	start, end, _ := timehandler.ParseDateRange(startDate, endDate)

	service := c.service

	if status != "" {
		service = service.WithStatus(status)
	}
	if account != "" {
		service = service.WithAccountID(account)
	}
	if startDate != "" || endDate != "" {
		service = service.WithDate(*start, *end)
	}
	if studio != "" {
		service = service.WithStudioID(studio)
	}

	result, err := service.FindAll(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all transaction successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// FindAllGrouped implements TransactionController.
func (c *TransactionControllerImpl) FindAllGrouped(ctx *gin.Context) {
	status := ctx.Query("status")
	account := ctx.Query("account")

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	start, end, _ := timehandler.ParseDateRange(startDate, endDate)

	service := c.service

	if status != "" {
		service = service.WithStatus(status)
	}
	if account != "" {
		service = service.WithAccountID(account)
	}
	if startDate != "" || endDate != "" {
		service = service.WithDate(*start, *end)
	}

	result, err := service.FindAllGrouped(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all transaction successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TransactionControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.WithID(id).FindOne(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("transaction found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TransactionControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted transaction successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

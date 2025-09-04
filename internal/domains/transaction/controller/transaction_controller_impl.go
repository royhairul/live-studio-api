package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/helpers/response"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/service"
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

	result, err := c.service.Create(req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("created Transaction successfully", result)
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

	result, err := c.service.Update(id, req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("updated Transaction successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TransactionControllerImpl) FindAll(ctx *gin.Context) {
	status := ctx.Query("status")

	var (
		result []*params.TransactionResponse
		err    error
	)
	log.Printf("status: %s", status)

	if status != "" {
		log.Println("status")
		result, err = c.service.FindAllByStatus(status)
	} else {
		log.Println("all")
		result, err = c.service.FindAll()
	}

	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all Transaction successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TransactionControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("Transaction found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *TransactionControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted Transaction successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

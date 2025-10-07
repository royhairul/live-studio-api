package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/product/params"
	"github.com/royhairul/live-studio-api/internal/domains/product/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type ProductControllerImpl struct {
	service  service.ProductService
	validate *validator.Validate
}

func NewProductController(service service.ProductService, validate *validator.Validate) ProductController {
	return &ProductControllerImpl{service, validate}
}

func (c *ProductControllerImpl) Create(ctx *gin.Context) {
	var req params.CreateProductRequest
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

	resp := response.NewBaseResponse("created Product successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *ProductControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdateProductRequest

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

	resp := response.NewBaseResponse("updated Product successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ProductControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all Product successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ProductControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("Product found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ProductControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted Product successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

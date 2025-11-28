package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/internal/domains/account/params"
	"github.com/royhairul/live-studio-api/internal/domains/account/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type AccountControllerImpl struct {
	service  service.AccountService
	validate *validator.Validate
}

func NewAccountController(service service.AccountService) AccountController {
	return &AccountControllerImpl{service: service}
}

func (a *AccountControllerImpl) FindAll(ctx *gin.Context) {
	studioId := ctx.Query("studio")
	uniqueId := ctx.Query("uniqueId")

	service := a.service

	if studioId != "" {
		service = service.WithStudioID(studioId)
	}
	if uniqueId != "" {
		service = service.WithUniqueID(uniqueId)
	}

	accounts, err := service.FindAll(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved accounts successfully", accounts)
	ctx.JSON(http.StatusOK, resp)
}

func (a *AccountControllerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	account, err := a.service.WithID(id).FindOne(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved an account successfully", account)
	ctx.JSON(http.StatusOK, resp)
}

func (a *AccountControllerImpl) CreateOrUpdate(ctx *gin.Context) {
	var accountReq params.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&accountReq); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err.Error()))
		return
	}

	account, err := a.service.CreateOrUpdate(ctx.Request.Context(), accountReq)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("account created or updated successfully", account)
	ctx.JSON(http.StatusOK, resp)
}

// Update implements AccountController.
func (a *AccountControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.UpdateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
		return
	}

	if err := a.validate.Struct(req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	result, err := a.service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("updated accountads successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// Delete implements AccountController.
func (a *AccountControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := a.service.Delete(ctx.Request.Context(), id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("account deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// FindByStudio implements AccountController.
func (a *AccountControllerImpl) FindByStudio(ctx *gin.Context) {
	studioId := ctx.Param("studioId")

	accounts, err := a.service.WithStudioID(studioId).FindAll(ctx.Request.Context())
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved accounts by studio successfully", accounts)
	ctx.JSON(http.StatusOK, resp)
}

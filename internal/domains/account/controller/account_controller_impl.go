package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/helpers/response"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
	"github.com/royhairul/live-studio-api/internal/domains/account/service"
)

type AccountControllerImpl struct {
	AccountService service.AccountService
	validate       *validator.Validate
}

func NewAccountController(accountSvc service.AccountService) AccountController {
	return &AccountControllerImpl{AccountService: accountSvc}
}

func (a *AccountControllerImpl) FindAll(ctx *gin.Context) {
	studioId := ctx.Query("studio")

	if studioId != "" {
		accounts, err := a.AccountService.WithStudioID(studioId).FindAll()
		if err != nil {
			errorhandler.HandleError(ctx, err)
			return
		}
		resp := response.NewBaseResponse("retrieved accounts by studio successfully", accounts)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	accounts, err := a.AccountService.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved accounts successfully", accounts)
	ctx.JSON(http.StatusOK, resp)
}

func (a *AccountControllerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	account, err := a.AccountService.WithID(id).FindOne()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved detail successfully", account)
	ctx.JSON(http.StatusOK, resp)
}

func (a *AccountControllerImpl) CreateOrUpdate(ctx *gin.Context) {
	var accountReq params.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&accountReq); err != nil {
		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err.Error()))
		return
	}

	account, err := a.AccountService.CreateOrUpdate(accountReq)
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

	result, err := a.AccountService.Update(id, req)
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

	if err := a.AccountService.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("account deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

// FindByStudio implements AccountController.
func (a *AccountControllerImpl) FindByStudio(ctx *gin.Context) {
	studioId := ctx.Param("studioId")

	accounts, err := a.AccountService.WithStudioID(studioId).FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved accounts by studio successfully", accounts)
	ctx.JSON(http.StatusOK, resp)
}

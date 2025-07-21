package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/helpers/response"
	"github.com/royhairul/live-studio-api/internal/domains/account/params"
	"github.com/royhairul/live-studio-api/internal/domains/account/service"
)

type AccountControllerImpl struct {
	AccountService service.AccountService
}

func NewAccountController(accountSvc service.AccountService) AccountController {
	return &AccountControllerImpl{AccountService: accountSvc}
}

func (a *AccountControllerImpl) FindAll(ctx *gin.Context) {
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

	account, err := a.AccountService.FindById(id)
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

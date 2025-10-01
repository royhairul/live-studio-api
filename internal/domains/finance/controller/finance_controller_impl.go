package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/helpers/response"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	"github.com/royhairul/live-studio-api/internal/domains/finance/service"
)

type FinanceControllerImpl struct {
	FinanceService service.FinanceService
}

func NewFinanceController(financeSvc service.FinanceService) FinanceController {
	return &FinanceControllerImpl{FinanceService: financeSvc}
}

func (f *FinanceControllerImpl) GetLiveFinance(ctx *gin.Context) {
	var req params.ShopeeLiveFinanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	dataFinance, err := f.FinanceService.FindAll(req)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved data successfully", dataFinance)
	ctx.JSON(http.StatusOK, resp)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Sesuaikan policy production
	},
}

// FindAllCommission implements FinanceController.
func (f *FinanceControllerImpl) FindAllCommission(ctx *gin.Context) {
}

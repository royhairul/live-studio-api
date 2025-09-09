package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/helpers/response"
	"github.com/royhairul/live-studio-api/internal/domains/performa/service"
)

type PerformaControllerImpl struct {
	service  service.PerformaService
	validate *validator.Validate
}

func NewPerformaController(service service.PerformaService, validate *validator.Validate) PerformaController {
	return &PerformaControllerImpl{service, validate}
}

// GetHosts implements PerformaController.
func (p *PerformaControllerImpl) GetHosts(ctx *gin.Context) {
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	result, err := p.service.GetHosts(startTime, endTime)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all performa host successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetHostByID implements PerformaController.
func (p *PerformaControllerImpl) GetHostByID(ctx *gin.Context) {
	id := ctx.Param("id")

	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	// Jika salah satu atau keduanya kosong, isi dengan hari ini
	if startTime == "" || endTime == "" {
		today := time.Now().Format("2006-01-02")
		startTime = today
		endTime = today
	}

	result, err := p.service.GetHostByID(id, startTime, endTime)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("performa host found", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetAccounts implements PerformaController.
func (p *PerformaControllerImpl) GetAccounts(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	result, err := p.service.GetAccounts(startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all performa studio successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetStudios implements PerformaController.
func (p *PerformaControllerImpl) GetStudios(ctx *gin.Context) {
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	result, err := p.service.GetStudios(startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all performa studio successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

// GetStudioByID implements PerformaController.
func (p *PerformaControllerImpl) GetStudioByID(ctx *gin.Context) {
	id := ctx.Param("id")

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	result, err := p.service.GetStudioByID(id, startDate, endDate)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("performa host found", result)
	ctx.JSON(http.StatusOK, resp)
}

// func (c *PerformaControllerImpl) Create(ctx *gin.Context) {
// 	var req params.CreatePerformaRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
// 		return
// 	}

// 	if err := c.validate.Struct(req); err != nil {
// 		errorhandler.HandleError(ctx, err)
// 		return
// 	}

// 	result, err := c.service.Create(req)
// 	if err != nil {
// 		errorhandler.HandleError(ctx, err)
// 		return
// 	}

// 	resp := response.NewBaseResponse("created performa successfully", result)
// 	ctx.JSON(http.StatusCreated, resp)
// }

// func (c *PerformaControllerImpl) Update(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	var req params.UpdatePerformaRequest

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		errorhandler.HandleError(ctx, errorhandler.NewBadRequestError("invalid request data", err))
// 		return
// 	}

// 	// Validasi hanya field yang tidak nil (optional)
// 	if err := c.validate.Struct(req); err != nil {
// 		errorhandler.HandleError(ctx, err)
// 		return
// 	}

// 	result, err := c.service.Update(id, req)
// 	if err != nil {
// 		errorhandler.HandleError(ctx, err)
// 		return
// 	}

// 	resp := response.NewBaseResponse("updated performa successfully", result)
// 	ctx.JSON(http.StatusOK, resp)
// }

// func (c *PerformaControllerImpl) FindAll(ctx *gin.Context) {
// 	result, err := c.service.FindAll()
// 	if err != nil {
// 		errorhandler.HandleError(ctx, err)
// 		return
// 	}

// 	resp := response.NewBaseResponse("retrieved all performa successfully", result)
// 	ctx.JSON(http.StatusOK, resp)
// }

// func (c *PerformaControllerImpl) FindByID(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	result, err := c.service.FindByID(id)
// 	if err != nil {
// 		errorhandler.HandleError(ctx, err)
// 		return
// 	}

// 	resp := response.NewBaseResponse("performa found", result)
// 	ctx.JSON(http.StatusOK, resp)
// }

// func (c *PerformaControllerImpl) Delete(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	if err := c.service.Delete(id); err != nil {
// 		errorhandler.HandleError(ctx, err)
// 		return
// 	}

// 	resp := response.NewBaseResponse("deleted performa successfully", nil)
// 	ctx.JSON(http.StatusOK, resp)
// }

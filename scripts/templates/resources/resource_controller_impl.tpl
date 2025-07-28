import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"{{Module}}/internal/domains/{{feature}}/params"
	"{{Module}}/internal/domains/{{feature}}/service"
)

type {{Feature}}ControllerImpl struct {
	service  service.{{Feature}}Service
	validate *validator.Validate
}

func New{{Feature}}Controller(service service.{{Feature}}Service, validate *validator.Validate) {{Feature}}Controller {
	return &{{Feature}}ControllerImpl{service, validate}
}

func (c *{{Feature}}ControllerImpl) Create(ctx *gin.Context) {
	var req params.Create{{Feature}}Request
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

	resp := response.NewBaseResponse("created {{feature}} successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

func (c *{{Feature}}ControllerImpl) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req params.Update{{Feature}}Request

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

	resp := response.NewBaseResponse("updated {{feature}} successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *{{Feature}}ControllerImpl) FindAll(ctx *gin.Context) {
	result, err := c.service.FindAll()
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("retrieved all {{feature}} successfully", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *{{Feature}}ControllerImpl) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.service.FindByID(id)
	if err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("{{feature}} found", result)
	ctx.JSON(http.StatusOK, resp)
}

func (c *{{Feature}}ControllerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.Delete(id); err != nil {
		errorhandler.HandleError(ctx, err)
		return
	}

	resp := response.NewBaseResponse("deleted {{feature}} successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

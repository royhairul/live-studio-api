package errorhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(code int, message string, details any) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewBadRequestError(msg string, details any) *CustomError {
	return NewCustomError(http.StatusBadRequest, msg, details)
}

func NewNotFoundError(msg string) *CustomError {
	return NewCustomError(http.StatusNotFound, msg, nil)
}

func NewUnauthorizedError(msg string) *CustomError {
	return NewCustomError(http.StatusUnauthorized, msg, nil)
}

func NewInternalServerError(err error) *CustomError {
	return NewCustomError(http.StatusInternalServerError, "Internal server error", err.Error())
}

func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {

	case *CustomError:
		if e.Details == nil {
			c.AbortWithStatusJSON(e.Code, gin.H{
				"error": e.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(e.Code, gin.H{
			"error":   e.Error(),
			"details": e.Details,
		})
		return

	case validator.ValidationErrors:
		validationErrs := FormatValidationError(e)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Validation Error",
			"details": validationErrs,
		})
		return

	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": e.Error(),
		})
	}
}

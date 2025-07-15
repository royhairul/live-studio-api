package errorhandler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, e := range err {
		field := strings.ToLower(e.Field())
		errors[field] = msgForTag(e.Tag(), e.Param())
	}
	return errors
}

func msgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "field is required"
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("value is too short (minimum is %s characters)", param)
	case "max":
		return fmt.Sprintf("value is too long (maximum is %s characters)", param)
	case "isFile":
		return "field must be a file"
	case "image":
		return "field must be an image file"
	case "fileSize":
		return fmt.Sprintf("file size is too large (maximum is %s MB)", param)
	default:
		return "invalid value"
	}
}

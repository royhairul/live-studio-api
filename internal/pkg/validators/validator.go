package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() *validator.Validate {
	validate := validator.New()

	// Custom validation
	validate.RegisterValidation("phoneid", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		// Check if the phone number matches the regex pattern
		match, _ := regexp.MatchString(`^08[0-9]{8,11}$`, phone)
		return match
	})

	Validate = validate

	return validate
}

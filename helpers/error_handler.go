package helpers

import (
	"reflect"
	"regexp"
	"strings"
)

// ParseValidationError parses raw validation error strings into a map of field errors
func ParseValidationError(errorMessage string, dto interface{}) map[string][]string {
	errors := strings.Split(errorMessage, "\n")
	result := make(map[string][]string)

	// Ambil map field dari struct DTO
	jsonFieldMap := getJSONFieldMap(dto)

	// Regex untuk mengekstrak field dan tag dari error message
	re := regexp.MustCompile(`Key: '.*\.(\w+)' Error:Field validation for '.*' failed on the '(\w+)' tag`)

	for _, err := range errors {
		matches := re.FindStringSubmatch(err)
		if len(matches) == 3 {
			field := matches[1]
			tag := matches[2]

			jsonField := jsonFieldMap[field]
			if jsonField == "" {
				jsonField = strings.ToLower(field)
			}
			
			message := humanizeTag(tag, matches[1])
			result[jsonField] = append(result[jsonField], message)
		}
	}

	return result
}

// humanizeTag returns a friendly error message based on validation tag
func humanizeTag(tag string, field string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "numeric":
		return field + " must be a number"
	case "min":
		return field + " is too short"
	case "max":
		return field + " is too long"
	default:
		return field + " is invalid (" + tag + ")"
	}
}

// getJSONFieldMap returns a map of struct field name -> json tag
func getJSONFieldMap(dto interface{}) map[string]string {
	t := reflect.TypeOf(dto)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldMap := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonTag != "" && jsonTag != "-" {
			fieldMap[field.Name] = jsonTag
		}
	}

	return fieldMap
}
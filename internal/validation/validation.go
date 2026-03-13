package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatError(err error, target interface{}) interface{} {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fieldMap := extractJSONFieldMap(target)
		result := make(map[string]string)

		for _, fieldErr := range validationErrors {
			fieldName := fieldMap[fieldErr.Field()]
			if fieldName == "" {
				fieldName = strings.ToLower(fieldErr.Field())
			}

			result[fieldName] = messageForTag(fieldName, fieldErr)
		}

		return result
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return "invalid JSON format"
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		fieldName := strings.ToLower(typeErr.Field)
		if fieldName == "" {
			fieldName = "unknown"
		}
		return fmt.Sprintf("field %s has invalid type", fieldName)
	}

	return err.Error()
}

func extractJSONFieldMap(target interface{}) map[string]string {
	result := make(map[string]string)

	t := reflect.TypeOf(target)
	if t == nil {
		return result
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonName := strings.Split(jsonTag, ",")[0]
		result[field.Name] = jsonName
	}

	return result
}

func messageForTag(fieldName string, fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldName)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fieldName, fieldErr.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fieldName, fieldErr.Param())
	default:
		return fmt.Sprintf("%s is invalid", fieldName)
	}
}
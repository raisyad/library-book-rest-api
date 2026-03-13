package validation

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
)

type sampleRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func TestFormatError_ValidationErrors(t *testing.T) {
	validate := validator.New()

	req := sampleRequest{
		Name:  "",
		Email: "not-an-email",
	}

	err := validate.Struct(req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	got := FormatError(err, req)

	errorsMap, ok := got.(map[string]string)
	if !ok {
		t.Fatalf("expected map[string]string, got %T", got)
	}

	if errorsMap["name"] != "name is required" {
		t.Fatalf("unexpected name error: %v", errorsMap["name"])
	}

	if errorsMap["email"] != "email must be a valid email address" {
		t.Fatalf("unexpected email error: %v", errorsMap["email"])
	}
}

func TestFormatError_SyntaxError(t *testing.T) {
	var syntaxErr *json.SyntaxError = &json.SyntaxError{Offset: 10}

	got := FormatError(syntaxErr, sampleRequest{})

	msg, ok := got.(string)
	if !ok {
		t.Fatalf("expected string, got %T", got)
	}

	if msg != "invalid JSON format" {
		t.Fatalf("expected invalid JSON format, got %s", msg)
	}
}

func TestFormatError_UnmarshalTypeError(t *testing.T) {
	typeErr := &json.UnmarshalTypeError{
		Field: "stock",
	}

	got := FormatError(typeErr, sampleRequest{})

	msg, ok := got.(string)
	if !ok {
		t.Fatalf("expected string, got %T", got)
	}

	if msg != "field stock has invalid type" {
		t.Fatalf("expected field stock has invalid type, got %s", msg)
	}
}
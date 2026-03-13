package book

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateBook_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&Service{})

	router := gin.New()
	router.POST("/books", handler.Create)

	body := `{
		"title": "",
		"author": "",
		"isbn": ""
	}`

	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	respBody := rec.Body.String()

	if !strings.Contains(respBody, "validation failed") {
		t.Fatalf("expected response to contain validation failed, got %s", respBody)
	}
}

func TestCreateBook_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&Service{})

	router := gin.New()
	router.POST("/books", handler.Create)

	body := `{"title":"Clean Code",}`

	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	respBody := rec.Body.String()

	if !strings.Contains(respBody, "validation failed") {
		t.Fatalf("expected response to contain validation failed, got %s", respBody)
	}

	if !strings.Contains(respBody, "invalid JSON format") {
		t.Fatalf("expected response to contain invalid JSON format, got %s", respBody)
	}
}
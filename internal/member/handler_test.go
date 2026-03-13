package member

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateMember_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&Service{})

	router := gin.New()
	router.POST("/members", handler.Create)

	body := `{
		"name": "Rizal",
		"email": "bukan-email"
	}`

	req := httptest.NewRequest(http.MethodPost, "/members", strings.NewReader(body))
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

	if !strings.Contains(respBody, "valid email address") {
		t.Fatalf("expected response to contain valid email address, got %s", respBody)
	}
}
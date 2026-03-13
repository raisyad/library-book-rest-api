package borrowing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateBorrowing_ValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&Service{})

	router := gin.New()
	router.POST("/borrowings", handler.Create)

	body := `{
		"member_id": 0,
		"book_id": 0
	}`

	req := httptest.NewRequest(http.MethodPost, "/borrowings", strings.NewReader(body))
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
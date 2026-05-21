package handle

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddBotHTTPValidation(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/bots", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	HTTPHandle.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

func TestDeleteBotHTTPValidation(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/bots", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	HTTPHandle.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

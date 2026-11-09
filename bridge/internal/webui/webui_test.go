package webui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDashboardServesIndexHTML(t *testing.T) {
	handler := Handler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Nerve Dispatch") {
		t.Errorf("expected dashboard HTML to contain 'Nerve Dispatch', got: %s", body)
	}
	if !strings.Contains(body, "SSE /events") {
		t.Errorf("expected dashboard HTML to contain 'SSE /events', got: %s", body)
	}
}

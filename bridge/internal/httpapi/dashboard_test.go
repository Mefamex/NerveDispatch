package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"whatsapp-bridge/internal/eventbus"
	"whatsapp-bridge/internal/msgstore"
)

func TestDashboardEndpoint(t *testing.T) {
	wa := &fakeWA{loggedIn: true, connected: true}
	store := msgstore.New(10)
	bus := eventbus.New()
	srv := New(wa, store, bus, time.Now())

	req := httptest.NewRequest(http.MethodGet, "/dashboard/", nil)
	rr := httptest.NewRecorder()
	srv.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Nerve Dispatch") {
		t.Errorf("expected dashboard HTML to contain 'Nerve Dispatch', got: %s", body)
	}
}

package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"whatsapp-bridge/internal/channels"
	"whatsapp-bridge/internal/eventbus"
	"whatsapp-bridge/internal/msgstore"
)

func TestChannelsHandler(t *testing.T) {
	wa := &fakeWA{loggedIn: true, connected: true}
	store := msgstore.New(10)
	bus := eventbus.New()
	srv := New(wa, store, bus, time.Now())

	mgr := channels.NewManager()
	srv.SetChannelsManager(mgr)

	req := httptest.NewRequest(http.MethodGet, "/channels", nil)
	rr := httptest.NewRecorder()

	srv.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("beklenen HTTP 200, got %d", rr.Code)
	}

	var res struct {
		Status   string                   `json:"status"`
		Channels []channels.ChannelStatus `json:"channels"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("json parse hatası: %v", err)
	}

	if res.Status != "ok" {
		t.Errorf("beklenen status 'ok', got %q", res.Status)
	}

	if len(res.Channels) == 0 {
		t.Errorf("en az 1 kanal (whatsapp) bekleniyordu, got 0")
	}
	if res.Channels[0].Name != channels.ChannelWhatsApp {
		t.Errorf("ilk kanal whatsapp olmalıydı, got %s", res.Channels[0].Name)
	}
}

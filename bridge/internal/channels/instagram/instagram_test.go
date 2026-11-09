package instagram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"whatsapp-bridge/internal/channels"
)

func TestInstagramAdapterBasics(t *testing.T) {
	cfg := Config{
		AccessToken: "ig-test-token",
		PageID:      "page-12345",
	}

	adapter := NewAdapter(cfg, nil)
	if adapter.Name() != channels.ChannelInstagram {
		t.Errorf("beklenen kanal adı instagram, got %s", adapter.Name())
	}

	status := adapter.Status(context.Background())
	if status.Name != channels.ChannelInstagram || !status.Paired {
		t.Errorf("beklenen status.Paired=true, got %+v", status)
	}
}

func TestInstagramConnectAndSendText(t *testing.T) {
	var receivedRecipient string
	var receivedText string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Connect endpoint
		if r.URL.Path == "/page-12345" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"name": "NerveDispatch Page", "username": "nervedispatch"}`))
			return
		}

		// Send endpoint
		if r.URL.Path == "/page-12345/messages" && r.Method == http.MethodPost {
			var body struct {
				Recipient struct {
					ID string `json:"id"`
				} `json:"recipient"`
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			receivedRecipient = body.Recipient.ID
			receivedText = body.Message.Text

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"recipient_id": "user-999", "message_id": "mid.123"}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	adapter := NewAdapter(Config{
		AccessToken: "ig-test-token",
		PageID:      "page-12345",
	}, nil)
	adapter.apiBase = ts.URL

	ctx := context.Background()

	// 1. Connect
	if err := adapter.Connect(ctx); err != nil {
		t.Fatalf("Connect hatası: %v", err)
	}
	status := adapter.Status(ctx)
	if !status.Connected || status.Account != "@nervedispatch" {
		t.Errorf("Connect sonrası beklenen @nervedispatch bağlı durumu, got %+v", status)
	}

	// 2. SendText
	if err := adapter.SendText(ctx, "user-999", "Merhaba Instagram Direct"); err != nil {
		t.Fatalf("SendText hatası: %v", err)
	}
	if receivedRecipient != "user-999" || receivedText != "Merhaba Instagram Direct" {
		t.Errorf("beklenen recipient=user-999, text='Merhaba Instagram Direct', got %s, %s", receivedRecipient, receivedText)
	}
}

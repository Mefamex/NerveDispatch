package telegram

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"whatsapp-bridge/internal/channels"
)

func TestTelegramAdapterBasics(t *testing.T) {
	cfg := Config{
		Token:   "test-token-123",
		AdminID: 975678139,
	}

	adapter := NewAdapter(cfg, nil)
	if adapter.Name() != channels.ChannelTelegram {
		t.Errorf("beklenen kanal adı telegram, got %s", adapter.Name())
	}
	if adapter.cfg.AdminID != 975678139 {
		t.Errorf("beklenen admin ID 975678139, got %d", adapter.cfg.AdminID)
	}

	status := adapter.Status(context.Background())
	if status.Name != channels.ChannelTelegram || !status.Paired {
		t.Errorf("beklenen status.Paired=true, got %+v", status)
	}
}

func TestTelegramSendText(t *testing.T) {
	var receivedChatID string
	var receivedText string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bottest-token-123/sendMessage" {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			receivedChatID = body["chat_id"].(string)
			receivedText = body["text"].(string)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok": true, "result": {"message_id": 1001}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	adapter := NewAdapter(Config{Token: "test-token-123", AdminID: 975678139}, nil)
	adapter.apiBase = ts.URL

	ctx := context.Background()
	err := adapter.SendText(ctx, "975678139", "NerveDispatch Telegram test mesajı")
	if err != nil {
		t.Fatalf("SendText hatası: %v", err)
	}

	if receivedChatID != "975678139" {
		t.Errorf("chat_id 975678139 olmalı, got %q", receivedChatID)
	}
	if receivedText != "NerveDispatch Telegram test mesajı" {
		t.Errorf("text eşleşmedi, got %q", receivedText)
	}
}

package discord

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"whatsapp-bridge/internal/channels"
)

func TestDiscordAdapterBasics(t *testing.T) {
	cfg := Config{
		Token:   "discord-test-token",
		GuildID: "guild-123",
	}

	adapter := NewAdapter(cfg, nil)
	if adapter.Name() != channels.ChannelDiscord {
		t.Errorf("beklenen kanal adı discord, got %s", adapter.Name())
	}

	status := adapter.Status(context.Background())
	if status.Name != channels.ChannelDiscord || !status.Paired {
		t.Errorf("beklenen status.Paired=true, got %+v", status)
	}
}

func TestDiscordSendTextAndEmbed(t *testing.T) {
	var receivedContent string
	var receivedEmbedTitle string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/channels/987654321/messages" {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if c, ok := body["content"].(string); ok {
				receivedContent = c
			}
			if embeds, ok := body["embeds"].([]any); ok && len(embeds) > 0 {
				emMap := embeds[0].(map[string]any)
				receivedEmbedTitle, _ = emMap["title"].(string)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id": "msg-111"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	adapter := NewAdapter(Config{Token: "discord-test-token"}, nil)
	adapter.apiBase = ts.URL

	ctx := context.Background()

	// 1. Düz metin testi
	if err := adapter.SendText(ctx, "987654321", "NerveDispatch Discord Mesajı"); err != nil {
		t.Fatalf("SendText hatası: %v", err)
	}
	if receivedContent != "NerveDispatch Discord Mesajı" {
		t.Errorf("content eşleşmedi, got %q", receivedContent)
	}

	// 2. Zengin Embed testi
	embed := Embed{
		Title:       "Sistem Durumu",
		Description: "Tüm servisler aktif ve çalışıyor.",
		Color:       0x00FF00,
	}
	if err := adapter.SendEmbed(ctx, "987654321", embed); err != nil {
		t.Fatalf("SendEmbed hatası: %v", err)
	}
	if receivedEmbedTitle != "Sistem Durumu" {
		t.Errorf("embed title eşleşmedi, got %q", receivedEmbedTitle)
	}
}

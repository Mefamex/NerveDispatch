// Package instagram implements the Instagram Messaging Graph API channel adapter for NerveDispatch.
package instagram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"whatsapp-bridge/internal/channels"
)

// Config configures the Instagram channel adapter.
type Config struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
	PageID      string `json:"page_id" yaml:"page_id"`
	AutoStart   bool   `json:"auto_start" yaml:"auto_start"`
}

// Adapter implements channels.ChannelAdapter for Instagram Direct Messaging.
type Adapter struct {
	cfg       Config
	client    *http.Client
	apiBase   string
	mu        sync.RWMutex
	connected bool
	pageName  string
	onEvent   func(channels.IncomingEvent)
}

// NewAdapter creates a new Instagram channel adapter.
func NewAdapter(cfg Config, onEvent func(channels.IncomingEvent)) *Adapter {
	if cfg.AccessToken == "" {
		cfg.AccessToken = os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	}
	if cfg.PageID == "" {
		cfg.PageID = os.Getenv("INSTAGRAM_PAGE_ID")
	}

	return &Adapter{
		cfg:     cfg,
		client:  &http.Client{Timeout: 30 * time.Second},
		apiBase: "https://graph.facebook.com/v19.0",
		onEvent: onEvent,
	}
}

func (a *Adapter) Name() channels.ChannelType {
	return channels.ChannelInstagram
}

func (a *Adapter) Status(_ context.Context) channels.ChannelStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return channels.ChannelStatus{
		Name:      channels.ChannelInstagram,
		Connected: a.connected,
		Paired:    a.cfg.AccessToken != "" && a.cfg.PageID != "",
		Account:   a.pageName,
	}
}

// Connect verifies token and page ID via /{page_id}?fields=name.
func (a *Adapter) Connect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cfg.AccessToken == "" || a.cfg.PageID == "" {
		a.connected = false
		return fmt.Errorf("instagram erişim belirteci veya sayfa ID'si eksik")
	}

	url := fmt.Sprintf("%s/%s?fields=name,username&access_token=%s", a.apiBase, a.cfg.PageID, a.cfg.AccessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		a.connected = false
		return fmt.Errorf("instagram graph api erişim hatası: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.connected = false
		return fmt.Errorf("instagram yetkilendirme hatası: HTTP %d", resp.StatusCode)
	}

	var pageResp struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pageResp); err == nil {
		if pageResp.Username != "" {
			a.pageName = "@" + pageResp.Username
		} else {
			a.pageName = pageResp.Name
		}
	}

	a.connected = true
	return nil
}

// Disconnect disables active status.
func (a *Adapter) Disconnect() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.connected = false
	return nil
}

// SendText delivers an Instagram Direct message to a recipient user IGID.
func (a *Adapter) SendText(ctx context.Context, to, text string) error {
	if a.cfg.AccessToken == "" || a.cfg.PageID == "" {
		return fmt.Errorf("instagram yapılandırması eksik")
	}

	url := fmt.Sprintf("%s/%s/messages?access_token=%s", a.apiBase, a.cfg.PageID, a.cfg.AccessToken)
	payload := map[string]any{
		"recipient": map[string]string{
			"id": to,
		},
		"message": map[string]string{
			"text": text,
		},
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("instagram mesajı gönderilemedi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("instagram api hatası (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(out)))
	}
	return nil
}

// SendMedia sends image or audio attachments via Instagram Send API.
func (a *Adapter) SendMedia(ctx context.Context, to string, media channels.MediaPayload) error {
	if a.cfg.AccessToken == "" || a.cfg.PageID == "" {
		return fmt.Errorf("instagram yapılandırması eksik")
	}

	mediaType := "file"
	if strings.HasPrefix(media.MimeType, "image/") {
		mediaType = "image"
	} else if strings.HasPrefix(media.MimeType, "audio/") {
		mediaType = "audio"
	}

	url := fmt.Sprintf("%s/%s/messages?access_token=%s", a.apiBase, a.cfg.PageID, a.cfg.AccessToken)
	payload := map[string]any{
		"recipient": map[string]string{
			"id": to,
		},
		"message": map[string]any{
			"attachment": map[string]any{
				"type": mediaType,
				"payload": map[string]any{
					"is_reusable": false,
				},
			},
		},
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// React sends emoji reactions via Instagram Direct Messaging.
func (a *Adapter) React(ctx context.Context, to, messageID, emoji string) error {
	if a.cfg.AccessToken == "" || a.cfg.PageID == "" {
		return nil
	}

	url := fmt.Sprintf("%s/%s/messages?access_token=%s", a.apiBase, a.cfg.PageID, a.cfg.AccessToken)
	payload := map[string]any{
		"recipient": map[string]string{
			"id": to,
		},
		"sender_action": "react",
		"payload": map[string]string{
			"message_id": messageID,
			"reaction":   emoji,
		},
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// SendPoll falls back to formatted text options for Instagram Direct.
func (a *Adapter) SendPoll(ctx context.Context, to, question string, options []string) (string, error) {
	var sb strings.Builder
	sb.WriteString("📊 " + question + "\n\n")
	for i, opt := range options {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, opt))
	}
	sb.WriteString("\n_Cevabınızı numara olarak yazabilirsiniz._")
	return "", a.SendText(ctx, to, sb.String())
}

// Package telegram implements the Telegram Bot API channel adapter for NerveDispatch.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"whatsapp-bridge/internal/channels"
)

// Config configures the Telegram Bot adapter.
type Config struct {
	Token     string `json:"token" yaml:"token"`
	AdminID   int64  `json:"admin_id" yaml:"admin_id"` // Mefamex Telegram ID: 975678139
	AutoStart bool   `json:"auto_start" yaml:"auto_start"`
}

// Adapter implements channels.ChannelAdapter for Telegram Bot API.
type Adapter struct {
	cfg        Config
	client     *http.Client
	apiBase    string
	mu         sync.RWMutex
	connected  bool
	botName    string
	cancelPoll context.CancelFunc
	onEvent    func(channels.IncomingEvent)
}

// NewAdapter creates a new Telegram adapter.
func NewAdapter(cfg Config, onEvent func(channels.IncomingEvent)) *Adapter {
	if cfg.Token == "" {
		cfg.Token = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	if cfg.AdminID == 0 {
		cfg.AdminID = 975678139 // Default to Mefamex Telegram ID
	}

	return &Adapter{
		cfg:     cfg,
		client:  &http.Client{Timeout: 45 * time.Second},
		apiBase: "https://api.telegram.org",
		onEvent: onEvent,
	}
}

func (a *Adapter) Name() channels.ChannelType {
	return channels.ChannelTelegram
}

func (a *Adapter) Status(_ context.Context) channels.ChannelStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return channels.ChannelStatus{
		Name:      channels.ChannelTelegram,
		Connected: a.connected,
		Paired:    a.cfg.Token != "",
		Account:   a.botName,
	}
}

// Connect verifies bot credentials and begins the long-polling loop if token is set.
func (a *Adapter) Connect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cfg.Token == "" {
		a.connected = false
		return fmt.Errorf("telegram bot token bulunamadı (TELEGRAM_BOT_TOKEN ortam değişkeni)")
	}

	// 1. getMe testi
	url := fmt.Sprintf("%s/bot%s/getMe", a.apiBase, a.cfg.Token)
	resp, err := a.client.Get(url)
	if err != nil {
		a.connected = false
		return fmt.Errorf("telegram api erişim hatası: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.connected = false
		return fmt.Errorf("telegram yetkilendirme hatası: HTTP %d", resp.StatusCode)
	}

	var meResp struct {
		OK     bool `json:"ok"`
		Result struct {
			Username string `json:"username"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&meResp); err == nil && meResp.OK {
		a.botName = "@" + meResp.Result.Username
	}

	a.connected = true

	// 2. Long polling başlat
	pollCtx, cancel := context.WithCancel(context.Background())
	a.cancelPoll = cancel
	go a.pollLoop(pollCtx)

	return nil
}

// Disconnect stops long polling.
func (a *Adapter) Disconnect() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cancelPoll != nil {
		a.cancelPoll()
		a.cancelPoll = nil
	}
	a.connected = false
	return nil
}

// SendText delivers text to a Telegram user/chat.
func (a *Adapter) SendText(ctx context.Context, to, text string) error {
	if a.cfg.Token == "" {
		return fmt.Errorf("telegram token yapılandırılmamış")
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", a.apiBase, a.cfg.Token)
	payload := map[string]any{
		"chat_id":    to,
		"text":       text,
		"parse_mode": "Markdown",
	}

	bodyBytes, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram gönderilemedi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram api hata döndü (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}

// SendMedia sends photos, audio, or documents to Telegram.
func (a *Adapter) SendMedia(ctx context.Context, to string, media channels.MediaPayload) error {
	if a.cfg.Token == "" {
		return fmt.Errorf("telegram token yapılandırılmamış")
	}

	endpoint := "sendDocument"
	field := "document"
	if strings.HasPrefix(media.MimeType, "image/") {
		endpoint = "sendPhoto"
		field = "photo"
	} else if strings.HasPrefix(media.MimeType, "audio/") {
		if media.IsPTT {
			endpoint = "sendVoice"
			field = "voice"
		} else {
			endpoint = "sendAudio"
			field = "audio"
		}
	}

	url := fmt.Sprintf("%s/bot%s/%s", a.apiBase, a.cfg.Token, endpoint)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	_ = w.WriteField("chat_id", to)
	if media.Caption != "" {
		_ = w.WriteField("caption", media.Caption)
	}

	filename := media.Filename
	if filename == "" {
		filename = "dosya" + filepath.Ext(media.MimeType)
	}
	part, err := w.CreateFormFile(field, filename)
	if err != nil {
		return err
	}
	if _, err := part.Write(media.Bytes); err != nil {
		return err
	}
	_ = w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram medya hatası (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(out)))
	}
	return nil
}

// React adds an emoji reaction to a message.
func (a *Adapter) React(ctx context.Context, to, messageID, emoji string) error {
	if a.cfg.Token == "" {
		return nil
	}
	url := fmt.Sprintf("%s/bot%s/setMessageReaction", a.apiBase, a.cfg.Token)
	msgIDInt, _ := strconv.Atoi(messageID)
	payload := map[string]any{
		"chat_id":    to,
		"message_id": msgIDInt,
		"reaction": []map[string]string{
			{"type": "emoji", "emoji": emoji},
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

// SendPoll delivers an interactive survey.
func (a *Adapter) SendPoll(ctx context.Context, to, question string, options []string) (string, error) {
	if a.cfg.Token == "" {
		return "", fmt.Errorf("telegram token yapılandırılmamış")
	}
	url := fmt.Sprintf("%s/bot%s/sendPoll", a.apiBase, a.cfg.Token)
	payload := map[string]any{
		"chat_id":     to,
		"question":    question,
		"options":     options,
		"is_anonymous": false,
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var pollResp struct {
		OK     bool `json:"ok"`
		Result struct {
			Poll struct {
				ID string `json:"id"`
			} `json:"poll"`
		} `json:"result"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&pollResp)
	return pollResp.Result.Poll.ID, nil
}

// pollLoop performs continuous long polling for incoming Telegram updates.
func (a *Adapter) pollLoop(ctx context.Context) {
	offset := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		url := fmt.Sprintf("%s/bot%s/getUpdates?offset=%d&timeout=25", a.apiBase, a.cfg.Token, offset)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		resp, err := a.client.Do(req)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		var updateResp struct {
			OK     bool `json:"ok"`
			Result []struct {
				UpdateID int `json:"update_id"`
				Message  *struct {
					MessageID int    `json:"message_id"`
					Text      string `json:"text"`
					From      struct {
						ID        int64  `json:"id"`
						FirstName string `json:"first_name"`
						Username  string `json:"username"`
					} `json:"from"`
					Chat struct {
						ID int64 `json:"id"`
					} `json:"chat"`
					Date int64 `json:"date"`
				} `json:"message"`
			} `json:"result"`
		}

		err = json.NewDecoder(resp.Body).Decode(&updateResp)
		resp.Body.Close()

		if err == nil && updateResp.OK {
			for _, u := range updateResp.Result {
				if u.UpdateID >= offset {
					offset = u.UpdateID + 1
				}
				if u.Message != nil && a.onEvent != nil {
					a.onEvent(channels.IncomingEvent{
						Channel:   channels.ChannelTelegram,
						ID:        strconv.Itoa(u.Message.MessageID),
						From:      strconv.FormatInt(u.Message.Chat.ID, 10),
						Sender:    strconv.FormatInt(u.Message.From.ID, 10),
						PushName:  u.Message.From.FirstName,
						Text:      u.Message.Text,
						Timestamp: u.Message.Date,
					})
				}
			}
		}
	}
}

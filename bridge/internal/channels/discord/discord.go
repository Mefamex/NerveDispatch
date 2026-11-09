// Package discord implements the Discord Bot channel adapter for NerveDispatch.
package discord

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
	"strings"
	"sync"
	"time"

	"whatsapp-bridge/internal/channels"
)

// Config configures the Discord Bot adapter.
type Config struct {
	Token     string `json:"token" yaml:"token"`
	GuildID   string `json:"guild_id" yaml:"guild_id"`
	AutoStart bool   `json:"auto_start" yaml:"auto_start"`
}

// EmbedField represents a field inside a rich Discord embed.
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Embed represents a Discord rich embed message.
type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Color       int          `json:"color,omitempty"` // RGB int (e.g. 0x00FF00 for green)
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      *struct {
		Text string `json:"text"`
	} `json:"footer,omitempty"`
}

// Adapter implements channels.ChannelAdapter for Discord.
type Adapter struct {
	cfg       Config
	client    *http.Client
	apiBase   string
	mu        sync.RWMutex
	connected bool
	botName   string
	onEvent   func(channels.IncomingEvent)
}

// NewAdapter creates a new Discord channel adapter.
func NewAdapter(cfg Config, onEvent func(channels.IncomingEvent)) *Adapter {
	if cfg.Token == "" {
		cfg.Token = os.Getenv("DISCORD_BOT_TOKEN")
	}

	return &Adapter{
		cfg:     cfg,
		client:  &http.Client{Timeout: 30 * time.Second},
		apiBase: "https://discord.com/api/v10",
		onEvent: onEvent,
	}
}

func (a *Adapter) Name() channels.ChannelType {
	return channels.ChannelDiscord
}

func (a *Adapter) Status(_ context.Context) channels.ChannelStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return channels.ChannelStatus{
		Name:      channels.ChannelDiscord,
		Connected: a.connected,
		Paired:    a.cfg.Token != "",
		Account:   a.botName,
	}
}

// Connect verifies bot credentials via /users/@me.
func (a *Adapter) Connect(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cfg.Token == "" {
		a.connected = false
		return fmt.Errorf("discord bot token bulunamadı (DISCORD_BOT_TOKEN veya config.yaml)")
	}

	url := fmt.Sprintf("%s/users/@me", a.apiBase)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+a.cfg.Token)

	resp, err := a.client.Do(req)
	if err != nil {
		a.connected = false
		return fmt.Errorf("discord api erişim hatası: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.connected = false
		return fmt.Errorf("discord yetkilendirme hatası: HTTP %d", resp.StatusCode)
	}

	var userResp struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err == nil {
		a.botName = "@" + userResp.Username
	}

	a.connected = true
	return nil
}

// Disconnect closes active sessions.
func (a *Adapter) Disconnect() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.connected = false
	return nil
}

// SendText posts a message to a Discord channel or DM.
func (a *Adapter) SendText(ctx context.Context, to, text string) error {
	if a.cfg.Token == "" {
		return fmt.Errorf("discord token yapılandırılmamış")
	}

	url := fmt.Sprintf("%s/channels/%s/messages", a.apiBase, to)
	payload := map[string]any{
		"content": text,
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+a.cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("discord mesajı gönderilemedi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		out, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord api hatası (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(out)))
	}
	return nil
}

// SendEmbed sends a rich embedded message.
func (a *Adapter) SendEmbed(ctx context.Context, to string, embed Embed) error {
	if a.cfg.Token == "" {
		return fmt.Errorf("discord token yapılandırılmamış")
	}

	url := fmt.Sprintf("%s/channels/%s/messages", a.apiBase, to)
	payload := map[string]any{
		"embeds": []Embed{embed},
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+a.cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		out, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord embed hatası (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(out)))
	}
	return nil
}

// SendMedia uploads an attachment to a Discord channel.
func (a *Adapter) SendMedia(ctx context.Context, to string, media channels.MediaPayload) error {
	if a.cfg.Token == "" {
		return fmt.Errorf("discord token yapılandırılmamış")
	}

	url := fmt.Sprintf("%s/channels/%s/messages", a.apiBase, to)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if media.Caption != "" {
		_ = w.WriteField("content", media.Caption)
	}

	filename := media.Filename
	if filename == "" {
		filename = "attachment" + filepath.Ext(media.MimeType)
	}

	part, err := w.CreateFormFile("files[0]", filename)
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
	req.Header.Set("Authorization", "Bot "+a.cfg.Token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		out, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord medya hatası (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(out)))
	}
	return nil
}

// React adds an emoji reaction to a Discord message.
func (a *Adapter) React(ctx context.Context, to, messageID, emoji string) error {
	if a.cfg.Token == "" {
		return nil
	}

	url := fmt.Sprintf("%s/channels/%s/messages/%s/reactions/%s/@me", a.apiBase, to, messageID, emoji)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+a.cfg.Token)

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// SendPoll sends a native Discord survey/poll.
func (a *Adapter) SendPoll(ctx context.Context, to, question string, options []string) (string, error) {
	if a.cfg.Token == "" {
		return "", fmt.Errorf("discord token yapılandırılmamış")
	}

	url := fmt.Sprintf("%s/channels/%s/messages", a.apiBase, to)
	answers := make([]map[string]any, 0, len(options))
	for _, opt := range options {
		answers = append(answers, map[string]any{
			"poll_media": map[string]string{
				"text": opt,
			},
		})
	}

	payload := map[string]any{
		"poll": map[string]any{
			"question": map[string]string{
				"text": question,
			},
			"answers":  answers,
			"duration": 24, // 24 hours
		},
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bot "+a.cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var msgResp struct {
		ID string `json:"id"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&msgResp)
	return msgResp.ID, nil
}

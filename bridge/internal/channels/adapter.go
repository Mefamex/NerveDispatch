// Package channels defines the unified multi-channel adapter interface for NerveDispatch.
package channels

import (
	"context"
	"fmt"
	"sync"
)

// ChannelType identifies an omnichannel provider.
type ChannelType string

const (
	ChannelWhatsApp  ChannelType = "whatsapp"
	ChannelTelegram  ChannelType = "telegram"
	ChannelDiscord   ChannelType = "discord"
	ChannelInstagram ChannelType = "instagram"
)

// ChannelStatus represents the live health of an omnichannel adapter.
type ChannelStatus struct {
	Name      ChannelType `json:"name"`
	Connected bool        `json:"connected"`
	Paired    bool        `json:"paired"`
	Account   string      `json:"account,omitempty"`
	LastError string      `json:"last_error,omitempty"`
}

// MediaPayload wraps outbound media bytes and metadata.
type MediaPayload struct {
	Bytes    []byte
	MimeType string
	Filename string
	Caption  string
	IsPTT    bool
}

// IncomingEvent is the normalized cross-channel event emitted to the SSE and dispatcher bus.
type IncomingEvent struct {
	Channel   ChannelType `json:"channel"`
	ID        string      `json:"id"`
	Type      string      `json:"type,omitempty"`
	From      string      `json:"from"`
	Sender    string      `json:"sender"`
	PushName  string      `json:"push_name"`
	Text      string      `json:"text"`
	Timestamp int64       `json:"timestamp"`
	IsGroup   bool        `json:"is_group"`
}

// ChannelAdapter is the standard interface implemented by all channel providers.
type ChannelAdapter interface {
	Name() ChannelType
	Connect(ctx context.Context) error
	Disconnect() error
	SendText(ctx context.Context, to, text string) error
	SendMedia(ctx context.Context, to string, media MediaPayload) error
	React(ctx context.Context, to, messageID, emoji string) error
	SendPoll(ctx context.Context, to, question string, options []string) (string, error)
	Status(ctx context.Context) ChannelStatus
}

// Manager orchestrates all registered channel adapters.
type Manager struct {
	mu       sync.RWMutex
	adapters map[ChannelType]ChannelAdapter
}

// NewManager returns an empty channel manager.
func NewManager() *Manager {
	return &Manager{
		adapters: make(map[ChannelType]ChannelAdapter),
	}
}

// Register adds a channel adapter.
func (m *Manager) Register(adapter ChannelAdapter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.adapters[adapter.Name()] = adapter
}

// Get returns the adapter for the requested channel.
func (m *Manager) Get(name ChannelType) (ChannelAdapter, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.adapters[name]
	return a, ok
}

// ListStatus collects statuses across all registered adapters.
func (m *Manager) ListStatus(ctx context.Context) []ChannelStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]ChannelStatus, 0, len(m.adapters))
	for _, a := range m.adapters {
		res = append(res, a.Status(ctx))
	}
	return res
}

// SendText sends text via the specified channel, falling back to other channels if needed.
func (m *Manager) SendText(ctx context.Context, channel ChannelType, to, text string) error {
	a, ok := m.Get(channel)
	if !ok {
		return fmt.Errorf("kanal adaptörü bulunamadı: %s", channel)
	}
	return a.SendText(ctx, to, text)
}

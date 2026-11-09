package channels

import (
	"context"
	"testing"
)

type mockAdapter struct {
	name   ChannelType
	status ChannelStatus
	sent   []string
}

func (m *mockAdapter) Name() ChannelType { return m.name }
func (m *mockAdapter) Connect(_ context.Context) error { return nil }
func (m *mockAdapter) Disconnect() error { return nil }
func (m *mockAdapter) SendText(_ context.Context, to, text string) error {
	m.sent = append(m.sent, to+":"+text)
	return nil
}
func (m *mockAdapter) SendMedia(_ context.Context, _ string, _ MediaPayload) error { return nil }
func (m *mockAdapter) React(_ context.Context, _, _, _ string) error { return nil }
func (m *mockAdapter) SendPoll(_ context.Context, _, _ string, _ []string) (string, error) {
	return "poll-1", nil
}
func (m *mockAdapter) Status(_ context.Context) ChannelStatus { return m.status }

func TestChannelManager(t *testing.T) {
	mgr := NewManager()

	wa := &mockAdapter{name: ChannelWhatsApp, status: ChannelStatus{Name: ChannelWhatsApp, Connected: true}}
	tg := &mockAdapter{name: ChannelTelegram, status: ChannelStatus{Name: ChannelTelegram, Connected: true, Account: "@mefamex_bot"}}

	mgr.Register(wa)
	mgr.Register(tg)

	statuses := mgr.ListStatus(context.Background())
	if len(statuses) != 2 {
		t.Fatalf("beklenen 2 kanal durumu, got %d", len(statuses))
	}

	ctx := context.Background()
	if err := mgr.SendText(ctx, ChannelTelegram, "975678139", "merhaba telegram"); err != nil {
		t.Fatalf("telegram gönderimi başarısız: %v", err)
	}

	if len(tg.sent) != 1 || tg.sent[0] != "975678139:merhaba telegram" {
		t.Errorf("beklenen mesaj iletilmedi: %v", tg.sent)
	}
}

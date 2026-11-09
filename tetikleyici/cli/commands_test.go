package cli

import (
	"testing"
)

func TestDispatcherHelp(t *testing.T) {
	handled, err := Dispatcher([]string{"help"})
	if err != nil {
		t.Fatalf("help failed: %v", err)
	}
	if !handled {
		t.Errorf("Dispatcher(help) = false, want true")
	}
}

func TestDispatcherFlagBypass(t *testing.T) {
	// Bayraklı çağrılar daemon moduna düşmeli (handled == false)
	handled, err := Dispatcher([]string{"--config", "tetikleyici.yaml"})
	if err != nil {
		t.Fatalf("flag bypass failed: %v", err)
	}
	if handled {
		t.Errorf("Dispatcher(--config) = true, want false")
	}
}

func TestDispatcherConfigCheck(t *testing.T) {
	handled, err := Dispatcher([]string{"config", "check"})
	if err != nil {
		t.Fatalf("config check failed: %v", err)
	}
	if !handled {
		t.Errorf("Dispatcher(config check) = false, want true")
	}
}

func TestDispatcherChannelList(t *testing.T) {
	handled, err := Dispatcher([]string{"channel", "list"})
	if err != nil {
		t.Fatalf("channel list failed: %v", err)
	}
	if !handled {
		t.Errorf("Dispatcher(channel list) = false, want true")
	}
}

func TestDispatcherSessionList(t *testing.T) {
	handled, err := Dispatcher([]string{"session", "list"})
	if err != nil {
		t.Fatalf("session list failed: %v", err)
	}
	if !handled {
		t.Errorf("Dispatcher(session list) = false, want true")
	}
}

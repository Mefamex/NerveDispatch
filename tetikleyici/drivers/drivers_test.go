package drivers

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSessionStore(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewSessionStore(tmpDir)

	// Başlangıçta boş olmalı
	if got := store.Get("opencode"); got != "" {
		t.Errorf("Get(opencode) = %q, want empty", got)
	}

	// Set & Get opencode
	if err := store.Set("opencode", "ses_test123"); err != nil {
		t.Fatalf("Set(opencode) failed: %v", err)
	}
	if got := store.Get("opencode"); got != "ses_test123" {
		t.Errorf("Get(opencode) = %q, want ses_test123", got)
	}

	// Set & Get agy (conversation_id)
	if err := store.Set("agy", "conv_abc789"); err != nil {
		t.Fatalf("Set(agy) failed: %v", err)
	}
	if got := store.Get("agy"); got != "conv_abc789" {
		t.Errorf("Get(agy) = %q, want conv_abc789", got)
	}

	// Set & Get claude
	if err := store.Set("claude", "uuid_c99"); err != nil {
		t.Fatalf("Set(claude) failed: %v", err)
	}
	if got := store.Get("claude"); got != "uuid_c99" {
		t.Errorf("Get(claude) = %q, want uuid_c99", got)
	}
}

func TestDriverNames(t *testing.T) {
	agy := NewAgyDriver("")
	if agy.Name() != "agy" {
		t.Errorf("AgyDriver.Name() = %q, want agy", agy.Name())
	}

	claude := NewClaudeDriver("")
	if claude.Name() != "claude" {
		t.Errorf("ClaudeDriver.Name() = %q, want claude", claude.Name())
	}

	opencode := NewOpencodeDriver("")
	if opencode.Name() != "opencode" {
		t.Errorf("OpencodeDriver.Name() = %q, want opencode", opencode.Name())
	}
}

func TestMockDriverExecution(t *testing.T) {
	// echo ile mock driver testi
	tmpDir := t.TempDir()
	mockBin := filepath.Join(tmpDir, "mock_driver.sh")
	script := `#!/usr/bin/env bash
echo "REPLY: $@"
`
	if err := os.WriteFile(mockBin, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to write mock script: %v", err)
	}

	driver := &AgyDriver{BinPath: mockBin}
	out, err := driver.Run(context.Background(), "conv1", "merhaba dünya", nil, tmpDir)
	if err != nil {
		t.Fatalf("driver.Run failed: %v", err)
	}

	if out == "" {
		t.Error("expected non-empty output")
	}
}

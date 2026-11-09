package drivers

import (
	"context"
	"os"
	"path/filepath"
	"strings"
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

	if out.Text == "" {
		t.Error("expected non-empty output")
	}
}

// argsBin, aldığı argümanları dosyaya yazıp verilen çıktıyı basan sahte CLI kurar.
func argsBin(t *testing.T, stdout string) (bin, argsFile string) {
	t.Helper()
	dir := t.TempDir()
	argsFile = filepath.Join(dir, "args.txt")
	bin = filepath.Join(dir, "fake.sh")
	script := "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\" > " + argsFile + "\ncat <<'OUT'\n" + stdout + "\nOUT\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return bin, argsFile
}

func readArgs(t *testing.T, path string) []string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimRight(string(b), "\n"), "\n")
}

func argValue(args []string, flag string) (string, bool) {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1], true
		}
	}
	return "", false
}

func TestAgyNewConversationReturnsIDFromJSON(t *testing.T) {
	bin, af := argsBin(t, `{"conversation_id":"conv-yeni","status":"SUCCESS","response":"merhaba\n"}`)
	res, err := (&AgyDriver{BinPath: bin}).Run(context.Background(), "", "selam", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if res.SessionID != "conv-yeni" || strings.TrimSpace(res.Text) != "merhaba" {
		t.Fatalf("res = %+v", res)
	}
	args := readArgs(t, af)
	if _, ok := argValue(args, "--conversation"); ok {
		t.Errorf("yeni sohbette --conversation verilmemeli: %v", args)
	}
	if v, _ := argValue(args, "--output-format"); v != "json" {
		t.Errorf("--output-format json bekleniyordu: %v", args)
	}
}

func TestAgyResumesGivenConversation(t *testing.T) {
	bin, af := argsBin(t, `{"conversation_id":"8372b68f","status":"SUCCESS","response":"ok"}`)
	res, err := (&AgyDriver{BinPath: bin}).Run(context.Background(), "8372b68f", "selam", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := argValue(readArgs(t, af), "--conversation"); v != "8372b68f" {
		t.Errorf("--conversation = %q", v)
	}
	if res.SessionID != "8372b68f" {
		t.Errorf("SessionID = %q", res.SessionID)
	}
}

func TestAgyNonJSONOutputFallsBackToRawText(t *testing.T) {
	bin, _ := argsBin(t, "düz metin")
	res, err := (&AgyDriver{BinPath: bin}).Run(context.Background(), "c1", "x", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Text, "düz metin") || res.SessionID != "c1" {
		t.Errorf("res = %+v", res)
	}
}

func TestClaudeNewSessionUsesGeneratedUUID(t *testing.T) {
	bin, af := argsBin(t, "cevap")
	res, err := (&ClaudeDriver{BinPath: bin}).Run(context.Background(), "", "selam", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	args := readArgs(t, af)
	id, ok := argValue(args, "--session-id")
	if !ok || len(id) != 36 || id[14] != '4' {
		t.Fatalf("--session-id UUIDv4 bekleniyordu: %v", args)
	}
	if _, ok := argValue(args, "--resume"); ok {
		t.Errorf("yeni sohbette --resume verilmemeli")
	}
	if res.SessionID != id {
		t.Errorf("SessionID = %q, want %q", res.SessionID, id)
	}
}

func TestClaudeResumesGivenSession(t *testing.T) {
	bin, af := argsBin(t, "cevap")
	res, err := (&ClaudeDriver{BinPath: bin}).Run(context.Background(), "a9ce2f25", "selam", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	args := readArgs(t, af)
	if v, _ := argValue(args, "--resume"); v != "a9ce2f25" {
		t.Errorf("--resume = %q", v)
	}
	if _, ok := argValue(args, "--session-id"); ok {
		t.Errorf("devamda --session-id verilmemeli")
	}
	if res.SessionID != "a9ce2f25" {
		t.Errorf("SessionID = %q", res.SessionID)
	}
}

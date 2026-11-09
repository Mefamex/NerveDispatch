package drivers

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ClaudeDriver manages Claude Code CLI executions.
type ClaudeDriver struct {
	BinPath string
}

// NewClaudeDriver creates a ClaudeDriver locating the claude binary.
func NewClaudeDriver(customBin string) *ClaudeDriver {
	bin := customBin
	if bin == "" {
		home, _ := os.UserHomeDir()
		candidate := filepath.Join(home, ".local/bin/claude")
		if _, err := os.Stat(candidate); err == nil {
			bin = candidate
		} else if resolved, err := exec.LookPath("claude"); err == nil {
			bin = resolved
		} else {
			bin = "claude"
		}
	}
	return &ClaudeDriver{BinPath: bin}
}

func (d *ClaudeDriver) Name() string {
	return "claude"
}

// Run executes claude in print mode. With a sessionID it resumes that
// conversation (--resume); without one it creates a new conversation under a
// freshly generated UUID (--session-id) and returns that UUID for persisting.
func (d *ClaudeDriver) Run(ctx context.Context, sessionID string, prompt string, env []string, dir string) (Result, error) {
	args := []string{"-p", prompt, "--dangerously-skip-permissions"}
	used := sessionID
	if sessionID != "" {
		args = append(args, "--resume", sessionID)
	} else {
		id, err := newUUID()
		if err != nil {
			return Result{}, fmt.Errorf("claude oturum kimliği üretilemedi: %w", err)
		}
		used = id
		args = append(args, "--session-id", id)
	}

	cmd := exec.CommandContext(ctx, d.BinPath, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if len(env) > 0 {
		cmd.Env = env
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	res := Result{SessionID: used}
	if err := cmd.Run(); err != nil {
		res.Text = stdout.String()
		errOut := strings.TrimSpace(stderr.String())
		if errOut != "" {
			return res, fmt.Errorf("claude error: %w (stderr: %s)", err, errOut)
		}
		return res, err
	}
	res.Text = stdout.String()
	return res, nil
}

// newUUID returns a random RFC 4122 version-4 UUID.
func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

package drivers

import (
	"bytes"
	"context"
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

// Run executes claude with session resume support.
func (d *ClaudeDriver) Run(ctx context.Context, sessionID string, prompt string, env []string, dir string) (string, error) {
	args := []string{"-p", prompt, "--dangerously-skip-permissions"}
	if sessionID != "" {
		args = append(args, "--resume", sessionID)
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

	if err := cmd.Run(); err != nil {
		errOut := strings.TrimSpace(stderr.String())
		if errOut != "" {
			return stdout.String(), fmt.Errorf("claude error: %w (stderr: %s)", err, errOut)
		}
		return stdout.String(), err
	}

	return stdout.String(), nil
}

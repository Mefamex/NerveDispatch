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

// AgyDriver manages Antigravity CLI executions.
type AgyDriver struct {
	BinPath string
}

// NewAgyDriver creates an AgyDriver locating the agy binary.
func NewAgyDriver(customBin string) *AgyDriver {
	bin := customBin
	if bin == "" {
		home, _ := os.UserHomeDir()
		candidate := filepath.Join(home, ".local/bin/agy")
		if _, err := os.Stat(candidate); err == nil {
			bin = candidate
		} else if resolved, err := exec.LookPath("agy"); err == nil {
			bin = resolved
		} else {
			bin = "agy"
		}
	}
	return &AgyDriver{BinPath: bin}
}

func (d *AgyDriver) Name() string {
	return "agy"
}

// Run executes agy with optional conversation persistence.
func (d *AgyDriver) Run(ctx context.Context, sessionID string, prompt string, env []string, dir string) (string, error) {
	args := []string{"-p", prompt, "--dangerously-skip-permissions"}
	if sessionID != "" {
		args = append(args, "--conversation", sessionID)
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
			return stdout.String(), fmt.Errorf("agy error: %w (stderr: %s)", err, errOut)
		}
		return stdout.String(), err
	}

	return stdout.String(), nil
}

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

// OpencodeDriver manages OpenCode CLI executions.
type OpencodeDriver struct {
	BinPath string
}

// NewOpencodeDriver creates an OpencodeDriver locating the opencode binary.
func NewOpencodeDriver(customBin string) *OpencodeDriver {
	bin := customBin
	if bin == "" {
		home, _ := os.UserHomeDir()
		candidate := filepath.Join(home, ".opencode/bin/opencode")
		if _, err := os.Stat(candidate); err == nil {
			bin = candidate
		} else if resolved, err := exec.LookPath("opencode"); err == nil {
			bin = resolved
		} else {
			bin = "opencode"
		}
	}
	return &OpencodeDriver{BinPath: bin}
}

func (d *OpencodeDriver) Name() string {
	return "opencode"
}

// Run executes opencode with persistent session support.
func (d *OpencodeDriver) Run(ctx context.Context, sessionID string, prompt string, env []string, dir string) (string, error) {
	args := []string{"run"}
	if sessionID != "" {
		args = append(args, "-s", sessionID)
	}
	args = append(args, prompt)

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
			return stdout.String(), fmt.Errorf("opencode error: %w (stderr: %s)", err, errOut)
		}
		return stdout.String(), err
	}

	return stdout.String(), nil
}

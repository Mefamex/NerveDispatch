package drivers

import (
	"bytes"
	"context"
	"encoding/json"
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

// agyJSON is the subset of `agy -p --output-format json` we rely on.
type agyJSON struct {
	ConversationID string `json:"conversation_id"`
	Status         string `json:"status"`
	Response       string `json:"response"`
}

// Run executes agy in print mode. With a sessionID it resumes that
// conversation (--conversation); without one agy opens a new conversation and
// its ID is read back from the JSON output so it can be persisted.
func (d *AgyDriver) Run(ctx context.Context, sessionID string, prompt string, env []string, dir string) (Result, error) {
	args := []string{"-p", prompt, "--dangerously-skip-permissions", "--output-format", "json"}
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
	runErr := cmd.Run()

	res := Result{Text: stdout.String(), SessionID: sessionID}
	var parsed agyJSON
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &parsed); err == nil {
		res.Text = parsed.Response
		if parsed.ConversationID != "" {
			res.SessionID = parsed.ConversationID
		}
	}

	if runErr != nil {
		errOut := strings.TrimSpace(stderr.String())
		if errOut != "" {
			return res, fmt.Errorf("agy error: %w (stderr: %s)", runErr, errOut)
		}
		return res, runErr
	}
	return res, nil
}

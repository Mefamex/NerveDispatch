// Package drivers provides unified adapter interfaces for AI coding agents
// (opencode, agy, claude) with persistent session support.
package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Session represents a conversation session metadata.
type Session struct {
	ID        string `json:"id"`
	Title     string `json:"title,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Driver is the common interface implemented by all AI agent CLI adapters.
type Driver interface {
	Name() string
	Run(ctx context.Context, sessionID string, prompt string, env []string, dir string) (string, error)
}

// SessionEntry stores the persistent session mapping in oturumlar.json.
type SessionEntry struct {
	SessionID string `json:"session_id,omitempty"`
	ConvID    string `json:"conversation_id,omitempty"` // For agy
	UpdatedAt string `json:"updated_at,omitempty"`
}

// SessionStore handles thread-safe persistence of agent sessions.
type SessionStore struct {
	path string
	mu   sync.Mutex
}

// NewSessionStore creates or loads the session store at the given path.
func NewSessionStore(stateDir string) *SessionStore {
	return &SessionStore{
		path: filepath.Join(stateDir, "oturumlar.json"),
	}
}

// Get returns the persistent session ID for the given agent driver.
func (s *SessionStore) Get(driverName string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		return ""
	}

	var m map[string]SessionEntry
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}

	entry, ok := m[driverName]
	if !ok {
		return ""
	}
	if entry.SessionID != "" {
		return entry.SessionID
	}
	return entry.ConvID
}

// Set saves or updates the session ID for the given agent driver.
func (s *SessionStore) Set(driverName, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := make(map[string]SessionEntry)
	if data, err := os.ReadFile(s.path); err == nil {
		_ = json.Unmarshal(data, &m)
	}

	now := time.Now().Format(time.RFC3339)
	entry := m[driverName]
	if driverName == "agy" {
		entry.ConvID = id
	} else {
		entry.SessionID = id
	}
	entry.UpdatedAt = now
	m[driverName] = entry

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal oturumlar.json: %w", err)
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.path)
}

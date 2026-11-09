package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// QuotaStatus represents agent-specific quota limits.
type QuotaStatus struct {
	Agent       string    `json:"agent"`        // "agy", "claude", "opencode"
	FiveHourPct float64   `json:"five_hour_pct"`// 0.0 - 100.0
	WeeklyPct   float64   `json:"weekly_pct"`   // 0.0 - 100.0
	ResetAt     time.Time `json:"reset_at"`     // Kota sıfırlanma zamanı
	Blocked     bool      `json:"blocked"`      // %95 eşiğini aştı mı
}

// QuotaGuard enforces the 95% threshold safety rule for agents.
type QuotaGuard struct {
	mu        sync.RWMutex
	stateFile string
	statuses  map[string]QuotaStatus
}

const QuotaThresholdPct = 95.0

// NewQuotaGuard initializes the quota guard using the project state directory.
func NewQuotaGuard(stateDir string) *QuotaGuard {
	if stateDir == "" {
		stateDir = defaultState("")
	}
	f := filepath.Join(stateDir, "quota_status.json")
	g := &QuotaGuard{
		stateFile: f,
		statuses:  make(map[string]QuotaStatus),
	}
	g.load()
	return g
}

func (g *QuotaGuard) load() {
	data, err := os.ReadFile(g.stateFile)
	if err != nil {
		return
	}
	var m map[string]QuotaStatus
	if err := json.Unmarshal(data, &m); err == nil {
		g.statuses = m
	}
}

func (g *QuotaGuard) save() error {
	data, err := json.MarshalIndent(g.statuses, "", "  ")
	if err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(g.stateFile), 0o755)
	return os.WriteFile(g.stateFile, data, 0o644)
}

// UpdateQuota records live quota levels for an agent and determines if it must be blocked.
func (g *QuotaGuard) UpdateQuota(agent string, fiveHourPct, weeklyPct float64, resetAt time.Time) QuotaStatus {
	g.mu.Lock()
	defer g.mu.Unlock()

	blocked := fiveHourPct >= QuotaThresholdPct || weeklyPct >= QuotaThresholdPct

	// Eğer resetAt geçmişte kalmışsa blokajı kaldır
	if blocked && !resetAt.IsZero() && time.Now().After(resetAt) {
		blocked = false
	}

	st := QuotaStatus{
		Agent:       agent,
		FiveHourPct: fiveHourPct,
		WeeklyPct:   weeklyPct,
		ResetAt:     resetAt,
		Blocked:     blocked,
	}
	g.statuses[agent] = st
	_ = g.save()
	return st
}

// CanExecute checks if the given agent is safe to run.
func (g *QuotaGuard) CanExecute(agent string) (bool, string) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	st, exists := g.statuses[agent]
	if !exists {
		return true, ""
	}

	// Yenilenme süresi dolmuşsa otomatik aç
	if st.Blocked && !st.ResetAt.IsZero() && time.Now().After(st.ResetAt) {
		return true, ""
	}

	if st.Blocked {
		kalan := time.Until(st.ResetAt).Round(time.Minute)
		return false, fmt.Sprintf("%s kotası %%%.1f seviyesinde (%%95 sınırı aşıldı). Yenilenmeye kalan: %v", agent, st.FiveHourPct, kalan)
	}

	return true, ""
}

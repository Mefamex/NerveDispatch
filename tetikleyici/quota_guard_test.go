package main

import (
	"testing"
	"time"
)

func TestQuotaGuardNormal(t *testing.T) {
	tmpDir := t.TempDir()
	guard := NewQuotaGuard(tmpDir)

	guard.UpdateQuota("agy", 50.0, 40.0, time.Now().Add(2*time.Hour))
	ok, reason := guard.CanExecute("agy")
	if !ok {
		t.Fatalf("beklenen ok=true, got reason: %s", reason)
	}
}

func TestQuotaGuardBreach(t *testing.T) {
	tmpDir := t.TempDir()
	guard := NewQuotaGuard(tmpDir)

	resetTime := time.Now().Add(45 * time.Minute)
	guard.UpdateQuota("claude", 95.5, 30.0, resetTime)

	ok, reason := guard.CanExecute("claude")
	if ok {
		t.Fatal("beklenen %95 üzeri kotada bloklanmasıydı (ok=false)")
	}
	if reason == "" {
		t.Error("bloklanma sebebi boş dönmemeli")
	}
}

func TestQuotaGuardWeeklyBreach(t *testing.T) {
	tmpDir := t.TempDir()
	guard := NewQuotaGuard(tmpDir)

	guard.UpdateQuota("opencode", 40.0, 96.0, time.Now().Add(24*time.Hour))
	ok, _ := guard.CanExecute("opencode")
	if ok {
		t.Fatal("haftalık kota %95 üstündeyken engellenmeliydi")
	}
}

func TestQuotaGuardAutoRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	guard := NewQuotaGuard(tmpDir)

	// Geçmişte kalmış reset süresi
	pastReset := time.Now().Add(-1 * time.Minute)
	guard.UpdateQuota("agy", 98.0, 20.0, pastReset)

	ok, _ := guard.CanExecute("agy")
	if !ok {
		t.Fatal("süresi dolmuş kotanın otomatik açılması bekleniyordu")
	}
}

func TestQuotaGuardPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	guard1 := NewQuotaGuard(tmpDir)

	guard1.UpdateQuota("claude", 97.0, 10.0, time.Now().Add(1*time.Hour))

	// Yeni instance dosyadaki durumu okumalı
	guard2 := NewQuotaGuard(tmpDir)
	ok, _ := guard2.CanExecute("claude")
	if ok {
		t.Fatal("dosyadan okunan durumun bloke olarak korunması bekleniyordu")
	}
}

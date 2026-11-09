package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tetikleyici/drivers"
)

// sahteSurucu, aldığı oturum kimliğini kaydeder ve verilen kimliği döndürür.
type sahteSurucu struct {
	ad       string
	gelenID  string
	donenID  string
	cagrildi int
}

func (s *sahteSurucu) Name() string { return s.ad }
func (s *sahteSurucu) Run(_ context.Context, sessionID, _ string, _ []string, _ string) (drivers.Result, error) {
	s.cagrildi++
	s.gelenID = sessionID
	id := s.donenID
	if id == "" {
		id = sessionID
	}
	return drivers.Result{Text: "cevap", SessionID: id}, nil
}

// motorOrtami, geçici durum dizini kurar, aktif motoru yazar ve sahte sürücüyü bağlar.
func motorOrtami(t *testing.T, motor string, s *sahteSurucu) *drivers.SessionStore {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("TETIKLEYICI_STATE_DIR", dir)
	if err := os.WriteFile(filepath.Join(dir, "aktif_motor.txt"), []byte(motor+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	eski := yeniSurucu
	yeniSurucu = func(string) drivers.Driver { return s }
	t.Cleanup(func() { yeniSurucu = eski })
	return drivers.NewSessionStore(dir)
}

func TestKuraldakiOpencodeKimligiAgyyeGitmez(t *testing.T) {
	s := &sahteSurucu{ad: "agy", donenID: "agy-yeni"}
	store := motorOrtami(t, "agy", s)

	act := action{Kind: "session", SessionID: "ses_ef8590662ffeRSLXLL2R8VhJfM"}
	if _, err := runAgent(context.Background(), testOptions(), act, testPayload(), nil); err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(s.gelenID, "ses_") {
		t.Fatalf("agy'ye OpenCode kimliği gitti: %q", s.gelenID)
	}
	if s.gelenID != "" {
		t.Fatalf("kayıt yokken yeni sohbet (boş kimlik) bekleniyordu, gelen %q", s.gelenID)
	}
	if got := store.Get("agy"); got != "agy-yeni" {
		t.Fatalf("yeni agy sohbeti kaydedilmedi: %q", got)
	}
}

func TestKayitliSohbetVarsaOnaBaglanir(t *testing.T) {
	s := &sahteSurucu{ad: "claude"}
	store := motorOrtami(t, "claude", s)
	if err := store.Set("claude", "a9ce2f25-b315-46bd-ad5b-b3434c9246d1"); err != nil {
		t.Fatal(err)
	}

	act := action{Kind: "session", SessionID: "ses_kural"}
	for i := 0; i < 2; i++ {
		if _, err := runAgent(context.Background(), testOptions(), act, testPayload(), nil); err != nil {
			t.Fatal(err)
		}
		if s.gelenID != "a9ce2f25-b315-46bd-ad5b-b3434c9246d1" {
			t.Fatalf("tur %d: kayıtlı claude sohbeti kullanılmadı: %q", i, s.gelenID)
		}
	}
	if store.Get("claude") != "a9ce2f25-b315-46bd-ad5b-b3434c9246d1" {
		t.Fatal("kayıtlı sohbet değişmemeliydi")
	}
}

func TestYeniAcilanSohbetSonrakiMesajdaKullanilir(t *testing.T) {
	s := &sahteSurucu{ad: "agy", donenID: "ilk-sohbet"}
	motorOrtami(t, "agy", s)

	act := action{Kind: "session"}
	if _, err := runAgent(context.Background(), testOptions(), act, testPayload(), nil); err != nil {
		t.Fatal(err)
	}
	s.donenID = ""
	if _, err := runAgent(context.Background(), testOptions(), act, testPayload(), nil); err != nil {
		t.Fatal(err)
	}
	if s.gelenID != "ilk-sohbet" {
		t.Fatalf("ikinci mesaj ilk açılan sohbete gitmeliydi, gelen %q", s.gelenID)
	}
}

func TestMotorlarBirbirininSohbetiniKullanmaz(t *testing.T) {
	s := &sahteSurucu{ad: "claude"}
	store := motorOrtami(t, "claude", s)
	if err := store.Set("agy", "8372b68f-a31b-487d-bef9-144284811a30"); err != nil {
		t.Fatal(err)
	}
	s.donenID = "claude-yeni"
	if _, err := runAgent(context.Background(), testOptions(), action{Kind: "session"}, testPayload(), nil); err != nil {
		t.Fatal(err)
	}
	if s.gelenID != "" {
		t.Fatalf("claude, agy'nin sohbetini almamalı: %q", s.gelenID)
	}
	if store.Get("agy") != "8372b68f-a31b-487d-bef9-144284811a30" || store.Get("claude") != "claude-yeni" {
		t.Fatalf("kayıtlar karıştı: agy=%q claude=%q", store.Get("agy"), store.Get("claude"))
	}
}

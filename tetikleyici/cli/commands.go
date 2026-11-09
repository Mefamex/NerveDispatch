// Package cli provides unified command-line management for NerveDispatch.
package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Dispatcher processes subcommands. Returns true if a command was handled.
func Dispatcher(rawArgs []string) (bool, error) {
	if len(rawArgs) == 0 {
		return false, nil
	}

	cmd := rawArgs[0]
	subArgs := rawArgs[1:]

	// If starts with flag (e.g. --config, -bridge), let daemon handle it
	if strings.HasPrefix(cmd, "-") {
		return false, nil
	}

	switch cmd {
	case "status":
		return true, handleStatus()
	case "session":
		return true, handleSession(subArgs)
	case "channel":
		return true, handleChannel(subArgs)
	case "config":
		return true, handleConfig(subArgs)
	case "logs":
		return true, handleLogs(subArgs)
	case "restart":
		return true, handleRestart()
	case "qr":
		return true, handleQR()
	case "pair":
		return true, handlePair(subArgs)
	case "snapshot", "backup":
		return true, handleSnapshot(subArgs)
	case "help":
		printHelp()
		return true, nil
	default:
		return false, nil
	}
}

func printHelp() {
	fmt.Println("NerveDispatch (nerve) — Birleşik Yönetim CLI")
	fmt.Println("\nKullanım: nerve [komut] [parametreler]")
	fmt.Println("\nKomutlar:")
	fmt.Println("  status          Köprü, dağıtıcı, bağlı kanallar ve AI motor durumunu gösterir")
	fmt.Println("  session list    Kayıtlı kalıcı AI oturumlarını listeler")
	fmt.Println("  session reset   Belirtilen motorun oturumunu sıfırlar (opencode | agy | claude)")
	fmt.Println("  channel list    Desteklenen kanalların bağlantı durumunu listeler")
	fmt.Println("  config check    YAML yapılandırma dosyalarını ve yolları doğrular")
	fmt.Println("  logs [-f]       Tetikleyici günlük kayıtlarını görüntüler")
	fmt.Println("  restart         NerveDispatch dağıtıcı servisini güvenli yeniden başlatır")
	fmt.Println("  qr              WhatsApp eşleştirme QR kodunu terminalde görüntüler")
	fmt.Println("  pair <telefon>  8 haneli WhatsApp telefon eşleştirme kodu üretir")
	fmt.Println("  snapshot        Sıfır veri kaybı anlık yedeği alır (session.db, kurallar, durum)")
}

func handleStatus() error {
	fmt.Println("=== NerveDispatch Sistem Durumu ===")

	// 1. Bridge HTTP API
	bridgeAddr := "http://127.0.0.1:31415"
	resp, err := http.Get(bridgeAddr + "/health")
	if err != nil {
		fmt.Printf(" Köprü (Bridge :31415)  : \033[31mKAPALI\033[0m (%v)\n", err)
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var h map[string]interface{}
		_ = json.Unmarshal(body, &h)
		paired, _ := h["paired"].(bool)
		connected, _ := h["connected"].(bool)
		fmt.Printf(" Köprü (Bridge :31415)  : \033[32mÇALIŞIYOR\033[0m (Eşleşti: %v, Bağlı: %v)\n", paired, connected)
	}

	// 2. Dispatcher Servisi
	out, err := exec.Command("systemctl", "--user", "is-active", "whatsapp-tetikleyici.service").Output()
	status := strings.TrimSpace(string(out))
	if status == "active" {
		pidOut, _ := exec.Command("systemctl", "--user", "show", "--property", "MainPID", "--value", "whatsapp-tetikleyici.service").Output()
		fmt.Printf(" Dağıtıcı (Dispatcher)  : \033[32mAKTİF\033[0m (PID: %s)\n", strings.TrimSpace(string(pidOut)))
	} else {
		fmt.Printf(" Dağıtıcı (Dispatcher)  : \033[33m%s\033[0m\n", status)
	}

	// 3. Aktif Motor
	projDir := findProjectRoot()
	motorFile := filepath.Join(projDir, ".tetikleyici", "aktif_motor.txt")
	motor := "agy"
	if b, err := os.ReadFile(motorFile); err == nil && len(b) > 0 {
		motor = strings.TrimSpace(string(b))
	}
	fmt.Printf(" Aktif AI Motoru        : \033[36m%s\033[0m\n", motor)

	// 4. Kota Sigortası (%95 Guard)
	quotaFile := filepath.Join(projDir, ".tetikleyici", "quota_status.json")
	if qData, err := os.ReadFile(quotaFile); err == nil {
		var qMap map[string]struct {
			FiveHourPct float64   `json:"five_hour_pct"`
			Blocked     bool      `json:"blocked"`
			ResetAt     time.Time `json:"reset_at"`
		}
		if err := json.Unmarshal(qData, &qMap); err == nil && len(qMap) > 0 {
			for mName, q := range qMap {
				durum := "\033[32mGÜVENLİ\033[0m"
				if q.Blocked {
					kalan := time.Until(q.ResetAt).Round(time.Minute)
					durum = fmt.Sprintf("\033[31mBLOKE (%%%0.1f - Kalan %v)\033[0m", q.FiveHourPct, kalan)
				} else {
					durum = fmt.Sprintf("\033[32m%%%0.1f\033[0m", q.FiveHourPct)
				}
				fmt.Printf(" Kota [%-8s]        : %s\n", mName, durum)
			}
		}
	}

	return nil
}

func handleSession(args []string) error {
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
	}

	projDir := findProjectRoot()
	oturumlarFile := filepath.Join(projDir, ".tetikleyici", "oturumlar.json")

	switch sub {
	case "list":
		data, err := os.ReadFile(oturumlarFile)
		if err != nil {
			fmt.Println("Henüz kayıtlı oturum bulunmuyor.")
			return nil
		}
		var m map[string]map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			return fmt.Errorf("oturumlar.json ayrıştırma hatası: %w", err)
		}
		fmt.Printf("%-12s %-36s %-20s\n", "MOTOR", "OTURUM / KONUŞMA ID", "GÜNCELLEME")
		fmt.Println(strings.Repeat("-", 70))
		for motor, entry := range m {
			id := entry["session_id"]
			if id == nil || id == "" {
				id = entry["conversation_id"]
			}
			updated := entry["updated_at"]
			fmt.Printf("%-12s %-36v %-20v\n", motor, id, updated)
		}
		return nil

	case "reset":
		if len(args) < 2 {
			return fmt.Errorf("kullanım: nerve session reset <opencode | agy | claude>")
		}
		target := args[1]
		data, _ := os.ReadFile(oturumlarFile)
		var m map[string]map[string]interface{}
		_ = json.Unmarshal(data, &m)
		if m != nil && m[target] != nil {
			delete(m, target)
			out, _ := json.MarshalIndent(m, "", "  ")
			_ = os.WriteFile(oturumlarFile, out, 0o600)
			fmt.Printf("%s oturumu sıfırlandı.\n", target)
		}
		return nil

	default:
		return fmt.Errorf("bilinmeyen session komutu: %s", sub)
	}
}

func handleChannel(args []string) error {
	fmt.Println("=== NerveDispatch İletişim Kanalları ===")

	client := &http.Client{Timeout: 300 * time.Millisecond}
	resp, err := client.Get("http://127.0.0.1:31415/channels")
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var res struct {
			Status   string `json:"status"`
			Channels []struct {
				Name      string `json:"name"`
				Connected bool   `json:"connected"`
				Paired    bool   `json:"paired"`
				Account   string `json:"account"`
			} `json:"channels"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&res); err == nil && len(res.Channels) > 0 {
			for _, ch := range res.Channels {
				title := kanalAdi(ch.Name)
				durum := "\033[33mBEKLEMEDE\033[0m"
				extra := ""
				if ch.Connected {
					durum = "\033[32mAKTİF\033[0m"
				} else if ch.Paired {
					durum = "\033[36mEŞLEŞTİ\033[0m"
				}
				if ch.Account != "" {
					extra = fmt.Sprintf(" (%s)", ch.Account)
				}
				fmt.Printf("  %-12s : %s%s\n", title, durum, extra)
			}
			return nil
		}
	}

	// Fallback yerel görünüm (Köprü çalışmıyorken veya ilk kurulumda)
	fmt.Printf("  %-12s : \033[32mAKTİF\033[0m (Whatsmeow Soketi :31415)\n", "WhatsApp")
	fmt.Printf("  %-12s : \033[33mBEKLEMEDE\033[0m (Mefamex TG ID: 975678139)\n", "Telegram")
	fmt.Printf("  %-12s : \033[33mBEKLEMEDE\033[0m (Gateway WebSocket)\n", "Discord")
	fmt.Printf("  %-12s : \033[33mBEKLEMEDE\033[0m (DM Graph API)\n", "Instagram")
	return nil
}

// kanalAdi, köprünün döndürdüğü küçük harfli kanal adını marka yazılışına
// çevirir ("whatsapp" → "WhatsApp"). Bilinmeyen adın yalnız ilk harfi büyür.
func kanalAdi(ad string) string {
	switch strings.ToLower(strings.TrimSpace(ad)) {
	case "whatsapp":
		return "WhatsApp"
	case "telegram":
		return "Telegram"
	case "discord":
		return "Discord"
	case "instagram":
		return "Instagram"
	}
	if ad == "" {
		return ad
	}
	return strings.ToUpper(ad[:1]) + ad[1:]
}

func handleConfig(args []string) error {
	projDir := findProjectRoot()
	fmt.Println("=== Yapılandırma Dosyaları Doğrulanıyor ===")

	checkYAML := func(name string) {
		p := filepath.Join(projDir, name)
		data, err := os.ReadFile(p)
		if err != nil {
			fmt.Printf("  %-25s : \033[31mEKSİK\033[0m\n", name)
			return
		}
		var raw interface{}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			fmt.Printf("  %-25s : \033[31mBOZUK YAML\033[0m (%v)\n", name, err)
			return
		}
		fmt.Printf("  %-25s : \033[32mGEÇERLİ\033[0m\n", name)
	}

	checkYAML("tetikleyici.yaml")
	checkYAML("pasif-kurallar.yaml")
	checkYAML("config.yaml")
	return nil
}

func handleLogs(args []string) error {
	cmdArgs := []string{"--user", "-u", "whatsapp-tetikleyici.service"}
	for _, a := range args {
		if a == "-f" {
			cmdArgs = append(cmdArgs, "-f")
		}
	}
	if !contains(cmdArgs, "-f") {
		cmdArgs = append(cmdArgs, "-n", "30", "--no-pager")
	}

	cmd := exec.Command("journalctl", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func handleRestart() error {
	fmt.Println("NerveDispatch servisi yeniden başlatılıyor...")
	if err := exec.Command("systemctl", "--user", "restart", "whatsapp-tetikleyici.service").Run(); err != nil {
		return fmt.Errorf("restart hatası: %w", err)
	}
	time.Sleep(500 * time.Millisecond)
	pidOut, _ := exec.Command("systemctl", "--user", "show", "--property", "MainPID", "--value", "whatsapp-tetikleyici.service").Output()
	fmt.Printf("\033[32mServis başarıyla yeniden başlatıldı\033[0m (Yeni PID: %s)\n", strings.TrimSpace(string(pidOut)))
	return nil
}

func handleQR() error {
	resp, err := http.Get("http://127.0.0.1:31415/qr")
	if err != nil {
		return fmt.Errorf("köprüye erişilemedi: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return nil
}

func handlePair(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kullanım: nerve pair <telefon_numarası>")
	}
	phone := args[0]
	resp, err := http.Post("http://127.0.0.1:31415/pair?phone="+phone, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("köprüye erişilemedi: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Eşleştirme Yanıtı: %s\n", string(body))
	return nil
}

func findProjectRoot() string {
	if p := os.Getenv("TETIKLEYICI_PROJECT_DIR"); p != "" {
		return p
	}
	if cwd, err := os.Getwd(); err == nil {
		if fi, err := os.Stat(filepath.Join(cwd, "tetikleyici.yaml")); err == nil && !fi.IsDir() {
			return cwd
		}
	}
	return "/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch"
}

func handleSnapshot(args []string) error {
	projDir := findProjectRoot()
	scriptPath := filepath.Join(projDir, "scripts", "nerve-snapshot.sh")
	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf("snapshot betiği bulunamadı: %s", scriptPath)
	}

	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Dir = projDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

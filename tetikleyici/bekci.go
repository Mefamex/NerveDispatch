// Package main provides proactive system monitoring (bekci).
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// bekciRunner performs periodic disk, memory, temperature and bridge health checks.
type bekciRunner struct {
	client    *bridgeClient
	targets   []string
	interval  time.Duration
	cooldown  time.Duration
	mu        sync.Mutex
	sonAlarm  map[string]time.Time
}

func newBekci(client *bridgeClient, targets []string) *bekciRunner {
	return &bekciRunner{
		client:   client,
		targets:  targets,
		interval: 5 * time.Minute,
		cooldown: 1 * time.Hour,
		sonAlarm: make(map[string]time.Time),
	}
}

// start runs the watchdog loop until ctx is cancelled.
func (b *bekciRunner) start(ctx context.Context) {
	if len(b.targets) == 0 || b.client == nil {
		return
	}

	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.kontrolEt(ctx)
		}
	}
}

func (b *bekciRunner) kontrolEt(ctx context.Context) {
	// 1. Disk Doluluk Denetimi (> %90)
	for _, mount := range []string{"/", "/mnt/DiskD"} {
		if doluluk, err := diskDolulukOrani(mount); err == nil && doluluk >= 90 {
			b.alarmUret(ctx, "disk_"+mount, fmt.Sprintf("⚠️ Disk Doluluk Alarmı: `%s` bölümü %%%d dolu!", mount, doluluk))
		}
	}

	// 2. RAM Doluluk Denetimi (> %95)
	if ramDoluluk, err := ramDolulukOrani(); err == nil && ramDoluluk >= 95 {
		b.alarmUret(ctx, "ram", fmt.Sprintf("⚠️ RAM Tüketim Alarmı: Bellek doluluğu %%%d seviyesinde!", ramDoluluk))
	}

	// 3. CPU Sıcaklık Denetimi (> 85°C)
	if sicaklik, err := cpuSicaklik(); err == nil && sicaklik >= 85 {
		b.alarmUret(ctx, "temp", fmt.Sprintf("🔥 CPU Sıcaklık Alarmı: İşlemci sıcaklığı %d°C!", sicaklik))
	}
}

func (b *bekciRunner) alarmUret(ctx context.Context, kategori, mesaj string) {
	b.mu.Lock()
	son, varMi := b.sonAlarm[kategori]
	if varMi && time.Since(son) < b.cooldown {
		b.mu.Unlock()
		return
	}
	b.sonAlarm[kategori] = time.Now()
	b.mu.Unlock()

	logf("[BEKÇİ ALARMI] %s", mesaj)
	sendCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	uyariMetni := fmt.Sprintf("🚨 *[NerveDispatch Sistem Bekçisi]*\n\n%s\n\n_Zaman: %s_",
		mesaj, time.Now().Format("15:04:05"))
	if b.client != nil {
		for _, to := range b.targets {
			_ = b.client.sendText(sendCtx, to, uyariMetni)
		}
	}
}

func diskDolulukOrani(path string) (int, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	if total == 0 {
		return 0, fmt.Errorf("sıfır boyut")
	}
	used := total - free
	return int((used * 100) / total), nil
}

func ramDolulukOrani() (int, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var total, avail uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			total = parseKb(line)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			avail = parseKb(line)
		}
	}
	if total == 0 {
		return 0, fmt.Errorf("meminfo okunamadı")
	}
	used := total - avail
	return int((used * 100) / total), nil
}

func parseKb(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) >= 2 {
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		return val
	}
	return 0
}

func cpuSicaklik() (int, error) {
	// /sys/class/thermal/thermal_zone*/temp dosyalarından en yüksek sıcaklığı al
	zones, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil || len(zones) == 0 {
		return 0, fmt.Errorf("thermal zone bulunamadı")
	}

	maxTemp := 0
	for _, z := range zones {
		data, err := os.ReadFile(z)
		if err != nil {
			continue
		}
		raw, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil {
			celsius := raw / 1000
			if celsius > maxTemp && celsius < 130 { // 130 üzeri sahte donanım değerlerini ele
				maxTemp = celsius
			}
		}
	}
	if maxTemp == 0 {
		return 0, fmt.Errorf("geçerli sıcaklık okunamadı")
	}
	return maxTemp, nil
}

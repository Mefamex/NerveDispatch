package main

import (
	"context"
	"testing"
	"time"
)

func TestBekciMetricCollection(t *testing.T) {
	// 1. Disk doluluk testi
	orani, err := diskDolulukOrani("/")
	if err != nil {
		t.Logf("disk doluluk okunamadı (CI ortamı olabilir): %v", err)
	} else {
		if orani < 0 || orani > 100 {
			t.Errorf("geçersiz disk doluluk oranı: %d", orani)
		}
	}

	// 2. RAM doluluk testi
	ramOrani, err := ramDolulukOrani()
	if err != nil {
		t.Logf("ram doluluk okunamadı: %v", err)
	} else {
		if ramOrani < 0 || ramOrani > 100 {
			t.Errorf("geçersiz ram doluluk oranı: %d", ramOrani)
		}
	}

	// 3. Bekçi nesnesi ve cooldown denetimi
	bekci := newBekci(nil, []string{"test@s.whatsapp.net"})
	if bekci.cooldown != 1*time.Hour {
		t.Errorf("beklenen cooldown 1 saat, got %v", bekci.cooldown)
	}

	// Alarm simülasyonu
	bekci.alarmUret(context.Background(), "test_kat", "test mesajı")
	if len(bekci.sonAlarm) != 1 {
		t.Errorf("alarm kaydedilmedi")
	}

	// İkinci çağrıda cooldown engeli
	eskiZaman := bekci.sonAlarm["test_kat"]
	bekci.alarmUret(context.Background(), "test_kat", "ikinci mesaj")
	if bekci.sonAlarm["test_kat"] != eskiZaman {
		t.Errorf("cooldown sırasında alarm zamanı güncellenmemeliydi")
	}
}

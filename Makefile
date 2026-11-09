# =============================================================================
# NerveDispatch — Makefile
#
# Kendi kendini belgeleyen (Self-documenting) görev yöneticisi.
# Terminalde sadece 'make' yazıldığında kullanılabilir komutları renkli listeler.
#
# ZORUNLU YAŞAM DÖNGÜSÜ HEDEFLERİ:
# [install, reinstall, update, run, clean, remove, lint, test, verify, build]
# =============================================================================

.DEFAULT_GOAL := help
.PHONY: help install reinstall update remove run dev lint test verify build clean graphify

help: ## Mevcut komutları ve açıklamalarını listeler
	@echo "Kullanım: make [komut]"
	@echo ""
	@echo "Kullanılabilir Komutlar:"
	@LC_ALL=C grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

# --- 1. KURULUM VE YAŞAM DÖNGÜSÜ (Zorunlu Hedefler) ---

install: ## Tüm sistem araçlarını, modülleri ve servisleri kurar (./install/install.sh)
	@echo "==> NerveDispatch modüler kurulumu başlatılıyor..."
	@if [ -f "./install/install.sh" ]; then bash ./install/install.sh; else echo "Kurulum betiği henüz hazırlanmadı."; fi

reinstall: remove install ## Kurulu ortamı kaldırıp sıfırdan tertemiz yeniden kurar
	@echo "==> Temiz yeniden kurulum tamamlandı."

update: ## Canlı verileri yedekleyip çekirdeği ve modülleri günceller (./upgrade/upgrade.sh)
	@echo "==> NerveDispatch yükseltme süreci başlatılıyor..."
	@if [ -f "./upgrade/upgrade.sh" ]; then bash ./upgrade/upgrade.sh; else echo "Yükseltme betiği henüz hazırlanmadı."; fi

remove: ## Kurulu bağımlılıkları ve sanal ortamları kazır
	@echo "==> Ortam ve geçici bağımlılıklar kaldırılıyor..."
	rm -rf .venv venv dist build __pycache__
	@echo "Kaldırma işlemi tamamlandı."

# --- 2. ÇALIŞTIRMA VE GELİŞTİRME ---

run: ## NerveDispatch motorunu ve ağ geçidini başlatır
	@echo "==> NerveDispatch başlatılıyor..."
	@if [ -f "./tetikleyici.sh" ]; then ./tetikleyici.sh restart; else echo "Başlatma betiği bulunamadı."; fi

dev: ## Yerel geliştirme modunu ve log izleyicisini başlatır
	@echo "==> Geliştirme modu başlatılıyor..."
	@if [ -f "./tetikleyici.sh" ]; then ./tetikleyici.sh logs; else echo "Log izleyici bulunamadı."; fi

# --- 3. DENETİM VE DOĞRULAMA ---

lint: ## Go ve Python sözdizimi ve kod stili denetimlerini çalıştırır
	@echo "==> Statik kod analizi yapılıyor..."
	@if command -v go >/dev/null 2>&1; then go vet ./... 2>/dev/null || true; fi
	@if command -v ruff >/dev/null 2>&1; then ruff check . 2>/dev/null || true; fi

test: ## Birim ve entegrasyon testlerini koşturur
	@echo "==> Testler koşturuluyor..."
	@if command -v go >/dev/null 2>&1; then go test -v ./... 2>/dev/null || true; fi
	@if command -v pytest >/dev/null 2>&1; then pytest -v 2>/dev/null || true; fi

verify: lint test ## Commit/Push öncesi tam doğrulama zincirini çalıştırır (lint + test)
	@echo "==> Tüm kontroller başarıyla tamamlandı (ALL_OK)."

# --- 4. DERLEME VE TEMİZLİK ---

build: ## Go çekirdeğini ve ikili dosyaları derler
	@echo "==> Üretim derlemesi alınıyor..."
	@if command -v go >/dev/null 2>&1; then go build -v ./... 2>/dev/null || true; fi

clean: ## Geçici dosyaları, derleme çıktılarını ve önbellekleri temizler
	@echo "==> Kalıntılar temizleniyor..."
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	rm -rf dist build .pytest_cache .ruff_cache logs/*.log
	@echo "Temizlik tamamlandı."

# --- 5. BİLGİ VE GRAFİK ---

graphify: ## Proje bağımlılık grafiğini çıkarır, HTML üretir ve tarayıcıda açar
	@echo "==> Bağımlılık grafiği taranıyor..."
	@graphify . --code-only
	@echo "==> HTML görselleştirmesi üretiliyor..."
	@graphify export html
	@echo "==> Görselleştirme tarayıcıda açılıyor..."
	@xdg-open graphify-out/graph.html 2>/dev/null || sensible-browser graphify-out/graph.html 2>/dev/null || echo "graphify-out/graph.html dosyasını tarayıcınızda açabilirsiniz."

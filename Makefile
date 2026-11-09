# =============================================================================
# NerveDispatch — İnce Görev Yönlendiricisi (Thin Dispatcher Makefile)
# =============================================================================

.DEFAULT_GOAL := help
.PHONY: help install reinstall update upgrade remove run dev lint test verify build clean graphify snapshot

help: ## Mevcut komutları ve açıklamalarını listeler
	@echo "Kullanım: make [komut]"
	@echo ""
	@echo "Kullanılabilir Komutlar:"
	@LC_ALL=C grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

# --- 1. YAŞAM DÖNGÜSÜ & ORKESTRASYON ---

install: ## Tüm NerveDispatch ekosistemini modüler olarak kurar (ARGS="--only ..." destekler)
	@./install/install.sh $(ARGS)

upgrade: ## Canlı sistemi yedek alarak sıfır kesintiyle yükseltir (ARGS="--only ..." destekler)
	@./upgrade/upgrade.sh $(ARGS)

update: upgrade ## 'make upgrade' için uyumluluk takma adı

reinstall: ## Kurulu ortamı temizleyip sıfırdan kurar
	@./install/install.sh --reinstall $(ARGS)

remove: ## Derleme çıktılarını ve sanal ortamları kaldırır (kişisel verileri korur)
	@./install/install.sh --remove $(ARGS)

# --- 2. ÇALIŞTIRMA VE GELİŞTİRME ---

run: ## NerveDispatch motorunu ve ağ geçidini başlatır
	@if [ -f "./tetikleyici.sh" ]; then ./tetikleyici.sh restart; else ./bin/tetikleyici; fi

dev: ## Log izleyicisini başlatır
	@if [ -f "./tetikleyici.sh" ]; then ./tetikleyici.sh logs; else journalctl --user -fu whatsapp-tetikleyici; fi

# --- 3. DERLEME, DENETİM VE TEST ---

build: ## Go çekirdeğini ve ikililerini derler
	@mkdir -p bin
	CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o bin/whatsapp-bridge ./bridge
	go build -trimpath -ldflags="-s -w" -o bin/tetikleyici ./tetikleyici
	@ln -sf tetikleyici bin/nerve

lint: ## Go ve Python statik analizi yapar
	@cd bridge && go vet ./...
	@cd tetikleyici && go vet ./...
	@if command -v ruff >/dev/null 2>&1; then ruff check .; fi

test: ## Go ve Python testlerini koşturur (hataları yutmaz)
	@echo "==> Go Bridge Testleri..."
	@cd bridge && go test ./...
	@echo "==> Go Tetikleyici Testleri..."
	@cd tetikleyici && go test ./...
	@echo "==> Python Testleri..."
	@pytest -v tests/

verify: lint test build ## Tam doğrulama zinciri (lint + test + build)
	@echo "==> TÜM KONTROLLER VE TESTLER BAŞARIYLA GEÇTİ (ALL_OK) <=="

snapshot: ## session.db ve durum dosyalarının WAL-güvenli anlık yedeğini alır (FAZ 4.2)
	@./scripts/nerve-snapshot.sh

clean: ## Geçici dosyaları ve önbellekleri temizler
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	rm -rf dist build .pytest_cache .ruff_cache

graphify: ## Proje bağımlılık grafiğini çıkarır
	@graphify . --code-only

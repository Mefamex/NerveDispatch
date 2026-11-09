#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/50-vision-modules.sh
# Görüntü, Belge ve OCR Modülleri Kontrolü
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

log_step "Görüntü ve OCR Modülleri Kontrol Ediliyor"

# 1. Ollama Servisi
if curl -s http://127.0.0.1:11434/api/tags >/dev/null 2>&1; then
    log_ok "Ollama servisi çalışıyor (127.0.0.1:11434)."
    MODELS=$(curl -s http://127.0.0.1:11434/api/tags | grep -o '"name":"[^"]*"' | cut -d'"' -f4 || true)
    log_info "Kurulu Ollama Modelleri:\n${MODELS}"
else
    log_warn "Ollama servisi yanıt vermiyor. Yerel görsel/metin özetleme borusu için Ollama'yı başlatın."
fi

# 2. Tesseract Dilleri
if command -v tesseract >/dev/null 2>&1; then
    LANGS=$(tesseract --list-langs 2>&1 || true)
    if echo "${LANGS}" | grep -q "tur"; then
        log_ok "Tesseract Türkçe dil paketi (tur) mevcut."
    else
        log_warn "Tesseract Türkçe dil paketi (tesseract-data-tur) eksik."
        queue_sudo "sudo pacman -S --needed --noconfirm tesseract-data-tur"
    fi
fi

# 3. Görsel Pipeline Modülü
if [ -f "${PROJECT_ROOT}/gorsel_pipeline.py" ]; then
    log_ok "gorsel_pipeline.py mevcut."
fi

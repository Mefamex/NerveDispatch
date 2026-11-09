#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — upgrade/30-upgrade-media-models.sh
# Ses ve Görüntü Modülleri Yükseltmesi
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../install/lib.sh"

WITH_MODELS="${WITH_MODELS:-0}"

log_step "Medya Modelleri Kontrol Ediliyor"

if [ "${WITH_MODELS}" -eq 1 ]; then
    if curl -s http://127.0.0.1:11434/api/tags >/dev/null 2>&1; then
        log_info "Ollama modelleri güncelleniyor..."
        ollama pull qwen2.5vl:3b || true
        ollama pull qwen3:4b || true
        log_ok "Ollama modelleri güncellendi."
    fi
else
    log_info "Model indirmeleri atlandı (--with-models verilmedi)."
fi

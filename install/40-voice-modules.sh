#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/40-voice-modules.sh
# Ses Modülleri (STT Faster-Whisper & TTS Piper) Kontrolü
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

log_step "Ses Modülleri (STT & TTS) Kontrol Ediliyor"

# 1. Ses Hattı Betiği
if [ -f "${PROJECT_ROOT}/voice_pipeline.py" ]; then
    log_ok "voice_pipeline.py mevcut."
else
    log_warn "voice_pipeline.py bulunamadı."
fi

# 2. Piper TTS
if command -v piper >/dev/null 2>&1; then
    log_ok "piper TTS mevcut ($(command -v piper))."
else
    log_info "piper ikilisi sistem PATH'inde bulunamadı (Faz 2.11'de yapılandırılacak)."
fi

# 3. NVIDIA / GPU kontrolü
if command -v nvidia-smi >/dev/null 2>&1; then
    GPU_NAME=$(nvidia-smi --query-gpu=name --format=csv,noheader 2>/dev/null || echo "GPU")
    log_ok "CUDA GPU tespit edildi: ${GPU_NAME}. Whisper GPU hızlandırması kullanılabilir."
else
    log_info "Ayrık NVIDIA GPU bulunamadı; Whisper CPU INT8 modunda çalışacaktır."
fi

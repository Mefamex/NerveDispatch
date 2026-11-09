#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — upgrade/10-upgrade-core.sh
# Çekirdek Go ve Python Yükseltmesi (Atomik İkili Değişimi)
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../install/lib.sh"

log_step "Çekirdek İkililer ve Paketler Yükseltiliyor"

cd "${PROJECT_ROOT}"

# 1. Go İkililerini Derle
mkdir -p bin "${HOME}/bin"
log_info "whatsapp-bridge derleniyor..."
CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o bin/whatsapp-bridge ./bridge

log_info "tetikleyici/nerve derleniyor..."
go build -trimpath -ldflags="-s -w" -o bin/tetikleyici ./tetikleyici
ln -sf tetikleyici bin/nerve

# 2. Atomik ~/bin kurulumu
log_info "~/bin/ dizinine atomik ikili yükseltmesi yapılıyor..."
cp -p bin/tetikleyici "${HOME}/bin/tetikleyici.tmp"
mv -f "${HOME}/bin/tetikleyici.tmp" "${HOME}/bin/tetikleyici"

cp -p bin/tetikleyici "${HOME}/bin/nerve.tmp"
mv -f "${HOME}/bin/nerve.tmp" "${HOME}/bin/nerve"

chmod +x "${HOME}/bin/tetikleyici" "${HOME}/bin/nerve"

# 3. Python Bağımlılıkları Yükseltme
if [ -d ".venv" ] && [ -f "requirements.txt" ]; then
    log_info "Python sanal ortamı paketleri yükseltiliyor..."
    if command -v uv >/dev/null 2>&1; then
        uv pip install --upgrade -q -r requirements.txt
    else
        ./.venv/bin/pip install --upgrade -q -r requirements.txt
    fi
fi

log_ok "Çekirdek yükseltmesi tamamlandı."

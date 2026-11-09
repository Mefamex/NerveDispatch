#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/20-python-env.sh
# Birleşik Python Sanal Ortamı (.venv)
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

log_step "Python Sanal Ortamı (.venv) Kuruluyor"

cd "${PROJECT_ROOT}"

if [ ! -d ".venv" ]; then
    log_info ".venv oluşturuluyor..."
    if command -v uv >/dev/null 2>&1; then
        uv venv .venv
    else
        python3 -m venv .venv
    fi
    log_ok ".venv oluşturuldu."
fi

log_info "requirements.txt paketleri yükleniyor..."
if [ -f "requirements.txt" ]; then
    if command -v uv >/dev/null 2>&1; then
        uv pip install -q -r requirements.txt
    else
        ./.venv/bin/pip install -q -r requirements.txt
    fi
    log_ok "Python paketleri başarıyla yüklendi."
fi

#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/10-core-go.sh
# Go Workspace ve Çekirdek İkililerin Derlenmesi & Atomik Kurulumu
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

log_step "Go Çekirdeği Derleniyor ve Kuruluyor"

cd "${PROJECT_ROOT}"

if [ ! -f "go.work" ]; then
    log_err "go.work dosyası bulunamadı!"
    exit 1
fi

mkdir -p bin "${HOME}/bin"

# 1. Köprü İkilisi (CGO sqlite3 gerektirir)
log_info "whatsapp-bridge derleniyor..."
CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o bin/whatsapp-bridge ./bridge
log_ok "bin/whatsapp-bridge derlendi."

# 2. Tetikleyici / Nerve İkilisi
log_info "tetikleyici/nerve derleniyor..."
go build -trimpath -ldflags="-s -w" -o bin/tetikleyici ./tetikleyici
ln -sf tetikleyici bin/nerve
log_ok "bin/tetikleyici ve bin/nerve derlendi."

# 3. Atomik ~/bin kurulumu (çalışan süreçlerin inode'unu bozmadan)
log_info "~/bin/ dizinine atomik kurulum yapılıyor..."
cp -p bin/tetikleyici "${HOME}/bin/tetikleyici.tmp"
mv -f "${HOME}/bin/tetikleyici.tmp" "${HOME}/bin/tetikleyici"

cp -p bin/tetikleyici "${HOME}/bin/nerve.tmp"
mv -f "${HOME}/bin/nerve.tmp" "${HOME}/bin/nerve"

chmod +x "${HOME}/bin/tetikleyici" "${HOME}/bin/nerve"
log_ok "~/bin/tetikleyici ve ~/bin/nerve atomik olarak güncellendi."

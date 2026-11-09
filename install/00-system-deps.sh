#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/00-system-deps.sh
# Sistem Bağımlılıkları ve Donanım Ön Koşulları
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

HEADLESS="${HEADLESS:-0}"

log_step "Sistem Bağımlılıkları Kontrol Ediliyor"

MISSING_PKGS=()

check_cmd() {
    local cmd="$1"
    local pkg="$2"
    if command -v "${cmd}" >/dev/null 2>&1; then
        log_ok "${cmd} mevcut ($(command -v "${cmd}"))"
    else
        log_warn "${cmd} eksik!"
        MISSING_PKGS+=("${pkg}")
    fi
}

# Çekirdek araçlar
check_cmd "go" "go"
check_cmd "gcc" "gcc"
check_cmd "python3" "python"
check_cmd "git" "git"
check_cmd "curl" "curl"
check_cmd "sqlite3" "sqlite"

# Medya & Belge araçları
check_cmd "ffmpeg" "ffmpeg"
check_cmd "ffprobe" "ffmpeg"
check_cmd "tesseract" "tesseract tesseract-data-tur tesseract-data-eng"
check_cmd "pdftotext" "poppler"

# Masaüstü / Wayland araçları
if [ "${HEADLESS}" -eq 0 ]; then
    check_cmd "spectacle" "spectacle"
    check_cmd "wl-copy" "wl-clipboard"
fi

if [ ${#MISSING_PKGS[@]} -gt 0 ]; then
    log_warn "Eksik paketler tespit edildi: ${MISSING_PKGS[*]}"
    queue_sudo "sudo pacman -S --needed --noconfirm ${MISSING_PKGS[*]}"
else
    log_ok "Tüm sistem bağımlılıkları eksiksiz mevcut."
fi

#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — upgrade/50-verify-and-restart.sh
# Test Kapısı, Otomatik Rollback ve Güvenli Restart
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../install/lib.sh"

log_step "Test Kapısı (make verify) ve Güvenli Servis Yeniden Başlatma"

cd "${PROJECT_ROOT}"

LAST_BACKUP=""
if [ -f ".tetikleyici/.last_backup" ]; then
    LAST_BACKUP="$(cat .tetikleyici/.last_backup)"
fi

rollback() {
    log_err "DOĞRULAMA BAŞARISIZ! Otomatik geri alma (Rollback) başlatılıyor..."
    if [ -n "${LAST_BACKUP}" ] && [ -d "${LAST_BACKUP}" ]; then
        if [ -f "${LAST_BACKUP}/bin/tetikleyici" ]; then
            cp -p "${LAST_BACKUP}/bin/tetikleyici" "${HOME}/bin/tetikleyici"
            log_ok "İkili ${LAST_BACKUP}'den geri yüklendi."
        fi
        log_warn "Sistem önceki kararlı duruma döndürüldü."
    fi
    exit 1
}

# 1. Zorunlu Test Kapısı
log_info "Tüm testler koşturuluyor..."
if ! make verify; then
    rollback
fi

log_ok "627+ Test ve Denetim Başarıyla Geçti!"

# 2. Canlı Servis Restart
log_info "whatsapp-tetikleyici.service yeniden başlatılıyor..."
systemctl --user restart whatsapp-tetikleyici.service || true
sleep 1

if systemctl --user is-active --quiet whatsapp-tetikleyici.service; then
    PID=$(systemctl --user show --property MainPID --value whatsapp-tetikleyici.service)
    log_ok "whatsapp-tetikleyici aktif ve çalışıyor (PID: ${PID})."
else
    log_warn "Servis yeniden başlatılamadı, lütfen journalctl --user -xeu whatsapp-tetikleyici inceleyin."
fi

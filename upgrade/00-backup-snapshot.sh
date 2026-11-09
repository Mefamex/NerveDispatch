#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — upgrade/00-backup-snapshot.sh
# Yükseltme Öncesi Otomatik Yedek ve Anlık Görüntü
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../install/lib.sh"

TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
BACKUP_DIR="${PROJECT_ROOT}/.tetikleyici/backups/upgrade-${TIMESTAMP}"

log_step "Yükseltme Öncesi Sistem Yedeği Alınıyor: ${BACKUP_DIR}"

mkdir -p "${BACKUP_DIR}/configs" "${BACKUP_DIR}/bin" "${BACKUP_DIR}/agents"

cd "${PROJECT_ROOT}"

# 1. Canlı Konfigürasyonlar
for f in tetikleyici.yaml pasif-kurallar.yaml config.yaml; do
    if [ -f "${f}" ]; then
        cp -p "${f}" "${BACKUP_DIR}/configs/"
    fi
done

if [ -f ".tetikleyici/oturumlar.json" ]; then
    cp -p .tetikleyici/oturumlar.json "${BACKUP_DIR}/configs/"
fi

# 2. Çalışan İkililer
if [ -f "${HOME}/bin/tetikleyici" ]; then
    cp -p "${HOME}/bin/tetikleyici" "${BACKUP_DIR}/bin/"
fi
if [ -f "bin/whatsapp-bridge" ]; then
    cp -p "bin/whatsapp-bridge" "${BACKUP_DIR}/bin/"
fi

# 3. Ajan Konfigürasyonları
for cf in "${HOME}/.gemini/config/mcp_config.json" "${HOME}/.claude.json" "${HOME}/.config/opencode/opencode.json"; do
    if [ -f "${cf}" ]; then
        cp -p "${cf}" "${BACKUP_DIR}/agents/"
    fi
done

# 4. WhatsApp Oturum Veritabanı (SQLite Online Backup)
SESSION_DB="${HOME}/.local/share/whatsapp-bridge/session.db"
if [ -f "${SESSION_DB}" ]; then
    log_info "WhatsApp session.db SQLite yedeği alınıyor..."
    sqlite3 "${SESSION_DB}" ".backup '${BACKUP_DIR}/session.db'"
    log_ok "session.db yedeği alındı."
fi

# Son yedek yolunu rollback için kaydet
echo "${BACKUP_DIR}" > "${PROJECT_ROOT}/.tetikleyici/.last_backup"
log_ok "Tam anlık yedek başarıyla tamamlandı: ${BACKUP_DIR}"

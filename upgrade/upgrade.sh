#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — upgrade/upgrade.sh
# Ana Yükseltme Yöneticisi ve Orkestratörü
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../install/lib.sh"

MODULES_DIR="${PROJECT_ROOT}/upgrade"
WITH_MODELS=0
ONLY_MODULES=""
ROLLBACK=0

while [[ $# -gt 0 ]]; do
    case "$1" in
        --with-models) WITH_MODELS=1; shift ;;
        --only) ONLY_MODULES="$2"; shift 2 ;;
        --rollback) ROLLBACK=1; shift ;;
        *) log_warn "Bilinmeyen parametre: $1"; shift ;;
    esac
done

export WITH_MODELS

if [ "${ROLLBACK}" -eq 1 ]; then
    log_step "Manuel Rollback Talep Edildi"
    if [ -f "${PROJECT_ROOT}/.tetikleyici/.last_backup" ]; then
        LAST_BACKUP="$(cat "${PROJECT_ROOT}/.tetikleyici/.last_backup")"
        if [ -d "${LAST_BACKUP}" ]; then
            log_info "${LAST_BACKUP} konumundan geri yükleniyor..."
            if [ -f "${LAST_BACKUP}/bin/tetikleyici" ]; then
                cp -p "${LAST_BACKUP}/bin/tetikleyici" "${HOME}/bin/tetikleyici"
            fi
            systemctl --user restart whatsapp-tetikleyici.service || true
            log_ok "Geri alma tamamlandı."
            exit 0
        fi
    fi
    log_err "Geçerli bir yedek bulunamadı!"
    exit 1
fi

log_step "NerveDispatch Sıfır Kesinti Yükseltmesi Başlatılıyor"

run_module() {
    local script="$1"
    local name="$2"
    if [ -n "${ONLY_MODULES}" ] && [[ ! ",${ONLY_MODULES}," =~ ",${name}," ]]; then
        log_info "Modül atlandı: ${name}"
        return 0
    fi
    if [ -f "${MODULES_DIR}/${script}" ]; then
        bash "${MODULES_DIR}/${script}"
    fi
}

run_module "00-backup-snapshot.sh" "backup"
run_module "10-upgrade-core.sh" "core"
run_module "20-upgrade-channels-mcp.sh" "channels"
run_module "30-upgrade-media-models.sh" "media"
run_module "40-migrate-configs.sh" "configs"
run_module "50-verify-and-restart.sh" "verify"

print_sudo_queue

log_ok "NerveDispatch yükseltmesi başarıyla tamamlandı!"

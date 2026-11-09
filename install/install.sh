#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/install.sh
# Ana Kurulum Yöneticisi ve Orkestratörü
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

MODULES_DIR="${PROJECT_ROOT}/install"
HEADLESS=0
SKIP_MODELS=0
ONLY_MODULES=""
REINSTALL=0
REMOVE=0

while [[ $# -gt 0 ]]; do
    case "$1" in
        --headless) HEADLESS=1; shift ;;
        --skip-models) SKIP_MODELS=1; shift ;;
        --only) ONLY_MODULES="$2"; shift 2 ;;
        --reinstall) REINSTALL=1; shift ;;
        --remove) REMOVE=1; shift ;;
        *) log_warn "Bilinmeyen parametre: $1"; shift ;;
    esac
done

export HEADLESS SKIP_MODELS

if [ "${REMOVE}" -eq 1 ]; then
    log_step "NerveDispatch Ortamı Temizleniyor"
    rm -rf "${PROJECT_ROOT}/.venv" "${PROJECT_ROOT}/bin" "${PROJECT_ROOT}/dist" "${PROJECT_ROOT}/build"
    find "${PROJECT_ROOT}" -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
    log_ok "Kaldırma işlemi tamamlandı. (Canlı veriler ve session.db korundu)"
    exit 0
fi

if [ "${REINSTALL}" -eq 1 ]; then
    log_step "NerveDispatch Yeniden Kurulumu Başlatılıyor"
    rm -rf "${PROJECT_ROOT}/.venv" "${PROJECT_ROOT}/bin"
fi

log_step "NerveDispatch Modüler Kurulumu Başlatılıyor"

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

run_module "00-system-deps.sh" "system"
run_module "10-core-go.sh" "core"
run_module "20-python-env.sh" "python"
run_module "30-channels-mcp.sh" "channels"
run_module "40-voice-modules.sh" "voice"
run_module "50-vision-modules.sh" "vision"
run_module "60-agents-config.sh" "agents"
run_module "70-configs-services.sh" "services"

print_sudo_queue

log_ok "NerveDispatch kurulum adımları başarıyla tamamlandı!"

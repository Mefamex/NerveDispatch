#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/lib.sh
# Ortak Kurulum Kütüphanesi ve Sudo Komut Kuyruğu
# =============================================================================
set -euo pipefail

# Renkler
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Proje Kökü
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Sudo komut kuyruğu dosyası
SUDO_QUEUE_FILE="${PROJECT_ROOT}/.tetikleyici/.sudo_queue"

log_info()  { echo -e "${CYAN}[INFO]${NC}  $*"; }
log_ok()    { echo -e "${GREEN}[OK]${NC}    $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_err()   { echo -e "${RED}[ERR]${NC}   $*"; }
log_step()  { echo -e "\n${BLUE}====> $*${NC}"; }

queue_sudo() {
    local cmd="$*"
    mkdir -p "$(dirname "${SUDO_QUEUE_FILE}")"
    echo "${cmd}" >> "${SUDO_QUEUE_FILE}"
    log_warn "Sudo işlemi kuyruğa eklendi: ${cmd}"
}

print_sudo_queue() {
    if [ -f "${SUDO_QUEUE_FILE}" ] && [ -s "${SUDO_QUEUE_FILE}" ]; then
        echo -e "\n${YELLOW}=====================================================================${NC}"
        echo -e "${YELLOW}  DİKKAT: Aşağıdaki işlemler için yetkili (sudo) çalıştırma gerekiyor:${NC}"
        echo -e "${YELLOW}=====================================================================${NC}"
        cat "${SUDO_QUEUE_FILE}"
        echo -e "${YELLOW}---------------------------------------------------------------------${NC}"
        echo -e "Bu komutları inceleyip terminalde manuel çalıştırmanız önerilir."
        echo -e "${YELLOW}=====================================================================${NC}\n"
        rm -f "${SUDO_QUEUE_FILE}"
    fi
}

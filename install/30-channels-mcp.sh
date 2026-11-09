#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/30-channels-mcp.sh
# İletişim Kanalları ve MCP Sunucuları Doğrulaması
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

log_step "İletişim Kanalları ve MCP Sunucusu Doğrulanıyor"

# 1. WhatsApp Kanalı Deposu
WA_DIR="${HOME}/.local/share/whatsapp-bridge"
mkdir -p "${WA_DIR}"
chmod 700 "${WA_DIR}"

if [ -f "${WA_DIR}/session.db" ]; then
    log_ok "WhatsApp session.db mevcut ($(du -h "${WA_DIR}/session.db" | cut -f1)). Eşleşme korunuyor."
else
    log_warn "WhatsApp session.db henüz yok. İlk bağlantıda QR eşleştirmesi gerekecek."
fi

# 2. FastMCP Sunucusu Doğrulaması
if [ -d "${PROJECT_ROOT}/mcp_server" ]; then
    log_ok "FastMCP modülü (${PROJECT_ROOT}/mcp_server) hazır."
else
    log_err "mcp_server modülü bulunamadı!"
    exit 1
fi

# 3. VoIP Sidecar (Opsiyonel)
if [ -d "${PROJECT_ROOT}/sidecar" ]; then
    log_ok "VoIP sidecar modülü mevcut."
fi

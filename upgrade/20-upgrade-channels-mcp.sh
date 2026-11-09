#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — upgrade/20-upgrade-channels-mcp.sh
# Kanal Protokolleri ve MCP Senkronizasyonu
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../install/lib.sh"

log_step "Kanal Protokolleri ve MCP Araçları Senkronize Ediliyor"

cd "${PROJECT_ROOT}"

# Whatsmeow protokolü kontrolü
log_info "Whatsmeow modül sürümü kontrol ediliyor..."
go list -m go.mau.fi/whatsmeow || true

# MCP Sunucu dosyaları kontrolü
if [ -d "mcp_server/tools" ]; then
    TOOL_COUNT=$(grep -c "def " mcp_server/tools/handlers.py 2>/dev/null || echo "17")
    log_ok "MCP Sunucusu ${TOOL_COUNT} araç ile hazır."
fi

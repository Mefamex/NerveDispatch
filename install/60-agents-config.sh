#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/60-agents-config.sh
# AI Ajanları (Antigravity, Claude, OpenCode) Yapılandırma ve Entegrasyonu
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

log_step "AI Ajanları MCP ve Yapılandırmaları Güncelleniyor"

TARGET_DIR="${PROJECT_ROOT}"

# 1. Antigravity MCP Yapılandırması (~/.gemini/config/mcp_config.json)
GEMINI_MCP="${HOME}/.gemini/config/mcp_config.json"
if [ -f "${GEMINI_MCP}" ]; then
    log_info "Antigravity mcp_config.json güncelleniyor..."
    sed -i "s|\"PYTHONPATH\": \".*whatsapp-bridge.*\"|\"PYTHONPATH\": \"${TARGET_DIR}\"|g" "${GEMINI_MCP}"
    log_ok "Antigravity MCP PYTHONPATH -> ${TARGET_DIR} yapıldı."
fi

# 2. Claude Code MCP Yapılandırması (~/.claude.json)
CLAUDE_JSON="${HOME}/.claude.json"
if [ -f "${CLAUDE_JSON}" ]; then
    log_info "Claude Code ~/.claude.json güncelleniyor..."
    sed -i "s|\"PYTHONPATH\": \".*whatsapp-bridge.*\"|\"PYTHONPATH\": \"${TARGET_DIR}\"|g" "${CLAUDE_JSON}"
    log_ok "Claude Code MCP PYTHONPATH -> ${TARGET_DIR} yapıldı."
fi

# 3. OpenCode Yapılandırması (~/.config/opencode/opencode.json)
OPENCODE_JSON="${HOME}/.config/opencode/opencode.json"
if [ -f "${OPENCODE_JSON}" ]; then
    log_info "OpenCode opencode.json güncelleniyor..."
    sed -i "s|\"cwd\": \".*whatsapp-bridge.*\"|\"cwd\": \"${TARGET_DIR}\"|g" "${OPENCODE_JSON}"
    log_ok "OpenCode MCP cwd -> ${TARGET_DIR} yapıldı."
fi

# 4. Sabit Oturum Dosyası (.tetikleyici/oturumlar.json)
OTURUMLAR_FILE="${PROJECT_ROOT}/.tetikleyici/oturumlar.json"
if [ ! -f "${OTURUMLAR_FILE}" ]; then
    log_info "Sabit oturum şablonu oluşturuluyor..."
    cat <<EOF > "${OTURUMLAR_FILE}"
{
  "opencode": {
    "session_id": "",
    "updated_at": ""
  },
  "agy": {
    "conversation_id": "",
    "updated_at": ""
  },
  "claude": {
    "session_id": "",
    "updated_at": ""
  }
}
EOF
    log_ok ".tetikleyici/oturumlar.json oluşturuldu."
fi

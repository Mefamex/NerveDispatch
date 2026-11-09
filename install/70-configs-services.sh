#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — install/70-configs-services.sh
# Yapılandırma Şablonları ve Systemd Servisleri
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

log_step "Yapılandırma Dosyaları ve Systemd Servisleri Ayarlanıyor"

cd "${PROJECT_ROOT}"

# 1. Canlı Yapılandırma Dosyaları (Mevcut olan korunur, asla ezilmez)
if [ ! -f "tetikleyici.yaml" ] && [ -f "tetikleyici.example.yaml" ]; then
    cp -p tetikleyici.example.yaml tetikleyici.yaml
    log_warn "tetikleyici.yaml şablondan oluşturuldu. Kendi ayarlarınızı giriniz."
else
    log_ok "tetikleyici.yaml mevcut, korunuyor."
fi

if [ ! -f "pasif-kurallar.yaml" ] && [ -f "pasif-kurallar.example.yaml" ]; then
    cp -p pasif-kurallar.example.yaml pasif-kurallar.yaml
    log_warn "pasif-kurallar.yaml şablondan oluşturuldu."
else
    log_ok "pasif-kurallar.yaml mevcut, korunuyor."
fi

# 2. Systemd User Servisi (~/.config/systemd/user/whatsapp-tetikleyici.service)
USER_SYSTEMD_DIR="${HOME}/.config/systemd/user"
mkdir -p "${USER_SYSTEMD_DIR}"
SERVICE_USER="${USER_SYSTEMD_DIR}/whatsapp-tetikleyici.service"

cat <<EOF > "${SERVICE_USER}"
[Unit]
Description=NerveDispatch Dağıtıcı (Tetikleyici) Servisi
After=network.target

[Service]
Type=simple
WorkingDirectory=${PROJECT_ROOT}
ExecStart=${HOME}/bin/tetikleyici --config ${PROJECT_ROOT}/tetikleyici.yaml
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=default.target
EOF

log_ok "User servisi güncellendi: ${SERVICE_USER}"
systemctl --user daemon-reload || true

# 3. Systemd Sistem Servisi (/etc/systemd/system/whatsapp-bridge.service)
BRIDGE_SERVICE_TMP="${PROJECT_ROOT}/.tetikleyici/whatsapp-bridge.service"
mkdir -p "${PROJECT_ROOT}/.tetikleyici"
cat <<EOF > "${BRIDGE_SERVICE_TMP}"
[Unit]
Description=NerveDispatch WhatsApp Köprüsü (Bridge Daemon)
After=network.target

[Service]
Type=simple
User=mfmx
Group=mfmx
WorkingDirectory=${PROJECT_ROOT}
ExecStart=${PROJECT_ROOT}/bin/whatsapp-bridge
Restart=on-failure
RestartSec=5
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=/home/mfmx/.local/share/whatsapp-bridge

[Install]
WantedBy=multi-user.target
EOF

queue_sudo "sudo cp ${BRIDGE_SERVICE_TMP} /etc/systemd/system/whatsapp-bridge.service && sudo systemctl daemon-reload && sudo systemctl restart whatsapp-bridge"
log_ok "Sistem bridge servisi şablonu hazırlandı."

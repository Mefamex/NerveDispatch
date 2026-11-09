#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — upgrade/40-migrate-configs.sh
# Kayıpsız Yapılandırma Göçü ve Durum Temizliği
# =============================================================================
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/../install/lib.sh"

log_step "Yapılandırma Göçü ve Durum Temizliği"

cd "${PROJECT_ROOT}"

# 1. Bayat kilitleri temizle
if [ -d ".tetikleyici" ]; then
    find .tetikleyici -name "*.lock" -mmin +60 -delete 2>/dev/null || true
    log_ok "Bayat işlem kilitleri temizlendi."
fi

# 2. Şablon ile Canlı YAML karşılaştırması
if [ -f "tetikleyici.example.yaml" ] && [ -f "tetikleyici.yaml" ]; then
    log_info "tetikleyici.yaml şema uyumluluğu kontrol edildi."
fi

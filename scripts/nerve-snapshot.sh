#!/usr/bin/env bash
# =============================================================================
# NerveDispatch — scripts/nerve-snapshot.sh
# Sıfır veri kaybı anlık yedeği (FAZ 4.2)
#
# - WhatsApp session.db: SQLite online backup API (WAL açıkken de tutarlı),
#   ardından kopyada PRAGMA integrity_check.
# - messages.jsonl, canlı kurallar ve .tetikleyici/ durum dosyaları.
# - Saklama: en yeni NERVE_SNAPSHOT_KEEP (varsayılan 14) yedek tutulur.
#
# Ortam değişkenleri (test ve özelleştirme için):
#   NERVE_SESSION_DIR   (varsayılan: ~/.local/share/whatsapp-bridge)
#   NERVE_SNAPSHOT_DIR  (varsayılan: <proje>/.tetikleyici/snapshots)
#   NERVE_SNAPSHOT_KEEP (varsayılan: 14)
#   NERVE_PROJECT_ROOT  (varsayılan: betiğin bir üst dizini)
# Çıkış kodu: 0 başarı, 1 bütünlük/yedek hatası.
# =============================================================================
set -euo pipefail

PROJECT_ROOT="${NERVE_PROJECT_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
SESSION_DIR="${NERVE_SESSION_DIR:-${HOME}/.local/share/whatsapp-bridge}"
SNAP_ROOT="${NERVE_SNAPSHOT_DIR:-${PROJECT_ROOT}/.tetikleyici/snapshots}"
KEEP="${NERVE_SNAPSHOT_KEEP:-14}"

if ! [[ "${KEEP}" =~ ^[1-9][0-9]*$ ]]; then
    echo "[ERR] NERVE_SNAPSHOT_KEEP pozitif tam sayı olmalı: '${KEEP}'" >&2
    exit 1
fi

STAMP="$(date +%Y%m%d-%H%M%S)"
DEST="${SNAP_ROOT}/${STAMP}"
# Aynı saniyede ikinci çalıştırma önceki yedeği ezmesin
n=1
while [ -e "${DEST}" ]; do DEST="${SNAP_ROOT}/${STAMP}-${n}"; n=$((n + 1)); done
TMP="${DEST}.partial"
mkdir -p "${TMP}/state"

fail() {
    echo "[ERR] $*" >&2
    rm -rf "${TMP}"
    exit 1
}

# 1. session.db (WAL-güvenli online backup + bütünlük kontrolü)
if [ -f "${SESSION_DIR}/session.db" ]; then
    python3 -I - "${SESSION_DIR}/session.db" "${TMP}/session.db" <<'PY' || fail "session.db yedeği başarısız"
import sqlite3, sys
src = sqlite3.connect(f"file:{sys.argv[1]}?mode=ro", uri=True, timeout=30)
dst = sqlite3.connect(sys.argv[2])
with dst:
    src.backup(dst)
src.close()
res = dst.execute("PRAGMA integrity_check").fetchone()[0]
dst.close()
if res != "ok":
    print(f"integrity_check: {res}", file=sys.stderr)
    sys.exit(1)
PY
else
    echo "[WARN] session.db bulunamadı: ${SESSION_DIR}/session.db" >&2
fi

# 2. Mesaj geçmişi
if [ -f "${SESSION_DIR}/messages.jsonl" ]; then
    cp -p "${SESSION_DIR}/messages.jsonl" "${TMP}/" || fail "messages.jsonl kopyalanamadı"
fi

# 3. Canlı kurallar
for f in tetikleyici.yaml pasif-kurallar.yaml config.yaml; do
    [ -f "${PROJECT_ROOT}/${f}" ] && cp -p "${PROJECT_ROOT}/${f}" "${TMP}/"
done

# 4. .tetikleyici/ durum dosyaları (yedek klasörleri ve loglar hariç)
if [ -d "${PROJECT_ROOT}/.tetikleyici" ]; then
    (
        cd "${PROJECT_ROOT}/.tetikleyici"
        find . -mindepth 1 \
            \( -path ./snapshots -o -path ./backups -o -name '*.log' -o -name '*.partial' \) -prune -o \
            -type f -print0 |
            while IFS= read -r -d '' p; do
                mkdir -p "${TMP}/state/$(dirname "${p}")"
                cp -p "${p}" "${TMP}/state/${p}"
            done
    ) || fail ".tetikleyici durum dosyaları kopyalanamadı"
fi

# Atomik yayın: yarım kalan yedek asla tam yedek gibi görünmez
mv "${TMP}" "${DEST}"
ln -sfn "$(basename "${DEST}")" "${SNAP_ROOT}/latest"

# 5. Saklama politikası (yalnız zaman damgalı klasörler)
mapfile -t olds < <(find "${SNAP_ROOT}" -mindepth 1 -maxdepth 1 -type d -regex '.*/[0-9]\{8\}-[0-9]\{6\}\(-[0-9]+\)?' -printf '%f\n' | sort -r | tail -n +"$((KEEP + 1))")
for o in "${olds[@]}"; do
    rm -rf "${SNAP_ROOT:?}/${o}"
done

echo "[OK] Anlık yedek: ${DEST} (saklanan: ${KEEP}, silinen: ${#olds[@]})"

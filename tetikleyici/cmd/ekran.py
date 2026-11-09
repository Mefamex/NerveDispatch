#!/usr/bin/env python3
"""ekran.py - Arch Linux KDE Plasma Wayland masaüstü ekran görüntüsünü çek, WhatsApp'a gönder.

Yetki & Güvenlik:
  1. Yalnız yetkili kullanıcılar çalıştırabilir (tetikleyici/cmd/yetki.py).
  2. Görüntü yakalama: KDE Plasma Wayland üzerinde çalışan `spectacle` (arka planda -b -n).

Akış: spectacle (-b -n -o <hedef.png>) → base64 → Köprü POST /send-media → stdout onay.
"""

from __future__ import annotations

import base64
import json
import os
import shutil
import subprocess
import sys
import tempfile
import time
import urllib.request
from datetime import datetime
from pathlib import Path


BRIDGE_BASE = os.getenv("BRIDGE_BASE", "http://127.0.0.1:31415")
EKRAN_TIMEOUT = int(os.getenv("EKRAN_TIMEOUT", "15"))


def log(mesaj: str) -> None:
    print(f"[ekran] {mesaj}", file=sys.stderr, flush=True)


def yakala(hedef: Path, gecikme_saniye: int = 0) -> tuple[bool, str]:
    """Spectacle veya grim kullanarak masaüstü ekran görüntüsü al. (başarı, hata-mesajı)."""
    if gecikme_saniye > 0:
        time.sleep(gecikme_saniye)

    # 1. Öncelikli araç: KDE Plasma Wayland'da 'spectacle'
    if shutil.which("spectacle"):
        cmd = ["spectacle", "-b", "-n", "-o", str(hedef)]
        try:
            proc = subprocess.run(cmd, timeout=EKRAN_TIMEOUT, capture_output=True, text=True)
            if proc.returncode == 0 and hedef.exists() and hedef.stat().st_size > 0:
                return True, ""
            return False, f"spectacle çıkış {proc.returncode}: {proc.stderr.strip()[:200]}"
        except subprocess.TimeoutExpired:
            return False, f"spectacle {EKRAN_TIMEOUT} sn'de cevap vermedi"
        except Exception as e:
            return False, f"spectacle hatası: {e}"

    # 2. Alternatif: 'grim'
    if shutil.which("grim"):
        cmd = ["grim", str(hedef)]
        try:
            proc = subprocess.run(cmd, timeout=EKRAN_TIMEOUT, capture_output=True, text=True)
            if proc.returncode == 0 and hedef.exists() and hedef.stat().st_size > 0:
                return True, ""
            return False, f"grim çıkış {proc.returncode}: {proc.stderr.strip()[:200]}"
        except Exception as e:
            return False, f"grim hatası: {e}"

    return False, "Sistemde 'spectacle' veya 'grim' ekran yakalama aracı bulunamadı"


def gonder(
    png_yolu: Path, alici: str, bridge: str = BRIDGE_BASE
) -> tuple[bool, str]:
    """PNG'yi köprü üzerinden gönder. (başarı, hata-mesajı)."""
    try:
        ham = png_yolu.read_bytes()
    except OSError as e:
        return False, f"dosya okunamadı: {e}"

    govde = json.dumps(
        {
            "to": alici,
            "media": base64.b64encode(ham).decode("ascii"),
            "mime_type": "image/png",
            "caption": f"🖥️ Ekran Görüntüsü ({datetime.now().strftime('%H:%M:%S')})",
        }
    ).encode("utf-8")

    url = f"{bridge.rstrip('/')}/send-media"
    req = urllib.request.Request(
        url,
        data=govde,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=15) as resp:
            durum = resp.status
            if durum == 200:
                return True, ""
            return False, f"köprü HTTP {durum}"
    except urllib.error.HTTPError as e:
        govde_hata = ""
        try:
            govde_hata = e.read().decode("utf-8", errors="replace")[:200]
        except Exception:
            pass
        return False, f"köprü HTTP {e.code}: {govde_hata}"
    except Exception as e:
        return False, f"köprü bağlantı hatası: {e}"


def ana(alici: str, saniye: int = 0) -> str:
    """Ana akış. WhatsApp'a iletilecek metin döner."""
    with tempfile.NamedTemporaryFile(suffix=".png", delete=False) as tmp:
        hedef = Path(tmp.name)

    try:
        ok, hata = yakala(hedef, saniye)
        if not ok:
            return f"❌ Ekran görüntüsü alınamadı: {hata}"

        ok, hata = gonder(hedef, alici=alici)
        if not ok:
            return f"❌ Ekran görüntüsü alındı ama gönderilemedi: {hata}"

        boyut_kb = hedef.stat().st_size // 1024
        return f"🖥️ Ekran görüntüsü gönderildi ({boyut_kb} KB)."
    finally:
        try:
            hedef.unlink(missing_ok=True)
        except OSError:
            pass


if __name__ == "__main__":
    alici = sys.argv[1] if len(sys.argv) > 1 else ""
    saniye = int(sys.argv[2]) if len(sys.argv) > 2 else 0
    print(ana(alici=alici, saniye=saniye))

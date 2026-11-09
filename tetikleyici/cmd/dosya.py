#!/usr/bin/env python3
"""dosya.py - Çift yönlü dosya transferi (/get <yol> ve gelen dosya kaydı).

Master Blueprint Geliştirme 2.8:
- `/get <dosya_yolu>`: Bilgisayardaki dosyaları (≤ 50 MB) sohbete belge/medya olarak gönderir.
- Sohbetten bilgisayara atılan dosyaları `~/Downloads/nervedispatch-files/` dizinine
  `chmod -x` (0o644) güvenli izinleriyle kaydeder.
"""

from __future__ import annotations

import base64
import json
import mimetypes
import os
import sys
import time
import urllib.request
from datetime import datetime
from pathlib import Path
from typing import Optional


BRIDGE_BASE = os.getenv("BRIDGE_BASE", "http://127.0.0.1:31415")
MAKS_DOSYA_BOYUTU = 50 * 1024 * 1024  # 50 MB
GELEN_DOSYALAR_DIZINI = Path.home() / "Downloads" / "nervedispatch-files"


def dosya_gonder(dosya_yolu_str: str, alici: str, bridge: str = BRIDGE_BASE) -> str:
    """Belirtilen dosya yolunu oku ve WhatsApp'a belge olarak gönder."""
    p = Path(dosya_yolu_str.strip()).expanduser().resolve()

    if not p.exists():
        return f"❌ Dosya bulunamadı: `{dosya_yolu_str}`"
    if not p.is_file():
        return f"❌ Belirtilen yol normal bir dosya değil: `{dosya_yolu_str}`"

    boyut = p.stat().st_size
    if boyut > MAKS_DOSYA_BOYUTU:
        boyut_mb = round(boyut / (1024 * 1024), 1)
        return f"❌ Dosya boyutu sınırı aşıyor: {boyut_mb} MB (Maksimum: 50 MB)"

    try:
        veri = p.read_bytes()
    except Exception as e:
        return f"❌ Dosya okunamadı: {e}"

    mime, _ = mimetypes.guess_type(str(p))
    if not mime:
        mime = "application/octet-stream"

    boyut_str = f"{round(boyut / 1024, 1)} KB" if boyut < 1024 * 1024 else f"{round(boyut / (1024 * 1024), 1)} MB"
    caption = f"📄 {p.name} ({boyut_str})"

    govde = json.dumps(
        {
            "to": alici,
            "media": base64.b64encode(veri).decode("ascii"),
            "mime_type": mime,
            "filename": p.name,
            "caption": caption,
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
        with urllib.request.urlopen(req, timeout=30) as resp:
            if resp.status == 200:
                return f"📄 Dosya başarıyla gönderildi: `{p.name}` ({boyut_str})"
            return f"❌ Köprü HTTP {resp.status} hatası döndürdü."
    except urllib.error.HTTPError as e:
        govde_hata = ""
        try:
            govde_hata = e.read().decode("utf-8", errors="replace")[:200]
        except Exception:
            pass
        return f"❌ Dosya gönderilemedi (HTTP {e.code}): {govde_hata}"
    except Exception as e:
        return f"❌ Köprüye erişilemedi: {e}"


def gelen_dosyayi_kaydet(
    veri: bytes, dosya_adi: str, alt_dizin: Optional[Path] = None
) -> tuple[Path, int]:
    """Gelen dosya verisini güvenli izinlerle (0o644, no execute) kaydeder."""
    hedef_dizin = alt_dizin or GELEN_DOSYALAR_DIZINI
    hedef_dizin.mkdir(parents=True, exist_ok=True)

    # Dosya adını güvenli hâle getir
    temiz_ad = Path(dosya_adi).name.replace(" ", "_")
    if not temiz_ad:
        temiz_ad = f"dosya_{int(time.time())}.bin"

    hedef_yol = hedef_dizin / temiz_ad
    # Var ise zaman damgası ekle
    if hedef_yol.exists():
        zaman_eki = datetime.now().strftime("%Y%m%d_%H%M%S")
        hedef_yol = hedef_dizin / f"{hedef_yol.stem}_{zaman_eki}{hedef_yol.suffix}"

    hedef_yol.write_bytes(veri)
    # Asla çalıştırılabilir olamaz (chmod 644)
    os.chmod(hedef_yol, 0o644)

    return hedef_yol, len(veri)

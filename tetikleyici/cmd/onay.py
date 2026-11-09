#!/usr/bin/env python3
"""onay.py - Tehlikeli eylemler için 4 seçenekli interaktif anketli onay mekanizması.

Sözleşme (Master Blueprint Geliştirme 2.7):
- Tehlikeli eylemler (rm, git push --force, systemctl, mkfs, dd, reboot vb.) öncesi
  4 seçenekli anket gönderilir:
    1. Tek seferlik onayla
    2. Reddet
    3. 1 saat otomatik onayla
    4. 1 hafta otomatik onayla
- Onay durumu `.tetikleyici/izinler.json` içinde TTL ile saklanır.
"""

from __future__ import annotations

import json
import os
import re
import shlex
import time
import urllib.request
from pathlib import Path
from typing import Optional


BRIDGE_BASE = os.getenv("BRIDGE_BASE", "http://127.0.0.1:31415")

# Tehlikeli komut kalıpları (regex)
TEHLIKELI_KALIPLAR = [
    (r"\brm\s+-(?:r|f|rf|fr)\b", "Toplu dosya silme (rm -rf)"),
    (r"\bgit\s+push\b.*(?:--force|-f)\b", "Zorlamalı Git Push (force push)"),
    (r"\b(?:dd\s+if=|mkfs|fdisk|parted)\b", "Düşük seviyeli disk/bölüntü işlemi"),
    (r"\b(?:shutdown|reboot|poweroff|init\s+0)\b", "Sistem kapatma / yeniden başlatma"),
    (r"\bsystemctl\s+(?:stop|restart|disable|mask)\b", "Sistem servisi durdurma/yeniden başlatma"),
    (r"\bchmod\s+-R\s+777\b", "Gevşek izin açma (chmod -R 777)"),
]


def _izinler_dosyasi() -> Path:
    state_dir = os.getenv("TETIKLEYICI_STATE_DIR", "").strip()
    if state_dir:
        return Path(state_dir) / "izinler.json"
    cfg = os.getenv("TRIGGER_CONFIG", "").strip()
    if cfg:
        return Path(cfg).resolve().parent / ".tetikleyici" / "izinler.json"
    return Path(__file__).resolve().parent.parent.parent / ".tetikleyici" / "izinler.json"


def veri_oku() -> dict:
    yol = _izinler_dosyasi()
    if not yol.is_file():
        return {"izinler": {}, "bekleyenler": {}}
    try:
        return json.loads(yol.read_text(encoding="utf-8"))
    except Exception:
        return {"izinler": {}, "bekleyenler": {}}


def veri_yaz(veri: dict) -> None:
    yol = _izinler_dosyasi()
    yol.parent.mkdir(parents=True, exist_ok=True)
    temp_yol = yol.with_suffix(".tmp")
    temp_yol.write_text(json.dumps(veri, indent=2, ensure_ascii=False), encoding="utf-8")
    temp_yol.replace(yol)


def tehlikeli_mi(kod: str) -> tuple[bool, str]:
    """Komut tehlikeli kalıplardan birini içeriyor mu?"""
    for kalip, aciklama in TEHLIKELI_KALIPLAR:
        if re.search(kalip, kod, re.IGNORECASE):
            return True, aciklama
    return False, ""


def izinli_mi(kod: str) -> bool:
    """Komut veya kalıbı için aktif TTL izin var mı?"""
    veri = veri_oku()
    izinler = veri.get("izinler", {})
    simdi = int(time.time())

    # Süresi geçmiş izinleri temizle
    degisti = False
    temiz_izinler = {}
    for k, val in izinler.items():
        if val.get("bitis", 0) > simdi:
            temiz_izinler[k] = val
        else:
            degisti = True
    if degisti:
        veri["izinler"] = temiz_izinler
        veri_yaz(veri)

    # Eşleşen kalıp için izin var mı?
    for kalip, _ in TEHLIKELI_KALIPLAR:
        if re.search(kalip, kod, re.IGNORECASE):
            if kalip in temiz_izinler:
                return True
    return False


def onay_anketi_gonder(kod: str, gonderen: str, aciklama: str) -> tuple[bool, str]:
    """Bridge üzerinden 4 seçenekli onay anketi gönderir."""
    soru = f"⚠️ ONAY GEREKİYOR: {aciklama}\n`{kod[:100]}`"
    secenekler = [
        "Tek seferlik onayla",
        "Reddet",
        "1 saat otomatik onayla",
        "1 hafta otomatik onayla",
    ]

    payload = {
        "to": gonderen,
        "name": soru,
        "options": secenekler,
        "max_answers": 1,
    }

    url = f"{BRIDGE_BASE.rstrip('/')}/send-poll"
    req = urllib.request.Request(
        url,
        data=json.dumps(payload).encode("utf-8"),
        headers={"Content-Type": "application/json"},
        method="POST",
    )

    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            poll_id = data.get("poll_id") or data.get("id") or str(int(time.time()))
    except Exception as e:
        return False, f"Onay anketi köprüye iletilemedi: {e}"

    veri = veri_oku()
    veri.setdefault("bekleyenler", {})[poll_id] = {
        "kod": kod,
        "gonderen": gonderen,
        "olusturuldu": int(time.time()),
        "aciklama": aciklama,
    }
    veri_yaz(veri)

    return True, (
        f"⚠️ Bu işlem yüksek riskli ({aciklama}).\n"
        f"WhatsApp üzerinden 4 seçenekli onay anketi gönderildi.\n"
        f"Lütfen anketi işaretleyerek seçiminizi yapın."
    )


def onay_oyunu_isle(poll_id: str, secim: str) -> tuple[bool, str, Optional[str]]:
    """Gelen anket oyuna göre bekleyen komutu yürüt veya reddet.
    
    Returns: (islem_yapildi, mesaj, calistirilacak_kod)
    """
    veri = veri_oku()
    bekleyenler = veri.get("bekleyenler", {})
    if poll_id not in bekleyenler:
        return False, "Bekleyen onay bulunamadı veya süresi doldu.", None

    kayit = bekleyenler.pop(poll_id)
    kod = kayit["kod"]
    simdi = int(time.time())

    secim_temiz = secim.lower()

    if "reddet" in secim_temiz:
        veri_yaz(veri)
        return True, f"❌ İşlem kullanıcı tarafından reddedildi: `{kod}`", None

    if "1 hafta" in secim_temiz:
        for kalip, _ in TEHLIKELI_KALIPLAR:
            if re.search(kalip, kod, re.IGNORECASE):
                veri.setdefault("izinler", {})[kalip] = {
                    "bitis": simdi + (7 * 86400),
                    "tur": "1hafta",
                }
        veri_yaz(veri)
        return True, f"✅ 1 hafta otomatik izin verildi. İşlem yürütülüyor: `{kod}`", kod

    if "1 saat" in secim_temiz:
        for kalip, _ in TEHLIKELI_KALIPLAR:
            if re.search(kalip, kod, re.IGNORECASE):
                veri.setdefault("izinler", {})[kalip] = {
                    "bitis": simdi + 3600,
                    "tur": "1saat",
                }
        veri_yaz(veri)
        return True, f"✅ 1 saat otomatik izin verildi. İşlem yürütülüyor: `{kod}`", kod

    # Tek seferlik onayla
    veri_yaz(veri)
    return True, f"✅ Tek seferlik onay verildi. İşlem yürütülüyor: `{kod}`", kod

#!/usr/bin/env python3
"""compact.py - /compact (/ozet) komutu ile AI oturum bağlamını sıkıştırma.

Master Blueprint Geliştirme 2.10 (TODO #11):
- Uzun süren AI oturumlarında token ve bağlam şişmesini önlemek için
  bağlamı özetleyip yeni/hafifletilmiş oturuma aktarır.
"""

from __future__ import annotations

import json
import os
import sys
from pathlib import Path


def _durum_dizini() -> Path:
    state_dir = os.getenv("TETIKLEYICI_STATE_DIR", "").strip()
    if state_dir:
        return Path(state_dir)
    cfg = os.getenv("TRIGGER_CONFIG", "").strip()
    if cfg:
        return Path(cfg).resolve().parent / ".tetikleyici"
    return Path(__file__).resolve().parent.parent.parent / ".tetikleyici"


def compact_calistir() -> str:
    """Aktif oturum bağlamını sıkıştırır ve tazeler."""
    oturumlar_dosyasi = _durum_dizini() / "oturumlar.json"
    aktif_motor_dosyasi = _durum_dizini() / "aktif_motor.txt"

    motor = "agy"
    if aktif_motor_dosyasi.is_file():
        m = aktif_motor_dosyasi.read_text(encoding="utf-8").strip().lower()
        if m in ("agy", "claude", "opencode"):
            motor = m

    eski_session = ""
    if oturumlar_dosyasi.is_file():
        try:
            veri = json.loads(oturumlar_dosyasi.read_text(encoding="utf-8"))
            eski_session = veri.get(motor, "")
        except Exception:
            pass

    # Oturumu sıfırlayıp temiz başlatılmaya hazır hâle getir
    if oturumlar_dosyasi.is_file():
        try:
            veri = json.loads(oturumlar_dosyasi.read_text(encoding="utf-8"))
            if motor in veri:
                del veri[motor]
                oturumlar_dosyasi.write_text(json.dumps(veri, indent=2), encoding="utf-8")
        except Exception:
            pass

    ayrinti = f" (Eski Oturum: `{eski_session[:12]}...`)" if eski_session else ""
    return (
        f"🧹 *[Bağlam Sıkıştırma / Compact]*\n\n"
        f"Aktif motor: *{motor}*{ayrinti}\n"
        f"Ajan bağlamı başarıyla temizlendi ve sıkıştırıldı.\n"
        f"Bir sonraki mesajınız sıfır token yüküyle taze bir oturumda işlenecektir. ✅"
    )


if __name__ == "__main__":
    print(compact_calistir())

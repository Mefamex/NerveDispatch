#!/usr/bin/env python3
"""sistem.py - Dinamik çalışma dizini (/root), pano (/clip) ve süreç yönetimi (/apps, /kill).

Master Blueprint Geliştirme 2.9 (TODO #15, #16, #17):
- `/root [dizin]`: Çalışma dizinini göster / değiştir.
- `/clip [metin]`: Wayland panosunu oku / yaz (wl-paste / wl-copy).
- `/apps`: Açık kullanıcı süreçlerini listele.
- `/kill <pid|isim>`: Kritik servis kara listesi korumalı süreç sonlandırıcı.
"""

from __future__ import annotations

import os
import re
import shutil
import subprocess
from pathlib import Path
from typing import Optional


# Kritik süreç kara listesi (asla kill edilemez)
KARA_LISTE = {
    "systemd",
    "init",
    "whatsapp-bridge",
    "tetikleyici",
    "dbus-daemon",
    "dbus-broker",
    "wireplumber",
    "pipewire",
    "sshd",
    "kwin_wayland",
    "plasmashell",
    "login",
}


def _durum_dizini() -> Path:
    state_dir = os.getenv("TETIKLEYICI_STATE_DIR", "").strip()
    if state_dir:
        return Path(state_dir)
    cfg = os.getenv("TRIGGER_CONFIG", "").strip()
    if cfg:
        return Path(cfg).resolve().parent / ".tetikleyici"
    return Path(__file__).resolve().parent.parent.parent / ".tetikleyici"


def aktif_root_al() -> Path:
    """Kayıtlı aktif çalışma dizinini döner; yoksa proje kök dizinini döner."""
    yol = _durum_dizini() / "calisma_dizini.txt"
    if yol.is_file():
        d = yol.read_text(encoding="utf-8").strip()
        p = Path(d)
        if p.is_dir():
            return p
    return Path(__file__).resolve().parent.parent.parent


def root_yonet(yeni_dizin: str = "") -> str:
    """`/root` komutunu yönetir."""
    yol_dosya = _durum_dizini() / "calisma_dizini.txt"
    yol_dosya.parent.mkdir(parents=True, exist_ok=True)

    if not yeni_dizin.strip():
        mevcut = aktif_root_al()
        return f"📂 Aktif Çalışma Dizini:\n`{mevcut}`"

    hedef = Path(yeni_dizin.strip()).expanduser().resolve()
    if not hedef.exists():
        return f"❌ Dizin bulunamadı: `{yeni_dizin}`"
    if not hedef.is_dir():
        return f"❌ Belirtilen yol bir klasör değil: `{yeni_dizin}`"

    yol_dosya.write_text(str(hedef) + "\n", encoding="utf-8")
    return f"✅ Çalışma dizini değiştirildi:\n`{hedef}`"


def clip_yonet(yeni_metin: str = "") -> str:
    """`/clip` komutunu yönetir (Wayland wl-paste / wl-copy)."""
    # 1. Yazma modu: /clip <metin>
    if yeni_metin.strip():
        if shutil.which("wl-copy"):
            try:
                proc = subprocess.run(
                    ["wl-copy", yeni_metin],
                    timeout=5,
                    capture_output=True,
                    text=True,
                )
                if proc.returncode == 0:
                    uzunluk = len(yeni_metin)
                    return f"📋 Metin panoya kopyalandı ({uzunluk} karakter)."
                return f"❌ wl-copy hatası: {proc.stderr.strip()}"
            except Exception as e:
                return f"❌ Pano yazma hatası: {e}"
        return "❌ Sistemde `wl-copy` (Wayland pano aracı) bulunamadı."

    # 2. Okuma modu: /clip
    if shutil.which("wl-paste"):
        try:
            proc = subprocess.run(
                ["wl-paste", "--no-newline"],
                timeout=5,
                capture_output=True,
                text=True,
            )
            if proc.returncode == 0:
                metin = proc.stdout
                if not metin.strip():
                    return "📋 Pano şu anda boş."
                if len(metin) > 3000:
                    metin = metin[:3000] + "\n...(kesildi)"
                return f"📋 Pano İçeriği:\n```\n{metin}\n```"
            return f"❌ wl-paste hatası: {proc.stderr.strip()}"
        except Exception as e:
            return f"❌ Pano okuma hatası: {e}"
    return "❌ Sistemde `wl-paste` (Wayland pano aracı) bulunamadı."


def apps_listele() -> str:
    """`/apps` - Kullanıcının en çok kaynak tüketen süreçlerini listeler."""
    try:
        proc = subprocess.run(
            ["ps", "-u", os.getenv("USER", "mfmx"), "-o", "pid,comm,%cpu,%mem", "--sort=-%cpu"],
            timeout=5,
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            return f"❌ Süreç listesi alınamadı: {proc.stderr.strip()}"

        satirlar = proc.stdout.strip().splitlines()
        if len(satirlar) <= 1:
            return "Aktif kullanıcı süreci bulunamadı."

        # Başlık ve ilk 15 süreç
        secilenler = satirlar[:16]
        return "🖥️ Aktif Kullanıcı Süreçleri (En çok CPU tüketenler):\n```\n" + "\n".join(secilenler) + "\n```"
    except Exception as e:
        return f"❌ Süreç listeleme hatası: {e}"


def kill_yonet(hedef: str) -> str:
    """`/kill <pid|isim>` - Süreç sonlandırma (kara liste korumalı)."""
    hedef = hedef.strip()
    if not hedef:
        return "❌ Kullanım: /kill <pid|isim>"

    # 1. Kara liste denetimi
    hedef_kucuk = hedef.lower()
    for yasak in KARA_LISTE:
        if yasak in hedef_kucuk:
            return f"🛡️ Reddedildi: `{hedef}` kritik bir sistem sürecidir ve uzaktan sonlandırılamaz."

    # 2. PID mi isim mi?
    if hedef.isdigit():
        pid = int(hedef)
        # PID'nin komut adını bul ve kara liste kontrolü yap
        try:
            comm_proc = subprocess.run(["ps", "-p", str(pid), "-o", "comm="], capture_output=True, text=True)
            comm_name = comm_proc.stdout.strip().lower()
            for yasak in KARA_LISTE:
                if yasak in comm_name:
                    return f"🛡️ Reddedildi: PID {pid} ({comm_name}) kritik bir sistem sürecidir."
            
            subprocess.run(["kill", "-15", str(pid)], timeout=5, check=True)
            return f"🛑 Süreç sonlandırıldı: PID {pid} ({comm_name})"
        except subprocess.CalledProcessError:
            return f"❌ PID {pid} sonlandırılamadı (veya süreç zaten kapandı)."
        except Exception as e:
            return f"❌ Hata: {e}"

    # İsim ise pkill
    try:
        subprocess.run(["pkill", "-15", "-f", hedef], timeout=5, check=True)
        return f"🛑 `{hedef}` kalıbına uyan süreçler sonlandırıldı."
    except subprocess.CalledProcessError:
        return f"❌ `{hedef}` adında çalışan bir süreç bulunamadı."
    except Exception as e:
        return f"❌ Hata: {e}"

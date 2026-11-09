"""test_sistem.py - /root, /clip, /apps, /kill sistem yönetimi testleri (TODO #15, #16, #17)."""

from __future__ import annotations

from pathlib import Path
import sys
from unittest.mock import MagicMock, patch

import pytest

CMD_DIR = str(Path(__file__).resolve().parent.parent / "tetikleyici" / "cmd")
if CMD_DIR not in sys.path:
    sys.path.insert(0, CMD_DIR)

from sozdizimi import parse_cmd
import sistem


def test_parse_cmd_sistem():
    r1 = parse_cmd("/root")
    assert r1.komut_turu == "root"
    assert r1.kod == ""

    r2 = parse_cmd("/root /home/mfmx/proj")
    assert r2.komut_turu == "root"
    assert r2.kod == "/home/mfmx/proj"

    r3 = parse_cmd("/clip")
    assert r3.komut_turu == "clip"
    assert r3.metin == ""

    r4 = parse_cmd("/clip panoya yazılacak metin")
    assert r4.komut_turu == "clip"
    assert r4.metin == "panoya yazılacak metin"

    r5 = parse_cmd("/apps")
    assert r5.komut_turu == "apps"

    r6 = parse_cmd("/kill 45890")
    assert r6.komut_turu == "kill"
    assert r6.hedef == "45890"

    r7 = parse_cmd("/kill firefox")
    assert r7.komut_turu == "kill"
    assert r7.hedef == "firefox"


def test_root_yonet(tmp_path, monkeypatch):
    durum_dir = tmp_path / ".tetikleyici"
    monkeypatch.setattr(sistem, "_durum_dizini", lambda: durum_dir)

    # Başlangıç
    assert "Aktif Çalışma Dizini" in sistem.root_yonet()

    # Olmayan dizin
    yanit_hata = sistem.root_yonet("/tmp/olmayan_dizin_xyz_123")
    assert "bulunamadı" in yanit_hata

    # Geçerli dizin
    hedef_dir = tmp_path / "yeni_proje"
    hedef_dir.mkdir()
    yanit_ok = sistem.root_yonet(str(hedef_dir))
    assert "değiştirildi" in yanit_ok

    # Kaydedilmiş mi kontrol et
    assert sistem.aktif_root_al() == hedef_dir


def test_clip_yonet_wl_copy_paste():
    with patch("shutil.which", return_value="/usr/bin/wl-copy"):
        proc = MagicMock()
        proc.returncode = 0
        with patch("subprocess.run", return_value=proc):
            yanit = sistem.clip_yonet("test panosu")
            assert "kopyalandı" in yanit

    with patch("shutil.which", return_value="/usr/bin/wl-paste"):
        proc2 = MagicMock()
        proc2.returncode = 0
        proc2.stdout = "panodaki eski yazi"
        with patch("subprocess.run", return_value=proc2):
            yanit2 = sistem.clip_yonet()
            assert "panodaki eski yazi" in yanit2


def test_kill_kara_liste():
    # Kara liste isim testi
    yanit = sistem.kill_yonet("systemd")
    assert "Reddedildi" in yanit

    yanit_bridge = sistem.kill_yonet("whatsapp-bridge")
    assert "Reddedildi" in yanit_bridge

    yanit_kwin = sistem.kill_yonet("kwin_wayland")
    assert "Reddedildi" in yanit_kwin

"""test_onay.py - İnteraktif onay ve izin sistemi testleri (TODO #9)."""

from __future__ import annotations

import json
from pathlib import Path
import sys
import time
from unittest.mock import MagicMock, patch

import pytest

CMD_DIR = str(Path(__file__).resolve().parent.parent / "tetikleyici" / "cmd")
if CMD_DIR not in sys.path:
    sys.path.insert(0, CMD_DIR)

import onay


def test_tehlikeli_mi_tespit():
    assert onay.tehlikeli_mi("rm -rf /tmp/test")[0] is True
    assert onay.tehlikeli_mi("git push origin main --force")[0] is True
    assert onay.tehlikeli_mi("systemctl stop nginx")[0] is True
    assert onay.tehlikeli_mi("sudo reboot")[0] is True
    assert onay.tehlikeli_mi("echo merhaba")[0] is False
    assert onay.tehlikeli_mi("ls -la /home/mfmx")[0] is False


def test_onay_anketi_ve_oy_isle(tmp_path, monkeypatch):
    izin_dosya = tmp_path / "izinler.json"
    monkeypatch.setattr(onay, "_izinler_dosyasi", lambda: izin_dosya)

    mock_resp = MagicMock()
    mock_resp.read.return_value = json.dumps({"ok": True, "poll_id": "poll-456"}).encode("utf-8")
    mock_resp.__enter__.return_value = mock_resp

    with patch("urllib.request.urlopen", return_value=mock_resp):
        ok, msg = onay.onay_anketi_gonder("rm -rf /tmp/cop", "905550000001", "Toplu dosya silme")
        assert ok is True
        assert "onay anketi gönderildi" in msg

    veri = json.loads(izin_dosya.read_text(encoding="utf-8"))
    assert "poll-456" in veri["bekleyenler"]

    # Reddetme testi
    islem, bildirim, kod = onay.onay_oyunu_isle("poll-456", "Reddet")
    assert islem is True
    assert "reddedildi" in bildirim
    assert kod is None

    # Yeniden anket oluştur ve 1 saat otomatik onayla
    with patch("urllib.request.urlopen", return_value=mock_resp):
        onay.onay_anketi_gonder("rm -rf /tmp/cop2", "905550000001", "Toplu dosya silme")

    islem2, bildirim2, kod2 = onay.onay_oyunu_isle("poll-456", "1 saat otomatik onayla")
    assert islem2 is True
    assert kod2 == "rm -rf /tmp/cop2"

    # Şimdi izinli mi kontrolü
    assert onay.izinli_mi("rm -rf /tmp/cop3") is True

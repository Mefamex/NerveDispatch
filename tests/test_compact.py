"""test_compact.py - /compact (/ozet) bağlam sıkıştırma testleri (TODO #11)."""

from __future__ import annotations

import json
from pathlib import Path
import sys

import pytest

CMD_DIR = str(Path(__file__).resolve().parent.parent / "tetikleyici" / "cmd")
if CMD_DIR not in sys.path:
    sys.path.insert(0, CMD_DIR)

from sozdizimi import parse_cmd
import compact


def test_parse_cmd_compact():
    assert parse_cmd("/compact").komut_turu == "compact"
    assert parse_cmd("/ozet").komut_turu == "compact"
    assert parse_cmd("  /COMPACT  ").komut_turu == "compact"


def test_compact_calistir(tmp_path, monkeypatch):
    durum_dir = tmp_path / ".tetikleyici"
    durum_dir.mkdir(parents=True)
    monkeypatch.setattr(compact, "_durum_dizini", lambda: durum_dir)

    (durum_dir / "aktif_motor.txt").write_text("agy\n", encoding="utf-8")
    (durum_dir / "oturumlar.json").write_text(json.dumps({"agy": "sess-xyz-123", "claude": "sess-c-456"}), encoding="utf-8")

    yanit = compact.compact_calistir()
    assert "Bağlam Sıkıştırma" in yanit
    assert "agy" in yanit

    # agy oturumu silinmiş olmalı, claude kalmalı
    kalan = json.loads((durum_dir / "oturumlar.json").read_text(encoding="utf-8"))
    assert "agy" not in kalan
    assert kalan.get("claude") == "sess-c-456"

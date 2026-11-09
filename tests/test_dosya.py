"""test_dosya.py - /get komutu ve çift yönlü dosya transferi testleri (TODO #14)."""

from __future__ import annotations

import base64
import json
from pathlib import Path
import stat
import sys
from unittest.mock import MagicMock, patch

import pytest

CMD_DIR = str(Path(__file__).resolve().parent.parent / "tetikleyici" / "cmd")
if CMD_DIR not in sys.path:
    sys.path.insert(0, CMD_DIR)

from sozdizimi import parse_cmd
import dosya


def test_parse_cmd_get():
    res = parse_cmd("/get /home/user/document.pdf")
    assert res.komut_turu == "get"
    assert res.hedef == "/home/user/document.pdf"

    res2 = parse_cmd("  /GET   /tmp/data.csv  ")
    assert res2.komut_turu == "get"
    assert res2.hedef == "/tmp/data.csv"


def test_dosya_gonder_olmayan_dosya():
    yanit = dosya.dosya_gonder("/tmp/olmayan_dosya_xyz123.txt", "905550000001")
    assert "bulunamadı" in yanit


def test_dosya_gonder_basarili(tmp_path):
    test_file = tmp_path / "ornek.txt"
    test_file.write_text("merhaba dunya!", encoding="utf-8")

    mock_resp = MagicMock()
    mock_resp.status = 200
    mock_resp.__enter__.return_value = mock_resp

    with patch("urllib.request.urlopen", return_value=mock_resp) as mock_open:
        yanit = dosya.dosya_gonder(str(test_file), "905550000001", bridge="http://127.0.0.1:31415")
        assert "başarıyla gönderildi" in yanit
        assert mock_open.called
        req = mock_open.call_args[0][0]
        assert req.get_full_url() == "http://127.0.0.1:31415/send-media"
        body = json.loads(req.data.decode("utf-8"))
        assert body["to"] == "905550000001"
        assert body["filename"] == "ornek.txt"
        assert base64.b64decode(body["media"]) == b"merhaba dunya!"


def test_gelen_dosyayi_kaydet_guvenli_izin(tmp_path):
    veri = b"zararsiz metin icerigi"
    hedef_yol, boyut = dosya.gelen_dosyayi_kaydet(veri, "script.sh", alt_dizin=tmp_path)

    assert hedef_yol.exists()
    assert boyut == len(veri)
    assert hedef_yol.read_bytes() == veri

    # İzin kontrolü: çalıştırılabilir (+x) OLMAMALI (0o644)
    mode = stat.S_IMODE(hedef_yol.stat().st_mode)
    assert mode & stat.S_IXUSR == 0
    assert mode & stat.S_IXGRP == 0
    assert mode & stat.S_IXOTH == 0

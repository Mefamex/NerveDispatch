"""test_ekran.py - /ekran komutu ve ekran görüntüsü boru hattı testleri."""

from __future__ import annotations

import base64
import json
from pathlib import Path
import sys
from unittest.mock import MagicMock, patch

import pytest

CMD_DIR = str(Path(__file__).resolve().parent.parent / "tetikleyici" / "cmd")
if CMD_DIR not in sys.path:
    sys.path.insert(0, CMD_DIR)

from sozdizimi import parse_cmd
import ekran


def test_parse_cmd_ekran():
    res1 = parse_cmd("/ekran")
    assert res1.komut_turu == "ekran"
    assert res1.saniye == 0

    res2 = parse_cmd("/ekran 5")
    assert res2.komut_turu == "ekran"
    assert res2.saniye == 5

    res3 = parse_cmd("  /EKRAN   10  ")
    assert res3.komut_turu == "ekran"
    assert res3.saniye == 10


def test_ekran_yakala_spectacle(tmp_path):
    hedef = tmp_path / "ekran.png"

    def fake_run(cmd, **kwargs):
        hedef.write_bytes(b"\x89PNG\r\n\x1a\nfakeimage")
        proc = MagicMock()
        proc.returncode = 0
        return proc

    with patch("shutil.which", return_value="/usr/bin/spectacle"):
        with patch("subprocess.run", side_effect=fake_run) as mock_run:
            ok, err = ekran.yakala(hedef, gecikme_saniye=0)
            assert ok is True
            assert err == ""
            assert mock_run.called
            assert mock_run.call_args[0][0] == ["spectacle", "-b", "-n", "-o", str(hedef)]


def test_ekran_gonder_basarili(tmp_path):
    png = tmp_path / "test.png"
    png.write_bytes(b"test-png-content")

    mock_resp = MagicMock()
    mock_resp.status = 200
    mock_resp.__enter__.return_value = mock_resp

    with patch("urllib.request.urlopen", return_value=mock_resp) as mock_open:
        ok, err = ekran.gonder(png, alici="905550000001", bridge="http://127.0.0.1:31415")
        assert ok is True
        assert err == ""
        assert mock_open.called
        req = mock_open.call_args[0][0]
        assert req.get_full_url() == "http://127.0.0.1:31415/send-media"
        body = json.loads(req.data.decode("utf-8"))
        assert body["to"] == "905550000001"
        assert body["mime_type"] == "image/png"
        assert base64.b64decode(body["media"]) == b"test-png-content"


def test_ekran_ana_akis(tmp_path):
    with patch.object(ekran, "yakala", return_value=(True, "")) as mock_yakala:
        with patch.object(ekran, "gonder", return_value=(True, "")) as mock_gonder:
            cikti = ekran.ana(alici="905550000001", saniye=2)
            assert "gönderildi" in cikti
            assert mock_yakala.called
            assert mock_gonder.called

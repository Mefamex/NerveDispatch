"""test_tts.py - /tts komutu ve Piper sesli yanıt boru hattı testleri (TODO #18)."""

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

PROJE_ROOT = str(Path(__file__).resolve().parent.parent)
if PROJE_ROOT not in sys.path:
    sys.path.insert(0, PROJE_ROOT)

from sozdizimi import parse_cmd
import tts_pipeline


def test_parse_cmd_tts():
    r1 = parse_cmd("/tts")
    assert r1.komut_turu == "tts"
    assert r1.metin == ""

    r2 = parse_cmd("/tts on")
    assert r2.komut_turu == "tts"
    assert r2.metin == "on"

    r3 = parse_cmd("/tts off")
    assert r3.komut_turu == "tts"
    assert r3.metin == "off"

    r4 = parse_cmd("/tts merhaba nasılsın")
    assert r4.komut_turu == "tts"
    assert r4.metin == "merhaba nasılsın"


def test_tts_durum_ayarla(tmp_path, monkeypatch):
    durum_dir = tmp_path / ".tetikleyici"
    monkeypatch.setattr(tts_pipeline, "_durum_dizini", lambda: durum_dir)

    assert tts_pipeline.tts_aktif_mi() is False
    tts_pipeline.tts_durum_ayarla(True)
    assert tts_pipeline.tts_aktif_mi() is True
    tts_pipeline.tts_durum_ayarla(False)
    assert tts_pipeline.tts_aktif_mi() is False


def test_sentezle_ve_gonder_mock(tmp_path):
    mock_resp = MagicMock()
    mock_resp.status = 200
    mock_resp.__enter__.return_value = mock_resp

    def fake_run(cmd, **kwargs):
        p = MagicMock()
        p.returncode = 0
        if "ffmpeg" in cmd[0]:
            out_file = cmd[-1]
            Path(out_file).write_bytes(b"OggSfake_opus_data")
        elif "espeak" in cmd[0]:
            out_idx = cmd.index("-w")
            Path(cmd[out_idx + 1]).write_bytes(b"RIFFfake_wav_data")
        elif "piper" in cmd[0]:
            out_idx = cmd.index("--output_file")
            Path(cmd[out_idx + 1]).write_bytes(b"RIFFfake_wav_data")
        return p

    with patch("shutil.which") as mock_which:
        mock_which.side_effect = lambda prog: f"/usr/bin/{prog}"
        with patch("subprocess.run", side_effect=fake_run):
            with patch("urllib.request.urlopen", return_value=mock_resp) as mock_open:
                ok, msg = tts_pipeline.sentezle_ve_gonder("merhaba dünya", alici="905550000001", bridge="http://127.0.0.1:31415")
                assert ok is True
                assert "Sesli not gönderildi" in msg
                assert mock_open.called
                req = mock_open.call_args[0][0]
                body = json.loads(req.data.decode("utf-8"))
                assert body["to"] == "905550000001"
                assert body["ptt"] is True
                assert body["mime_type"] == "audio/ogg; codecs=opus"
                assert base64.b64decode(body["media"]) == b"OggSfake_opus_data"

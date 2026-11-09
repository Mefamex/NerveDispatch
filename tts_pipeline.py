#!/usr/bin/env python3
"""tts_pipeline.py - Metinden Türkçe sesli yanıt üretimi (Piper TTS + ffmpeg ogg opus).

Master Blueprint Geliştirme 2.11 (TODO #18):
- Piper TTS ile CPU üzerinde sıfır VRAM tüketerek Türkçe ses sentezler.
- ffmpeg ile WhatsApp uyumlu OGG Opus (.ogg, 48kHz, mono) sesli nota dönüştürür.
- Köprü POST /send-media ile sesli not (ptt: true) olarak WhatsApp'a iletir.
"""

from __future__ import annotations

import base64
import json
import os
import shutil
import subprocess
import sys
import tempfile
import urllib.request
from pathlib import Path
from typing import Optional


BRIDGE_BASE = os.getenv("BRIDGE_BASE", "http://127.0.0.1:31415")
MODEL_YOLU = os.getenv("PIPER_MODEL", "/mnt/DiskD/SOFTWARE/models/piper/tr_TR-dfki-medium.onnx")


def _durum_dizini() -> Path:
    state_dir = os.getenv("TETIKLEYICI_STATE_DIR", "").strip()
    if state_dir:
        return Path(state_dir)
    cfg = os.getenv("TRIGGER_CONFIG", "").strip()
    if cfg:
        return Path(cfg).resolve().parent / ".tetikleyici"
    return Path(__file__).resolve().parent / ".tetikleyici"


def tts_aktif_mi() -> bool:
    yol = _durum_dizini() / "tts_aktif.txt"
    if yol.is_file():
        return yol.read_text(encoding="utf-8").strip().lower() in ("true", "1", "on", "acik")
    return False


def tts_durum_ayarla(acik: bool) -> str:
    yol = _durum_dizini() / "tts_aktif.txt"
    yol.parent.mkdir(parents=True, exist_ok=True)
    if acik:
        yol.write_text("true\n", encoding="utf-8")
        return "🔊 Sesli yanıt modu açıldı (/tts on). Cevaplar artık sesli not olarak da iletilecek."
    else:
        yol.write_text("false\n", encoding="utf-8")
        return "🔇 Sesli yanıt modu kapatıldı (/tts off). Cevaplar normal metin olarak iletilecek."


def sentezle_ve_gonder(metin: str, alici: str, bridge: str = BRIDGE_BASE) -> tuple[bool, str]:
    """Metni seslendirir ve WhatsApp'a sesli not olarak gönderir."""
    if not metin.strip():
        return False, "Seslendirilecek metin boş."

    piper_cmd = shutil.which("piper")
    espeak_cmd = shutil.which("espeak-ng")

    if not piper_cmd and not espeak_cmd:
        return False, "Sistemde Piper TTS (`piper`) veya `espeak-ng` bulunamadı. Kurulum: `pip install piper-tts`"

    with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as tmp_wav, \
         tempfile.NamedTemporaryFile(suffix=".ogg", delete=False) as tmp_ogg:
        wav_path = Path(tmp_wav.name)
        ogg_path = Path(tmp_ogg.name)

    try:
        # 1. Ses sentezi (Piper veya fallback espeak)
        if piper_cmd and Path(MODEL_YOLU).exists():
            cmd = [piper_cmd, "--model", MODEL_YOLU, "--output_file", str(wav_path)]
            proc = subprocess.run(cmd, input=metin, text=True, capture_output=True, timeout=30)
            if proc.returncode != 0:
                return False, f"Piper hatası: {proc.stderr.strip()[:200]}"
        elif espeak_cmd:
            cmd = [espeak_cmd, "-v", "tr", "-w", str(wav_path), metin]
            proc = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
            if proc.returncode != 0:
                return False, f"espeak-ng hatası: {proc.stderr.strip()[:200]}"
        elif piper_cmd:
            cmd = [piper_cmd, "--output_file", str(wav_path)]
            proc = subprocess.run(cmd, input=metin, text=True, capture_output=True, timeout=30)
            if proc.returncode != 0:
                return False, f"Piper sentez hatası: {proc.stderr.strip()[:200]}"

        if not wav_path.exists() or wav_path.stat().st_size == 0:
            return False, "WAV ses dosyası üretilemedi."

        # 2. ffmpeg ile WhatsApp uyumlu OGG Opus formatına çevir
        ffmpeg_cmd = shutil.which("ffmpeg")
        if not ffmpeg_cmd:
            return False, "Sistemde ffmpeg bulunamadı."

        ff_proc = subprocess.run(
            [ffmpeg_cmd, "-y", "-i", str(wav_path), "-c:a", "libopus", "-b:a", "32k", "-ar", "48000", str(ogg_path)],
            capture_output=True,
            text=True,
            timeout=20,
        )
        if ff_proc.returncode != 0:
            return False, f"ffmpeg opus dönüştürme hatası: {ff_proc.stderr.strip()[:200]}"

        # 3. Köprü üzerinden POST /send-media (ptt: true) ile gönder
        ses_baytlari = ogg_path.read_bytes()
        govde = json.dumps(
            {
                "to": alici,
                "media": base64.b64encode(ses_baytlari).decode("ascii"),
                "mime_type": "audio/ogg; codecs=opus",
                "ptt": True,
            }
        ).encode("utf-8")

        url = f"{bridge.rstrip('/')}/send-media"
        req = urllib.request.Request(
            url,
            data=govde,
            headers={"Content-Type": "application/json"},
            method="POST",
        )
        with urllib.request.urlopen(req, timeout=15) as resp:
            if resp.status == 200:
                return True, "🎙️ Sesli not gönderildi."
            return False, f"Köprü HTTP {resp.status}"

    except Exception as e:
        return False, f"TTS gönderme hatası: {e}"
    finally:
        wav_path.unlink(missing_ok=True)
        ogg_path.unlink(missing_ok=True)


if __name__ == "__main__":
    if len(sys.argv) > 1:
        komut = sys.argv[1].lower()
        if komut == "on":
            print(tts_durum_ayarla(True))
        elif komut == "off":
            print(tts_durum_ayarla(False))
        else:
            metin = " ".join(sys.argv[1:])
            ok, msj = sentezle_ve_gonder(metin, alici="me")
            print(msj)

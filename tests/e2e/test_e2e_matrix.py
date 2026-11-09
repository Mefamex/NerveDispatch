#!/usr/bin/env python3
"""FAZ 4.3 - Uçtan uca test matrisi.

Derlenmiş GERÇEK `bin/tetikleyici` ikilisi, Python ile yazılmış sahte bir köprüye
(`/events` SSE, `/send`, `/messages`, `/health`, `/mark-read`, `/reaction`) bağlanır.
Canlı WhatsApp oturumuna, gerçek `tetikleyici.yaml`'a ve `.tetikleyici/` durumuna
DOKUNMAZ: config, durum dizini ve köprü adresi geçici dizine yönlendirilir.

Matris:
  kanal    : whatsapp, telegram (aynı kural motoru, kanal filtresi)
  eşleşme  : metin tutar / tutmaz / not_text dışlaması
  teslim   : komut çıktısı sohbete döner; aynı mesaj kimliği iki kez işlenmez
  kurtarma : köprü akışı koparsa yeniden bağlanıp sonraki mesajı işler;
             bozuk SSE satırı akışı kesmez

Çalıştırma: make e2e   (ikili yoksa önce `make build`)
"""

import json
import os
import queue
import shutil
import signal
import socket
import subprocess
import tempfile
import threading
import time
import unittest
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
BINARY = ROOT / "bin" / "tetikleyici"
SAHIP = "905551234567@s.whatsapp.net"


class SahteKopru:
    """Tetikleyicinin konuştuğu köprü uçlarının en küçük taklidi."""

    def __init__(self):
        self.gonderilen = []          # /send gövdeleri
        self.olaylar = queue.Queue()  # /events akışına yazılacak ham satırlar
        self.akis_baglanti = 0        # /events bağlantı sayısı
        self.kilit = threading.Lock()
        self.kopar = threading.Event()
        kopru = self

        class Isleyici(BaseHTTPRequestHandler):
            protocol_version = "HTTP/1.1"

            def log_message(self, *a):
                pass

            def _json(self, kod, govde):
                veri = json.dumps(govde).encode()
                self.send_response(kod)
                self.send_header("Content-Type", "application/json")
                self.send_header("Content-Length", str(len(veri)))
                self.end_headers()
                self.wfile.write(veri)

            def do_GET(self):
                if self.path.startswith("/health"):
                    return self._json(200, {"paired": True, "connected": True})
                if self.path.startswith("/messages/"):
                    return self._json(200, {"count": 0, "messages": []})
                if self.path.startswith("/events"):
                    return self._akis()
                return self._json(404, {"error": "yok"})

            def do_POST(self):
                uzunluk = int(self.headers.get("Content-Length") or 0)
                govde = self.rfile.read(uzunluk) if uzunluk else b"{}"
                if self.path.startswith("/send"):
                    with kopru.kilit:
                        kopru.gonderilen.append(json.loads(govde or b"{}"))
                return self._json(200, {"ok": True})

            def _akis(self):
                with kopru.kilit:
                    kopru.akis_baglanti += 1
                self.send_response(200)
                self.send_header("Content-Type", "text/event-stream")
                self.send_header("Cache-Control", "no-cache")
                self.send_header("Connection", "close")
                self.end_headers()
                self.close_connection = True
                try:
                    while True:
                        if kopru.kopar.is_set():
                            kopru.kopar.clear()
                            return  # akışı sunucu tarafında kes
                        try:
                            satir = kopru.olaylar.get(timeout=0.1)
                        except queue.Empty:
                            continue
                        self.wfile.write((satir + "\n\n").encode())
                        self.wfile.flush()
                except (BrokenPipeError, ConnectionResetError):
                    return

        with socket.socket() as s:
            s.bind(("127.0.0.1", 0))
            self.port = s.getsockname()[1]
        self.sunucu = ThreadingHTTPServer(("127.0.0.1", self.port), Isleyici)
        self.sunucu.daemon_threads = True
        threading.Thread(target=self.sunucu.serve_forever, daemon=True).start()

    @property
    def adres(self):
        return f"http://127.0.0.1:{self.port}"

    def mesaj(self, mid, metin, kanal="whatsapp", kimden=SAHIP):
        olay = {
            "event": "message",
            "timestamp": int(time.time()),
            "data": {
                "channel": kanal,
                "id": mid,
                "type": "text",
                "direction": "incoming",
                "from": kimden,
                "sender": kimden,
                "timestamp": int(time.time()),
                "text": metin,
                "is_from_me": False,
                "is_group": False,
                "push_name": "Test",
            },
        }
        self.olaylar.put("data: " + json.dumps(olay))

    def ham(self, satir):
        self.olaylar.put(satir)

    def metinler(self):
        with self.kilit:
            return [g.get("text", "") for g in self.gonderilen]

    def kapat(self):
        self.sunucu.shutdown()
        self.sunucu.server_close()


KURALLAR = """
debug_cikti: false
okundu_isaretle: false
ai_eti_ekle: false
catchup_max_age: 0s
rules:
  - name: ping
    match:
      text: ping
      not_text: sessiz
    do:
      - type: command
        code: "printf 'pong'"
        reply: true
  - name: telegram-yalniz
    match:
      channel: telegram
      text: durum
    do:
      - type: command
        code: "printf 'telegram-tamam'"
        reply: true
"""


def bekle(kosul, sure=10.0):
    son = time.monotonic() + sure
    while time.monotonic() < son:
        if kosul():
            return True
        time.sleep(0.05)
    return kosul()


@unittest.skipUnless(BINARY.exists(), "bin/tetikleyici yok: önce `make build`")
class UctanUcaMatris(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        if not shutil.which("bash"):
            raise unittest.SkipTest("bash yok")

    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        base = Path(self.tmp.name)
        self.config = base / "tetikleyici.yaml"
        self.config.write_text(KURALLAR)
        self.state = base / "state"
        self.state.mkdir()
        self.kopru = SahteKopru()
        env = dict(
            os.environ,
            HOME=str(base),  # ~/ altındaki gerçek dosyalara asla yazmasın
            TETIKLEYICI_CONFIG=str(self.config),
            TETIKLEYICI_STATE_DIR=str(self.state),
        )
        self.log = open(base / "stdout.log", "w")
        self.surec = subprocess.Popen(
            [
                str(BINARY),
                "--bridge", self.kopru.adres,
                "--config", str(self.config),
                "--seen", str(self.state / "seen.txt"),
                "--log", str(self.state / "tetikleyici.log"),
                "--history", "0",
            ],
            cwd=str(base),
            env=env,
            stdout=self.log,
            stderr=subprocess.STDOUT,
            start_new_session=True,
        )
        self.assertTrue(
            bekle(lambda: self.kopru.akis_baglanti >= 1),
            "tetikleyici köprü akışına bağlanmadı:\n" + self._gunluk(),
        )

    def tearDown(self):
        if self.surec.poll() is None:
            os.killpg(self.surec.pid, signal.SIGTERM)
            try:
                self.surec.wait(timeout=5)
            except subprocess.TimeoutExpired:
                os.killpg(self.surec.pid, signal.SIGKILL)
        self.log.close()
        self.kopru.kapat()
        self.tmp.cleanup()

    def _gunluk(self):
        self.log.flush()
        metin = (Path(self.tmp.name) / "stdout.log").read_text(errors="replace")
        dosya = self.state / "tetikleyici.log"
        if dosya.exists():
            metin += "\n--- tetikleyici.log ---\n" + dosya.read_text(errors="replace")
        return metin[-3000:]

    def _cevap_bekle(self, metin, adet=1, sure=10.0):
        ok = bekle(lambda: self.kopru.metinler().count(metin) >= adet, sure)
        self.assertTrue(ok, f"'{metin}' cevabı gelmedi; gönderilen={self.kopru.metinler()}\n{self._gunluk()}")

    # --- eşleşme ve teslim ---
    def test_whatsapp_eslesen_mesaja_komut_ciktisi_doner(self):
        self.kopru.mesaj("E2E-1", "ping")
        self._cevap_bekle("pong")
        self.assertEqual(self.kopru.gonderilen[0]["to"], SAHIP)

    def test_eslesmeyen_ve_dislanan_mesaj_cevapsiz_kalir(self):
        self.kopru.mesaj("E2E-2", "merhaba")
        self.kopru.mesaj("E2E-3", "ping sessiz")
        self.kopru.mesaj("E2E-4", "ping")  # işaretçi: bu gelince öncekiler çoktan işlendi
        self._cevap_bekle("pong")
        time.sleep(0.5)
        self.assertEqual(self.kopru.metinler(), ["pong"])

    def test_ayni_mesaj_kimligi_iki_kez_islenmez(self):
        self.kopru.mesaj("E2E-5", "ping")
        self._cevap_bekle("pong")
        self.kopru.mesaj("E2E-5", "ping")
        self.kopru.mesaj("E2E-6", "ping")
        self._cevap_bekle("pong", adet=2)
        time.sleep(0.5)
        self.assertEqual(self.kopru.metinler().count("pong"), 2)

    # --- kanal matrisi ---
    def test_telegram_kanali_kendi_kuralini_tetikler(self):
        self.kopru.mesaj("E2E-7", "durum", kanal="telegram", kimden="975678139")
        self._cevap_bekle("telegram-tamam")

    def test_kanal_filtresi_whatsapp_mesajini_telegram_kuralina_sokmaz(self):
        self.kopru.mesaj("E2E-8", "durum")
        self.kopru.mesaj("E2E-9", "ping")
        self._cevap_bekle("pong")
        time.sleep(0.5)
        self.assertNotIn("telegram-tamam", self.kopru.metinler())

    # --- hata kurtarma ---
    def test_bozuk_sse_satiri_akisi_kesmez(self):
        self.kopru.ham("data: {bozuk json")
        self.kopru.ham(": yorum satiri")
        self.kopru.mesaj("E2E-10", "ping")
        self._cevap_bekle("pong")

    def test_akis_koparsa_yeniden_baglanip_devam_eder(self):
        self.kopru.kopar.set()
        self.assertTrue(
            bekle(lambda: self.kopru.akis_baglanti >= 2, 30.0),
            "akış koptuktan sonra yeniden bağlanmadı\n" + self._gunluk(),
        )
        self.kopru.mesaj("E2E-11", "ping")
        self._cevap_bekle("pong", sure=15.0)


if __name__ == "__main__":
    unittest.main()

#!/usr/bin/env python3
"""/agent durum: her motorun bağlı sohbet kimliğini gösterir (madde 23)."""

import importlib.util
import json
import os
import sys
import tempfile
import unittest
from pathlib import Path

KOK = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(KOK / "tetikleyici" / "cmd"))

# Standart kütüphanedeki `cmd` modülüyle çakışmasın diye dosyadan yüklenir.
_spec = importlib.util.spec_from_file_location("tetik_cmd_agent", KOK / "tetikleyici" / "cmd" / "cmd.py")
cmd_py = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(cmd_py)


class AgentSohbetTestleri(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.eski = os.environ.get("TETIKLEYICI_STATE_DIR")
        os.environ["TETIKLEYICI_STATE_DIR"] = self.tmp.name
        self.dizin = Path(self.tmp.name)

    def tearDown(self):
        if self.eski is None:
            os.environ.pop("TETIKLEYICI_STATE_DIR", None)
        else:
            os.environ["TETIKLEYICI_STATE_DIR"] = self.eski
        self.tmp.cleanup()

    def test_kayitli_sohbetler_listelenir(self):
        (self.dizin / "aktif_motor.txt").write_text("claude\n")
        (self.dizin / "oturumlar.json").write_text(json.dumps({
            "agy": {"conversation_id": "8372b68f-a31b-487d-bef9-144284811a30"},
            "claude": {"session_id": "a9ce2f25-b315-46bd-ad5b-b3434c9246d1"},
            "opencode": {"session_id": ""},
        }))
        cikti = cmd_py.agent_yonet("status")
        self.assertIn("Aktif Ajan Motoru: *claude*", cikti)
        self.assertIn("8372b68f-a31b-487d-bef9-144284811a30", cikti)
        self.assertIn("a9ce2f25-b315-46bd-ad5b-b3434c9246d1", cikti)
        self.assertIn("claude (aktif) ✅", cikti)
        opencode_satiri = cikti.split("• opencode")[1].split("\n")[1]
        self.assertIn("yeni sohbet", opencode_satiri)

    def test_dosya_yoksa_bozulmaz(self):
        cikti = cmd_py.agent_yonet("status")
        self.assertIn("Aktif Ajan Motoru: *agy*", cikti)
        self.assertEqual(cikti.count("yeni sohbet"), 3)


if __name__ == "__main__":
    unittest.main()

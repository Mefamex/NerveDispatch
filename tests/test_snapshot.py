#!/usr/bin/env python3
"""scripts/nerve-snapshot.sh (FAZ 4.2) testleri.

Gerçek session.db'ye ve proje durumuna DOKUNMAZ: tüm yollar env ile tmp'e yönlendirilir.
Doğrulanan: WAL modunda açık bir veritabanının tutarlı yedeği, durum dosyalarının
kopyalanması (yedek klasörleri hariç), latest bağlantısı, saklama politikası,
bozuk veritabanında hata ve yarım yedek bırakmama.
"""

import os
import sqlite3
import subprocess
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SCRIPT = ROOT / "scripts" / "nerve-snapshot.sh"


class SnapshotTestleri(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        base = Path(self.tmp.name)
        self.session = base / "session"
        self.project = base / "proje"
        self.snaps = base / "snaps"
        self.session.mkdir()
        (self.project / ".tetikleyici" / "zamanlanmis").mkdir(parents=True)
        (self.project / ".tetikleyici" / "backups" / "eski").mkdir(parents=True)
        (self.project / ".tetikleyici" / "seen.txt").write_text("m1\nm2\n")
        (self.project / ".tetikleyici" / "zamanlanmis" / "1.json").write_text("{}")
        (self.project / ".tetikleyici" / "backups" / "eski" / "x").write_text("x")
        (self.project / ".tetikleyici" / "tetikleyici.log").write_text("log")
        (self.project / "tetikleyici.yaml").write_text("kurallar: []\n")
        (self.session / "messages.jsonl").write_text('{"id":1}\n')

        # WAL modunda, kapatılmamış (canlı) bir bağlantı: değişiklikler -wal dosyasında
        self.db = sqlite3.connect(self.session / "session.db")
        self.db.execute("PRAGMA journal_mode=WAL")
        self.db.execute("CREATE TABLE t(v TEXT)")
        self.db.executemany("INSERT INTO t VALUES (?)", [(str(i),) for i in range(500)])
        self.db.commit()

    def tearDown(self):
        self.db.close()
        self.tmp.cleanup()

    def calistir(self, keep="14"):
        env = dict(
            os.environ,
            NERVE_SESSION_DIR=str(self.session),
            NERVE_PROJECT_ROOT=str(self.project),
            NERVE_SNAPSHOT_DIR=str(self.snaps),
            NERVE_SNAPSHOT_KEEP=keep,
        )
        return subprocess.run(["bash", str(SCRIPT)], env=env, capture_output=True, text=True)

    def test_wal_acikken_tutarli_yedek(self):
        self.assertTrue((self.session / "session.db-wal").exists())
        r = self.calistir()
        self.assertEqual(r.returncode, 0, r.stderr)
        latest = self.snaps / "latest"
        kopya = sqlite3.connect(latest / "session.db")
        self.assertEqual(kopya.execute("SELECT count(*) FROM t").fetchone()[0], 500)
        self.assertEqual(kopya.execute("PRAGMA integrity_check").fetchone()[0], "ok")
        kopya.close()

    def test_durum_dosyalari_kopyalanir_yedekler_ve_loglar_haric(self):
        self.assertEqual(self.calistir().returncode, 0)
        state = self.snaps / "latest" / "state"
        self.assertEqual((state / "seen.txt").read_text(), "m1\nm2\n")
        self.assertTrue((state / "zamanlanmis" / "1.json").exists())
        self.assertFalse((state / "backups").exists())
        self.assertFalse((state / "tetikleyici.log").exists())
        self.assertTrue((self.snaps / "latest" / "tetikleyici.yaml").exists())
        self.assertTrue((self.snaps / "latest" / "messages.jsonl").exists())

    def test_saklama_politikasi_eskileri_siler(self):
        for ad in ("20200101-000000", "20200102-000000", "20200103-000000"):
            (self.snaps / ad).mkdir(parents=True)
        (self.snaps / "elle-tutulan").mkdir()
        self.assertEqual(self.calistir(keep="2").returncode, 0)
        kalan = sorted(p.name for p in self.snaps.iterdir() if p.is_dir() and not p.is_symlink())
        self.assertIn("20200103-000000", kalan)
        self.assertNotIn("20200101-000000", kalan)
        self.assertNotIn("20200102-000000", kalan)
        self.assertIn("elle-tutulan", kalan)  # zaman damgası olmayan klasöre dokunulmaz
        self.assertEqual(len([k for k in kalan if k[0].isdigit()]), 2)

    def test_bozuk_veritabaninda_hata_ve_yarim_yedek_kalmaz(self):
        self.db.close()
        for ek in ("", "-wal", "-shm"):
            p = self.session / f"session.db{ek}"
            if p.exists():
                p.unlink()
        (self.session / "session.db").write_bytes(b"bu bir sqlite dosyasi degil" * 100)
        self.db = sqlite3.connect(":memory:")
        r = self.calistir()
        self.assertNotEqual(r.returncode, 0)
        if self.snaps.exists():
            self.assertEqual([p for p in self.snaps.iterdir()], [])

    def test_gecersiz_keep_reddedilir(self):
        r = self.calistir(keep="0")
        self.assertNotEqual(r.returncode, 0)
        self.assertIn("NERVE_SNAPSHOT_KEEP", r.stderr)


def test_nerve_snapshot_flow(tmp_path):
    """Fonksiyonel pytest snapshot akışı ve bütünlük kontrolü."""
    session_dir = tmp_path / "session_dir"
    session_dir.mkdir()
    db_file = session_dir / "session.db"

    conn = sqlite3.connect(str(db_file))
    conn.execute("PRAGMA journal_mode = WAL")
    conn.execute("CREATE TABLE test_sessions (id TEXT PRIMARY KEY, data TEXT)")
    conn.execute("INSERT INTO test_sessions VALUES ('s1', 'mefamex-secret')")
    conn.commit()
    conn.close()

    (session_dir / "messages.jsonl").write_text('{"id": "msg1", "text": "selam"}\n')

    snapshot_dir = tmp_path / "snapshots"
    script_path = Path(__file__).resolve().parent.parent / "scripts" / "nerve-snapshot.sh"

    env = os.environ.copy()
    env["NERVE_SESSION_DIR"] = str(session_dir)
    env["NERVE_SNAPSHOT_DIR"] = str(snapshot_dir)
    env["NERVE_SNAPSHOT_KEEP"] = "2"
    env["NERVE_PROJECT_ROOT"] = str(tmp_path)

    res1 = subprocess.run(["/bin/bash", str(script_path)], env=env, capture_output=True, text=True)
    assert res1.returncode == 0, f"Snapshot 1 hata: {res1.stderr}"

    latest_link = snapshot_dir / "latest"
    assert latest_link.exists()

    snap_db = latest_link / "session.db"
    assert snap_db.exists()

    conn_snap = sqlite3.connect(str(snap_db))
    cur = conn_snap.cursor()
    cur.execute("PRAGMA integrity_check")
    assert cur.fetchone()[0] == "ok"
    cur.execute("SELECT data FROM test_sessions WHERE id='s1'")
    assert cur.fetchone()[0] == "mefamex-secret"
    conn_snap.close()


if __name__ == "__main__":
    unittest.main()

import os
import sqlite3
import subprocess
from pathlib import Path


def test_nerve_snapshot_flow(tmp_path):
    # 1. Hazırlık: Geçici WhatsApp oturum dizini ve veritabanı
    session_dir = tmp_path / "session_dir"
    session_dir.mkdir()
    db_file = session_dir / "session.db"

    conn = sqlite3.connect(str(db_file))
    conn.execute("PRAGMA journal_mode = WAL")
    conn.execute("CREATE TABLE test_sessions (id TEXT PRIMARY KEY, data TEXT)")
    conn.execute("INSERT INTO test_sessions VALUES ('s1', 'mefamex-secret')")
    conn.commit()
    conn.close()

    # Sahte mesaj geçmişi
    (session_dir / "messages.jsonl").write_text('{"id": "msg1", "text": "selam"}\n')

    # Snapshot hedef dizini
    snapshot_dir = tmp_path / "snapshots"
    script_path = Path(__file__).resolve().parent.parent / "scripts" / "nerve-snapshot.sh"

    env = os.environ.copy()
    env["NERVE_SESSION_DIR"] = str(session_dir)
    env["NERVE_SNAPSHOT_DIR"] = str(snapshot_dir)
    env["NERVE_SNAPSHOT_KEEP"] = "2"
    env["NERVE_PROJECT_ROOT"] = str(tmp_path)

    # 2. Çalıştırma: 1. Snapshot
    res1 = subprocess.run(["/bin/bash", str(script_path)], env=env, capture_output=True, text=True)
    assert res1.returncode == 0, f"Snapshot 1 hata: {res1.stderr}"

    latest_link = snapshot_dir / "latest"
    assert latest_link.exists()

    # Bütünlük kontrolü
    snap_db = latest_link / "session.db"
    assert snap_db.exists()

    conn_snap = sqlite3.connect(str(snap_db))
    cur = conn_snap.cursor()
    cur.execute("PRAGMA integrity_check")
    assert cur.fetchone()[0] == "ok"
    cur.execute("SELECT data FROM test_sessions WHERE id='s1'")
    assert cur.fetchone()[0] == "mefamex-secret"
    conn_snap.close()

    # 3. Çalıştırma: 2. ve 3. Snapshot (Retention / Saklama Testi)
    subprocess.run(["/bin/bash", str(script_path)], env=env, check=True)
    subprocess.run(["/bin/bash", str(script_path)], env=env, check=True)

    # En fazla KEEP=2 snapshot klasörü kalmalı (+ latest symlink)
    dirs = [d for d in snapshot_dir.iterdir() if d.is_dir() and not d.is_symlink()]
    assert len(dirs) == 2, f"Beklenen 2 snapshot klasörü, bulunan: {len(dirs)}"

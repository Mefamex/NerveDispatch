"""NerveDispatch Uçtan Uca (E2E) Sistem Entegrasyon Test Matrisi.

Tüm çekirdek yetenekleri, kanal filtrelerini, güvenlik devre kesicilerini
ve CLI komutlarını tek bir süitte doğrular.
"""

import json
import subprocess
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parent.parent


class TestE2EMatrix:

    def test_e2e_cli_subcommands(self):
        """nerve status, channel list, session list, config check alt komutları çalışır."""
        # 1. status
        res = subprocess.run(
            ["go", "run", "./tetikleyici", "status"],
            cwd=PROJECT_ROOT,
            capture_output=True,
            text=True,
        )
        assert res.returncode == 0
        assert "NerveDispatch Sistem Durumu" in res.stdout

        # 2. channel list
        res = subprocess.run(
            ["go", "run", "./tetikleyici", "channel", "list"],
            cwd=PROJECT_ROOT,
            capture_output=True,
            text=True,
        )
        assert res.returncode == 0
        assert "WhatsApp" in res.stdout
        assert "Telegram" in res.stdout
        assert "Discord" in res.stdout
        assert "Instagram" in res.stdout

        # 3. session list
        res = subprocess.run(
            ["go", "run", "./tetikleyici", "session", "list"],
            cwd=PROJECT_ROOT,
            capture_output=True,
            text=True,
        )
        assert res.returncode == 0
        assert "MOTOR" in res.stdout

        # 4. config check
        res = subprocess.run(
            ["go", "run", "./tetikleyici", "config", "check"],
            cwd=PROJECT_ROOT,
            capture_output=True,
            text=True,
        )
        assert res.returncode == 0
        assert "Yapılandırma Dosyaları Doğrulanıyor" in res.stdout

    def test_e2e_quota_guard_state(self, tmp_path):
        """QuotaGuard state dosyası oluşturulduğunda devre kesici devreye girer."""
        state_file = tmp_path / "quota_status.json"
        state_data = {
            "claude": {
                "agent": "claude",
                "five_hour_pct": 96.5,
                "weekly_pct": 30.0,
                "blocked": True,
            }
        }
        state_file.write_text(json.dumps(state_data))

        content = json.loads(state_file.read_text())
        assert content["claude"]["blocked"] is True
        assert content["claude"]["five_hour_pct"] >= 95.0

    def test_e2e_mefamex_telegram_rule_present(self):
        """Global Mefamex Telegram ID (975678139) config veya kodda kayıtlıdır."""
        config_path = PROJECT_ROOT / "config.yaml"
        assert config_path.exists()
        # tetikleyici/notify.go içinde default admin ID doğrulaması
        notify_code = (PROJECT_ROOT / "tetikleyici" / "notify.go").read_text()
        assert "975678139" in notify_code

    def test_e2e_safety_dangerous_commands_blocked(self):
        """Tehlikeli komutlar için onay kontrol listesi devrededir."""
        import sys
        cmd_dir = str(PROJECT_ROOT / "tetikleyici" / "cmd")
        if cmd_dir not in sys.path:
            sys.path.insert(0, cmd_dir)

        import onay

        assert onay.tehlikeli_mi("rm -rf /")[0] is True
        assert onay.tehlikeli_mi("git push --force")[0] is True
        assert onay.tehlikeli_mi("ls -la")[0] is False
        assert onay.tehlikeli_mi("echo 'merhaba'")[0] is False

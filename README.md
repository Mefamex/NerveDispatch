# Nerve Dispatch (NerveDispatch)

> Çok kanallı (WhatsApp, Telegram, Discord, Instagram) otonom yapay zeka ajan motoru ve cihaz ağ geçidi - Linux terminal tüneli, yerel multimodal boru hatları (Whisper/Qwen) ve CLI kodlama ajanları entegrasyonu.

| | |
| --- | --- |
| **Author** | [Mehmet Akif Akkoç (Mefamex)](https://github.com/Mefamex) |
| **Website** | [https://www.mefamex.com/](https://www.mefamex.com/) |
| **License** | MIT |
| **Status** | Active Development (NerveDispatch Monorepo Migration) |

---

## 📖 Proje Hikayesi ve Evrimi

CHIMERA'nın 2024 yılından bu yana geçirdiği tüm evrim, mimari kırılma anları, geçmiş projelerle birleşme serüveni ve projenin **Nerve Dispatch** kimliğine kavuşma hikayesi için ana dokümanı inceleyin:

👉 **[NERVE-DISPATCH-THE-STORY.md](./NERVE-DISPATCH-THE-STORY.md)**

---

## 🏗️ Mimari ve Birleşme Planı

İki ayrı depo olarak başlayan WhatsApp Köprüsü (`whatsapp-bridge-mcp`) ve Dağıtıcı Motorunun (`whatsapp-tetikleyici`) tek çatı altında toplanarak evrensel bir ağ geçidine dönüştürülmesi planı:

👉 **[YENI-MIMARI-VE-BIRLESME-PLANI.md](./YENI-MIMARI-VE-BIRLESME-PLANI.md)**

Açık işler ve sprint yol haritası:
👉 **[TODO.md](./TODO.md)**

---

## ⚙️ Modüler Kurulum (`install/`) ve Yükseltme (`upgrade/`) Mimarisi

Nerve Dispatch; çoklu iletişim kanallarını, yerel ses/görsel boru hatlarını ve CLI kodlama ajanlarını (`agy`, `claude`, `opencode`) tek bir şişkin `Makefile` içine yığmaz. Kök `Makefile` yalnızca ince bir yönlendiricidir:

- **`make install`** -> `./install/install.sh`: Sistem araçları (`00`), Go çekirdeği (`10`), Python `.venv` (`20`), iletişim kanalları & MCP (`30`), ses STT/TTS modülleri (`40`), görüntü/OCR modülleri (`50`), AI ajan MCP/Skill ayarları (`60`) ve `systemd` servislerini (`70`) modüler olarak kurar.
- **`make upgrade`** -> `./upgrade/upgrade.sh`: Canlı verileri ve `session.db` oturumunu yedekler (`00`), çekirdeği atomik günceller (`10`), kanal protokollerini ve MCP araçlarını senkronize eder (`20`), medya modellerini günceller (`30`), kişisel verileri silmeden YAML şemasını göç ettirir (`40`) ve 627+ testi doğrulayıp (hata varsa otomatik rollback) servisleri yeniden başlatır (`50`).

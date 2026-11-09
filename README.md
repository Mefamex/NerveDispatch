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

👉 **[NERVE-DISPATCH-THE-STORY.md](./docs/NERVE-DISPATCH-THE-STORY.md)**

---

## 🏗️ Mimari ve Birleşme Planı

İki ayrı depo olarak başlayan WhatsApp Köprüsü (`whatsapp-bridge-mcp`) ve Dağıtıcı Motorunun (`whatsapp-tetikleyici`) tek çatı altında toplanarak evrensel bir ağ geçidine dönüştürülmesi planı:

👉 **[YENI-MIMARI-VE-BIRLESME-PLANI.md](./docs/YENI-MIMARI-VE-BIRLESME-PLANI.md)**

## 🗺️ Belge Haritası (Repository Doc Map)

| Belge | Amacı ve Sorumluluğu |
| :--- | :--- |
| [`README.md`](./README.md) | Proje vitrini, mimari özet ve modüler kurulum yapısı |
| [`TODO.md`](./docs/TODO.md) | Aktif sprint ve açık yapılacak işler listesi |
| [`TODO.DONE.md`](./docs/TODO.DONE.md) | Doğrulanmış iş arşivi (ters kronolojik, test kanıtı zorunlu) |
| [`NOTES.md`](./docs/NOTES.md) | Dahili teknik bellek, mühendislik kararları ve elenen yöntemler |
| [`ASK.md`](./docs/ASK.md) | Açık sorular, bloklayıcılar ve ajanın otonom inisiyatif kayıtları |
| [`GIT_STANDARDS.md`](./docs/GIT_STANDARDS.md) | Git & Commit standartları (Conventional Commits v1.0.0, gatekeeping) |
| [`pre_commit_checklist.md`](./docs/pre_commit_checklist.md) | Commit öncesi hızlı kontrol listesi (< 15 sn, secret/PII kapısı) |
| [`pre_push_checklist.md`](./docs/pre_push_checklist.md) | Push öncesi tam doğrulama listesi ve kullanıcı onay kapısı |
| [`YENI-MIMARI-VE-BIRLESME-PLANI.md`](./docs/YENI-MIMARI-VE-BIRLESME-PLANI.md) | WhatsApp köprüsü ve tetikleyicinin monorepo birleşme mimarisi |
| [`NERVE-DISPATCH-THE-STORY.md`](./docs/NERVE-DISPATCH-THE-STORY.md) | CHIMERA'dan NerveDispatch'e tarihsel evrim ve vizyon belgesi |

---

## ⚙️ Modüler Kurulum (`install/`) ve Yükseltme (`upgrade/`) Mimarisi

Nerve Dispatch; çoklu iletişim kanallarını, yerel ses/görsel boru hatlarını ve CLI kodlama ajanlarını (`agy`, `claude`, `opencode`) tek bir şişkin `Makefile` içine yığmaz. Kök `Makefile` yalnızca ince bir yönlendiricidir:

- **`make install`** -> `./install/install.sh`: Sistem araçları (`00`), Go çekirdeği (`10`), Python `.venv` (`20`), iletişim kanalları & MCP (`30`), ses STT/TTS modülleri (`40`), görüntü/OCR modülleri (`50`), AI ajan MCP/Skill ayarları (`60`) ve `systemd` servislerini (`70`) modüler olarak kurar.
- **`make upgrade`** -> `./upgrade/upgrade.sh`: Canlı verileri ve `session.db` oturumunu yedekler (`00`), çekirdeği atomik günceller (`10`), kanal protokollerini ve MCP araçlarını senkronize eder (`20`), medya modellerini günceller (`30`), kişisel verileri silmeden YAML şemasını göç ettirir (`40`) ve 627+ testi doğrulayıp (hata varsa otomatik rollback) servisleri yeniden başlatır (`50`).

---

## 💻 Birleşik Yönetim Arayüzü (`nerve` CLI)

NerveDispatch, tüm servisleri ve kanalları tek bir ikili (`nerve` veya `tetikleyici`) üzerinden yönetir:

```bash
# Sistem ve kanal durumunu, kota kullanımını ve aktif servisleri denetler
nerve status

# Eşzamanlı WAL modunda açık session.db ve durum dosyalarının atomik yedeğini alır
nerve snapshot

# Aktif kanalları (WhatsApp, Telegram, Discord, Instagram) listeler
nerve channel list

# Telefon numarası ile 8 haneli Whatsmeow eşleştirme kodu üretir (QR alternatifi)
nerve pair 90545xxxxxxx

# WhatsApp köprüsü eşleştirme QR kodunu terminalde ASCII olarak basar
nerve qr

# Sistem loglarını canlı takip eder
nerve logs -f

# Yapılandırma ve kural sözdizimini doğrular
nerve config check

# Ajan oturumlarını listeler veya sıfırlar
nerve session list
nerve session reset
```

---

## 🌐 Gömülü Web Paneli (Dashboard)

Ekstra web sunucusu (Nginx, Node vb.) gerektirmeden, doğrudan Go ikilisine gömülü (`//go:embed`) hafif karanlık tema yönetim paneli:

- **Erişim Adresi:** `http://127.0.0.1:31415/dashboard/`
- **Özellikler:**
  - Canlı bağlantı, bellek ve uptime telemetrisi
  - Anlık QR kodu ve 8 haneli eşleştirme desteği
  - Tüm kanalların (WhatsApp, Telegram, Discord, Instagram) canlı durum kartları
  - Son gelen ve giden mesaj akışı
  - Sistem kaynak tüketimi ve sistem bekçisi uyarıları

---

## 🛡️ Güvenlik, Devre Kesiciler ve Yedekleme

1. **QuotaGuard (%95 Kota Sigortası):**
   - Kodlama ajanlarının (`claude`, `agy`, `opencode`) 5 saatlik ve haftalık limitleri sürekli izlenir.
   - Kota $\ge \%95$ eşiğine ulaştığında komut yürütme otomatik durdurulur ve bir sonraki kota yenilenme periyoduna kadar güvenli uyku moduna geçilir.

2. **Zero-Data-Loss SQLite WAL Snapshot:**
   - WhatsApp mesaj geçmişi ve anahtar çiftleri (`session.db`) WAL modunda çalışırken kilitlenmeden, Python `sqlite3.backup()` API'si ile `%100` atomik yedeklenir.
   - Son 14 snapshot otomatik saklanır, `PRAGMA integrity_check` doğrulaması yapılır.

3. **İnteraktif İzin ve Onay Sistemi (`/onay`):**
   - Tehlikeli kabuk komutları (`rm -rf`, `git push --force`, `systemctl stop`, `reboot` vb.) doğrudan çalıştırılmaz.
   - WhatsApp üzerinden 4 seçenekli interaktif anket (poll) açılarak admin onayı istenir (tek seferlik, reddet, 1 saatlik, 1 haftalık otomatik onay).

4. **Cross-Channel Failover:**
   - WhatsApp bağlantısı koptuğunda kritik sistem alarmları otomatik olarak Telegram Bot API (`975678139`) üzerinden iletilir.

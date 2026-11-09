# Nerve Dispatch (NerveDispatch) - Çok Kanallı Mimari ve Birleşme Planı (Master Blueprint)

> **Proje Adı:** `Nerve Dispatch` (`NerveDispatch` - CLI: `nerve` / `tetikleyici`)  
> **Hedef Konum:** `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/`  
> **Tarihsel & Felsefi Arka Plan:** [`/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md`](file:///mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md)  
> **Birleştirilecek Kaynak Repolar:**
> 1. `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-bridge-mcp`
> 2. `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-tetikleyici`  
> **Durum:** Onaylandı - Sonraki Sohbette Doğrudan Uygulamaya Geçilecek (FAZ 1 → FAZ 2 → FAZ 3)  
> **Son Güncelleme:** 2026-11-09

---

## 🤖 YENİ SOHBETTE BAŞLAYACAK AJAN İÇİN DEVİR-TESLİM TALİMATI (AGENT HANDOFF)

Eğer yeni bir oturumda kullanıcı **"Nerve Dispatch birleşmesine başlayalım"** veya **"plana göre devam edelim"** derse:
1. **Önce şu iki belgeyi oku:**
   - Bu belge (`YENI-MIMARI-VE-BIRLESME-PLANI.md`) - Teknik mimari ve adım adım geçiş/geliştirme planı.
   - [`/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md`](file:///mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md) - Projenin v1 ve v2'den bugüne evrimi ve felsefesi.
2. **Stratejik Sırayı Bozmamaya Dikkat Et:**
   - **Eski repolarda (`whatsapp-bridge-mcp` ve `whatsapp-tetikleyici`) artık yeni özellik kodlanmayacaktır.**
   - **İlk iş (FAZ 1):** İki kaynak repoyu `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` içine birleştirmek, `install/` ve `upgrade/` modüler orkestrasyon klasörlerini ve ince `Makefile` yapısını kurmak, canlı `tetikleyici.yaml`, `.tetikleyici/` durum dosyalarını taşımak, `~/.local/share/whatsapp-bridge/session.db` oturumunu (QR bozmadan) korumak, `make verify` ile 627+ testi yeşil geçirmek ve `systemd` + MCP ayarlarını `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` yoluna bağlamaktır.
   - **İkinci iş (FAZ 2):** `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` üzerinde **önce WhatsApp hattını ve tüm çekirdek yetenekleri** (`Geliştirme 2.1` Sabit Session ID & Ajan Sürücülerinden başlayarak `2.11`'e kadar) sırasıyla tamamlamaktır.
   - **Üçüncü iş (FAZ 3):** Çekirdek tamamlandıktan sonra **Telegram (`975678139`), Discord ve Instagram** kanallarını eklemektir.
3. **Kırmızı Çizgileri Unutma:**
   - `sudo` komutlarını kendi başına çalıştırma; kullanıcıya tek satırda göster.
   - `~/.local/share/whatsapp-bridge/session.db` ve `messages.jsonl` dosyalarını asla silme veya sıfırlama (canlı WhatsApp eşleşmesi oradadır).
   - Python MCP klasörünün adı modül gölgelemesini (`from mcp.types import ...`) önlemek için mutlaka `mcp_server/` olarak kalmalıdır.

---

## 1. Nerve Dispatch Vizyonu: Çok Başlı, Çok Kanallı, Tek Zihinli Sistem

1. **Çok Kanallı Duyu Organları (Omnichannel Gateway):**
   - **WhatsApp** *(Birincil & Aktif Kanal - Whatsmeow Go soketi)*
   - **Telegram** *(İkincil Kanal - Faz 3.2, Mefamex Telegram ID: `975678139`)*
   - **Discord** *(Üçüncül Kanal - Faz 3.3, Gateway WebSocket & Slash Komutları)*
   - **Instagram** *(Dördüncül Kanal - Faz 3.4, DM Otomasyonu)*
2. **Tek Ortak Beyin ve Olay Dağıtıcı (Unified Dispatcher & Persistent Memory):**
   - Mesaj hangi kanaldan gelirse gelsin (WhatsApp, Telegram, Discord, Instagram), tek bir **Ortak Olay Zarfına (`ChimeraEvent`)** dönüştürülür.
   - **Kanal-Arası Ortak Ajan Hafızası:** Kullanıcı dışarıdayken WhatsApp'tan başlattığı bir kodlama veya sistem görevine, masasına geçtiğinde Telegram veya Discord üzerinden **aynı sabit AI oturumu (`opencode`, `agy`, `claude`)** içinde kesintisiz devam edebilir.
3. **Kanal-Arası Yedeklilik (Cross-Channel Failover):**
   - WhatsApp bağlantısı düştüğünde (`on_bridge_down`), WhatsApp kopuk olduğu için uyarı mesajı otomatik olarak **Telegram** üzerinden Mefamex'e (`975678139`) iletilir.
4. **Hafif, Güvenli ve Her Cihazda Çalışır:**
   - Ağır masaüstü GUI (Electron/Qt) içermez. **Headless Go Daemon (~30 MB RAM) + Birleşik `nerve` CLI + Go `embed` Gömülü Web Dashboard + ChatOps** mimarisiyle hem kişisel masaüstünde (Arch Linux KDE Plasma Wayland) hem de ekransız sunucularda (VPS, Raspberry Pi) çalışır.

---

## 2. Mevcut Sistemin Röntgeni ve Korunacak Canlı Varlıklar

`/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` dizinine geçerken tek bir satır kodun, testin veya canlı oturumun kaybolmaması için mevcut durumun tam envanteri:

### 2.1. Canlı Çalışan Süreçler ve Yetki Ayrımı
| Bileşen | Mevcut Servis | Yeni `NerveDispatch` Karşılığı | Yetki & Ortam Gerekçesi |
| :--- | :--- | :--- | :--- |
| **Köprü (Bridge)** | `/etc/systemd/system/whatsapp-bridge.service` | `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/bin/whatsapp-bridge` | Sistem servisi (`ProtectSystem=strict`, `ProtectHome=read-only`). Yalnızca protokol soketlerini ve `~/.local/share/whatsapp-bridge/` dizinini yönetir. |
| **Dağıtıcı (Dispatcher)** | `~/.config/systemd/user/whatsapp-tetikleyici.service` | `/home/mfmx/bin/tetikleyici --config /mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/tetikleyici.yaml` | Kullanıcı servisi (`systemd --user`). Wayland masaüstü (`spectacle`, `wl-copy`), kamera ve yerel AI CLI araçlarına (`agy`, `claude`, `opencode`) erişmek için kullanıcı oturumunda çalışmak zorundadır. |

### 2.2. Sıfır Kesintiyle Korunacak Canlı Veriler
1. **WhatsApp Oturum Veritabanı (`~/.local/share/whatsapp-bridge/`):**
   - `session.db` (~13 MB SQLite WAL) ve `messages.jsonl` (~5.5 MB).
   - **Kritik:** `NerveDispatch` geçişinde bu dizin yolu değiştirilmeyecek; böylece **QR kodu tekrar okutmaya gerek kalmadan** WhatsApp oturumu saniyesinde çalışmaya devam edecektir.
2. **Canlı Kural ve Kimlik Dosyaları (`/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-tetikleyici/`):**
   - `tetikleyici.yaml` ve `pasif-kurallar.yaml` (gerçek telefon numaraları, LID'ler ve yetki bayrakları).
3. **Çalışma Zamanı Durum Dizini (`.tetikleyici/`):**
   - `seen.txt` (işlenmiş mesaj ID'leri), `bildirim.txt` (bildirim throttle damgaları), `bekleyen_cevaplar.txt*` (`/reply` kuyruğu), `zamanlanmis/` (aktif zamanlanmış görevler), `tetikleyici.log`.
   - *(Mevcut dizindeki hatalı iç içe `.tetikleyici/.tetikleyici` klasörü geçişte temizlenecektir).*

### 2.3. Geçişte Onarılacak Ekosistem Kırıkları
Eski `~/Desktop/whatsapp-bridge-project` silindiğinde kırık kalan ve `NerveDispatch` geçişinde (FAZ 1) kalıcı olarak `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yoluna bağlanacak noktalar:
- `bridge/internal/appcfg/appcfg.go` içindeki `ProjectSubDir` ve `SidecarHomeSubDir` sabitleri.
- `~/.gemini/config/mcp_config.json` (Antigravity `PYTHONPATH`).
- `~/.config/opencode/opencode.json` (OpenCode `cwd` + `wp-sahip`/`wp-misafir` araç listesine 5 yeni MCP aracının eklenmesi).
- `~/.claude.json` (Claude Code `PYTHONPATH`).
- `~/.gemini/config/skills/whatsapp-bridge/SKILL.md`, `~/.config/opencode/skills/whatsapp-bridge/SKILL.md` ve `~/.config/opencode/INDEX.md`.

---

## 3. Nerve Dispatch - Çok Kanallı Hedef Mimari Şeması

```text
  ┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐
  │  WhatsApp (Aktif)  │  │ Telegram (Faz 3.2) │  │  Discord (Faz 3.3) │  │ Instagram (Faz 3.4)│
  │  Whatsmeow Soketi  │  │  Bot API / MTProto │  │  Gateway WebSocket │  │   DM / Graph API   │
  └─────────┬──────────┘  └─────────┬──────────┘  └─────────┬──────────┘  └─────────┬──────────┘
            │                       │                       │                       │
            └───────────────────────┴───────────┬───────────┴───────────────────────┘
                                                │
                                                ▼
┌──────────────────────────────────────────────────────────────────────────────────────────────┐
│ 1. KATMAN: NERVEDISPATCH KÖPRÜ MOTORU - bridge/ (Go Headless Daemon :31415)                  │
│  - Kanal Adaptörleri (ChannelAdapter: WhatsApp [Aktif], Telegram, Discord, Instagram)        │
│  - Ortak Olay Zarfı Üreticisi (Unified ChimeraEvent: channel, id, from, type, text, media)   │
│  - Kalıcı Oturum & Mesaj Deposu (session.db + messages.jsonl, ~30 MB RAM)                    │
│  - Gömülü Web Sunucusu (go:embed Web Dashboard: http://127.0.0.1:31415/dashboard)            │
│  - Birleşik HTTP REST API (/send, /reaction, /poll...) & Gerçek Zamanlı SSE (/events)        │
│  - FastMCP STDIO Sunucusu (mcp_server/ - AI Ajanları için Bildirim ve İletişim Araçları)     │
└───────────────────────────────┬──────────────────────────────────────────────┬───────────────┘
                                │ Gerçek Zamanlı SSE Olay Akışı (/events)      │
                                ▼                                              ▼
┌───────────────────────────────────────────────────────────────┐ ┌────────────────────────────┐
│ 2. KATMAN: NERVEDISPATCH DAĞITICI - tetikleyici/ (systemd)    │ │  Yönetim & Gözlem Katmanı  │
│  - Çok Kanallı Kimlik & RBAC (WP LID/Tel, TG ID, DC ID)       │ │ - Birleşik `nerve` CLI     │
│  - Catch-up Yaş Filtresi & Kanal-Arası Yedeklilik (Failover)  │ │   (status, qr, logs, pair) │
│  - Tünel Durum Makinesi (/cmdhub Bash & /session Tüneli)      │ │ - Yerel Web Dashboard      │
│  - Ajan Adaptör Katmanı (Sabit Session ID & Eşzamanlı Kilit)  │ │   (Canlı QR, Kanal Durumu) │
└───────┬───────────────────────┬───────────────────────┬───────┘ └────────────────────────────┘
        │                       │                       │
        ▼ [Sistem Komutları]    ▼ [Multimodal Borular]  ▼ [Sabit Oturumlu AI Ajan Sürücüleri]
 ┌──────────────┐        ┌─────────────────────────┐   ┌───────────────────────────────────────┐
 │ Doğrudan CLI │        │ Yerel Çıkarım Hatları   │   │ Ortak AgentDriver Katmanı (drivers/)  │
 │ /cmd, /camera│        │ - Whisper (Ses -> Metin)│   │ - OpenCode Driver (--session)         │
 │ /ekran, /clip│        │ - Piper TTS (Metin->Ses)│   │ - AGY Driver (--conversation)         │
 │ /get, /root  │        │ - Ollama V-LLM (Görsel) │   │ - Claude Driver (--resume / --session)│
 │ /apps, /kill │        │ (Otomatik VRAM Tahliye) │   │ (.tetikleyici/oturumlar.json)         │
 └──────┬───────┘        └────────────┬────────────┘   └───────────────────┬───────────────────┘
        │                             │                                    │
        └─────────────────────────────┼────────────────────────────────────┘
                                      │
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────────────────────────┐
 │ 3. KATMAN: KANALA ÖZEL FORMATLAMA VE TESLİMAT                                               │
 │  - Hedef kanala göre biçimlendirme (WhatsApp biçimi vs Telegram/Discord Markdown)           │
 │  - Canlı İlerleme Reaksiyonu (👀 İşleniyor -> ✅ Başarılı / ❌ Hatalı)                      │
 │  - POST /send (channel, to, text/media) -> İlgili Kanal Üzerinden Kullanıcıya Teslimat      │
 └─────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Çok Kanallı Veri ve Yapılandırma Standardı (Geriye Dönük %100 Uyumlu)

**Altın Kural:** Mevcut 627+ testi ve canlı `tetikleyici.yaml` dosyasını kırmamak için **tüm yeni alanlar opsiyoneldir ve varsayılan kanal `"whatsapp"` olarak kabul edilir.**

### 4.1. Ortak Olay Zarfı (`ChimeraEvent` / `message`)
Köprüden SSE `/events` akışına basılan JSON yapısı:
```json
{
  "channel": "whatsapp",
  "id": "3EB0A1B2C3D4",
  "from": "12345678901234@lid",
  "takma_ad": "905551234567@s.whatsapp.net",
  "sender": "12345678901234@lid",
  "push_name": "Mefamex",
  "type": "text",
  "text": "/durum",
  "timestamp": 1762638000,
  "is_from_me": false
}
```
- `channel` alanı boş gelirse veya belirtilmezse sistem otomatik olarak `"whatsapp"` kabul eder (mevcut tüm testler sıfır değişiklikle geçer).

### 4.2. Çok Kanallı Kullanıcı ve Kural Yapısı (`tetikleyici.yaml`)
Mevcut `number:` ve `wp_number:` alanları aynen çalışmaya devam ederken, Faz 3 için kanal kimlikleri eklenebilir:
```yaml
users:
  users:
    - name: ["Mefamex"]
      number: ["905551234567", "12345678901234@lid"]   # WhatsApp (mevcut, aynen korunur)
      tg_id: ["975678139"]                             # Telegram (Faz 3'te aktifleşir)
      dc_id: ["discord_user_id"]                       # Discord (Faz 3'te aktifleşir)
      ig_id: ["instagram_id"]                          # Instagram (Faz 3'te aktifleşir)
      allow_slash: true
      allow_cmd: true
      allow_camera: true
      allow_ekran: true
      allow_durum: true
      allow_zamanla: true
      allow_reply: true
```

---

## 5. `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` - Modüler Monorepo Dizin Hiyerarşisi

Mevcut testlerin (`realkonfig_test.go`, `shipped_config_test.go`, `tests/test_*.py`) beklediği göreli dosya yollarını (`../tetikleyici.yaml`, `../config.readme.md`, `Path(__file__).parent.parent`) hiç bozmayan, aynı zamanda çok kanallı genişlemeye hazır dizin yapısı:

```text
/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/
├── NERVE-DISPATCH-THE-STORY.md            # [MEVCUT] Chimera v1 • v2 • v3 Mimari Evrim Hikayesi
├── YENI-MIMARI-VE-BIRLESME-PLANI.md # [MEVCUT] Bu Master Blueprint belgesi
│
├── bridge/                         # Go Modülü 1: Çok Kanallı Köprü & HTTP/SSE Sunucusu (:31415)
│   ├── main.go                     # Köprü giriş noktası
│   ├── dispatch.go                 # Mesaj sınıflandırıcı ve olay yönlendirici
│   ├── go.mod & go.sum             # module whatsapp-bridge (go.work ile bağlı)
│   └── internal/
│       ├── appcfg/                 # Süreç yolları ve yapılandırma
│       ├── channels/               # [YENİ - Faz 3] Kanal Adaptör Arayüzü (telegram, discord, instagram)
│       ├── eventbus/               # SSE olay yayın otobüsü
│       ├── httpapi/                # REST ve SSE HTTP sunucusu (:31415)
│       ├── msgstore/               # Bellek içi + messages.jsonl kalıcı mesaj deposu
│       ├── voip/                   # Sesli arama (call) yöneticisi
│       ├── wabridge/               # WhatsApp (whatsmeow) kanal çekirdeği
│       ├── webhook/                # Webhook dağıtıcısı
│       └── webui/                  # [YENİ - Faz 2] go:embed ile gömülü yerel Web Dashboard (/dashboard)
│
├── tetikleyici/                    # Go Modülü 2: Chimera Olay Dağıtım Motoru & CLI
│   ├── main.go                     # Daemon ve CLI giriş noktası
│   ├── engine.go                   # SSE dinleyici, catch-up ve paralel kural yürütücü
│   ├── action.go                   # Eylem ve ajan çağrı orkestratörü
│   ├── drivers/                    # [YENİ - Faz 2.1] Ajan Sürücüleri (opencode, agy, claude - Sabit Oturum)
│   ├── cli/                        # [YENİ - Faz 2.2] Birleşik CLI (status, qr, pair, logs, session, channel)
│   ├── cmd/                        # Python: Sistem komutları (/cmd, /cmdhub, /camera, /ekran, /durum, /zamanla...)
│   ├── go.mod & go.sum             # module tetikleyici
│   └── *_test.go                   # 30 Go test dosyası (köprüde ayrıca 19)
│
├── mcp_server/                     # Python: FastMCP STDIO Sunucusu (`mcp` kütüphanesiyle çakışmaz)
│   ├── __main__.py                 # python -m mcp_server
│   ├── api/client.py               # Bridge HTTP REST istemcisi
│   ├── config/                     # Yapılandırma yükleyici
│   ├── server/                     # FastMCP sunucu tanımı
│   └── tools/                      # 17 MCP aracı (WhatsApp + Çok Kanallı bildirim araçları)
│
├── gorsel/                         # Python: Görsel, Belge (PDF/DOCX/XLSX), OCR, Video ve Sticker Paketi
├── gorsel_pipeline.py              # Görsel analiz boru hattı (Ollama Qwen2.5-VL + OCR)
├── voice_pipeline.py               # Sesli not transkripsiyon boru hattı (Faster-Whisper GPU/CPU + VRAM tahliye)
│
├── sidecar/                        # Opsiyonel Node.js VoIP / Arama WASM Yardımcı Modülü
├── install/                        # [YENİ] Modüler Kurulum Orkestratörü (make install -> ./install/install.sh)
│   ├── install.sh                  # Ana kurulum yöneticisi (--all, --only <modül>, --headless, --skip-models)
│   ├── lib.sh                      # Ortak yardımcılar (log, yedekleme, güvenli JSON/YAML yama, sudo komut toplayıcı)
│   ├── 00-system-deps.sh           # Sistem paketleri kontrolü (go, gcc, uv, ffmpeg, tesseract, pdftotext, spectacle, wl-clipboard)
│   ├── 10-core-go.sh               # Go Workspace, bridge (CGO sqlite3) ve tetikleyici/nerve derlemesi + ~/bin atomik kurulumu
│   ├── 20-python-env.sh            # Ortak .venv kurulumu ve requirements.txt paketlerinin yüklenmesi
│   ├── 30-channels-mcp.sh          # WhatsApp (ve Faz 3 Telegram/Discord/Instagram) kanal dizinleri, MCP sunucusu ve VoIP sidecar kurulumu
│   ├── 40-voice-modules.sh         # STT (faster-whisper + CUDA 12 venv + large-v3-tr model) ve TTS (Piper TTS + Türkçe ses modeli) kurulumu
│   ├── 50-vision-modules.sh        # Ollama servis kontrolü, görsel/metin modelleri (qwen2.5vl:3b, qwen3:4b) ve OCR (tur+eng) kurulumu
│   ├── 60-agents-config.sh         # AI Ajanları (Antigravity, Claude Code, OpenCode) MCP, 17 araç izni, SKILL.md ve INDEX.md otomatik kurulumu
│   └── 70-configs-services.sh      # Config şablonları (mevcudu ezmeden) ve systemd (--user + sistem) servislerinin kurulumu
│
├── upgrade/                        # [YENİ] Modüler Yükseltme Orkestratörü (make upgrade -> ./upgrade/upgrade.sh)
│   ├── upgrade.sh                  # Ana yükseltme yöneticisi (--all, --only <modül>, --rollback, --with-models)
│   ├── 00-backup-snapshot.sh       # Yükseltme öncesi tam yedek (.tetikleyici/backups/: config, session.db, ajan MCP, ikililer)
│   ├── 10-upgrade-core.sh          # Go ve Python bağımlılıklarını güncelleme ve çalışan süreci bozmayan atomik (cp + mv -f) ikili değişimi
│   ├── 20-upgrade-channels-mcp.sh  # Whatsmeow protokol güncellemesi (405 önleyici), diğer kanal SDK'ları ve yeni MCP araçlarının ajanlara senkronizasyonu
│   ├── 30-upgrade-media-models.sh  # Whisper/Piper ses kütüphaneleri ve Ollama görsel/metin modellerinin güncellenmesi
│   ├── 40-migrate-configs.sh       # Kişisel verileri/numaraları silmeden yeni YAML şema alanlarını canlı config'e güvenli birleştirme
│   └── 50-verify-and-restart.sh    # 627+ test doğrulaması (başarısızsa otomatik rollback), tek süreç kontrolü ve servis restart
│
├── docs/                           # Tasarım belgeleri, arşiv ve teknik analizler (PII temizlenmiş)
├── tests/                          # Birleşik Python Test Süiti (test_mcp_handlers.py + 10 tetikleyici testi)
│
├── bin/                            # Derlenmiş Go ikilileri (gitignore'da)
│   ├── whatsapp-bridge             # Köprü ikilisi (nerve-bridge)
│   └── tetikleyici                 # Dağıtıcı & CLI ikilisi (nerve / tetikleyici)
│
├── .tetikleyici/                   # Çalışma zamanı durum dizini (gitignore'da: seen.txt, oturumlar.json, backups/, log...)
├── go.work                         # Go Workspace dosyası (use ./bridge ./tetikleyici)
├── requirements.txt                # Tekil Python bağımlılıkları (.venv)
├── config.yaml                     # Köprü, kanal ve MCP çalışma parametreleri
├── tetikleyici.example.yaml        # Kural ve yetki şablonu (Git'te takip edilir)
├── tetikleyici.yaml                # Canlı kişisel kural ve yetki dosyası (gitignore'da)
├── pasif-kurallar.example.yaml     # Pasif misafir kuralları şablonu (Git'te takip edilir)
├── pasif-kurallar.yaml             # Canlı pasif kurallar dosyası (gitignore'da)
├── config.readme.md                # Kapsamlı yapılandırma kılavuzu
├── Makefile                        # İnce Görev Yönlendiricisi (make install -> install/, make upgrade -> upgrade/, build, test, verify)
├── README.md                       # Nerve Dispatch proje vitrini ve mimari şeması
├── TODO.md                         # Yalnızca açık işleri barındıran aktif sprint listesi
├── TODO.DONE.md                    # Tamamlanmış ve kanıtlanmış işler arşivi
├── NOTES.md                        # Kalıcı teknik bellek ve sistem kuralları
└── ASK.md                          # Mimari soru ve otonom karar kaydı
```

---

## 6. Birleşik CLI (`nerve` / `tetikleyici`) ve Gömülü Web Dashboard

### 6.1. CLI Komut Seti
Hem `nerve` (`nervedispatch`) hem de `tetikleyici` komutu aynı ikiliyi çalıştırır (geriye dönük uyumluluk tamdır):
- `nerve` (veya `tetikleyici --config ...`): Varsayılan daemon modunda olay dağıtıcıyı başlatır.
- `nerve status`: Köprü, dağıtıcı, bağlı kanallar (WhatsApp, Telegram, Discord, Instagram), aktif AI motoru ve bellek durumunu gösterir.
- `nerve qr` / `nerve pair <numara>`: Terminalde ANSI QR kod çizer veya 8 haneli telefon eşleştirme kodu üretir.
- `nerve logs [-f] [n]`: Tüm kanalların ve dağıtıcının günlüklerini UTF-8 güvenli şekilde akıtır.
- `nerve session list`: `opencode`, `agy` ve `claude` üzerindeki sabit ve açık oturumları listeler.
- `nerve channel list`: Yapılandırılmış kanalların bağlantı durumunu listeler.
- `nerve config check`: YAML dosyalarını, betik yollarını ve yetkileri doğrular.
- `nerve restart`: Servisleri güvenli şekilde yeniden başlatır.

### 6.2. Gömülü Yerel Web Dashboard (`http://127.0.0.1:31415/dashboard`)
- `bridge/internal/webui/` içinde `//go:embed` ile tek binary içine gömülü, **+0 MB RAM** maliyetli kontrol paneli.
- **Bölümler:**
  1. **Kanal Durum Kartları:** WhatsApp (`paired`/`connected`), Telegram, Discord, Instagram bağlantı rozetleri.
  2. **Eşleştirme (Pairing) Merkezi:** Büyük QR kod ve 8 haneli eşleştirme kodu alanı.
  3. **Canlı Olay Akışı (SSE Stream):** Gelen/giden mesajların, kanal etiketleriyle (`[WA]`, `[TG]`, `[DC]`) anlık izlenmesi.
  4. **Ajan & Kural Özeti:** Aktif AI motoru, sabit oturum ID'leri ve yetki matrisi.

---

## 7. Modüler Kurulum (`install/`) ve Yükseltme (`upgrade/`) Orkestrasyon Mimarisi

Nerve Dispatch; çok dilli çekirdeği (Go + Python + opsiyonel Node sidecar), çoklu iletişim kanalları (WhatsApp, Telegram, Discord, Instagram), yerel multimodal yapay zeka boru hatları (Whisper STT, Piper TTS, Ollama Vision, Tesseract OCR) ve harici CLI ajan entegrasyonları (`agy`, `claude`, `opencode`) nedeniyle kapsamlı bir ekosistemdir.

Tüm bu bileşenlerin kurulumunu ve güncellenmesini tek bir `Makefile` içine yığmak yerine, **İnce Makefile + Modüler Orkestratör Klasörleri (`install/` ve `upgrade/`)** mimarisi kullanılır:

### 7.1. İnce `Makefile` Prensibi (Thin Dispatcher)
Kök dizindeki `Makefile` yalnızca giriş noktasıdır; karmaşık kurulum ve yükseltme mantığını alt klasörlere devreder:
```makefile
install: ## Tüm NerveDispatch ekosistemini modüler olarak kurar (ARGS="--only ..." destekler)
	@./install/install.sh $(ARGS)

upgrade: ## Canlı sistemi yedek alarak sıfır kesintiyle yükseltir (ARGS="--only ..." destekler)
	@./upgrade/upgrade.sh $(ARGS)

update: upgrade ## 'make upgrade' için uyumluluk takma adı

reinstall: ## Kurulu ortamı temizleyip sıfırdan kurar
	@./install/install.sh --reinstall $(ARGS)

remove: ## Derleme çıktılarını ve sanal ortamları kaldırır (kişisel verileri ve session.db'yi korur)
	@./install/install.sh --remove $(ARGS)
```
- **Kullanım Örnekleri:**
  - `make install` -> Tüm alt modülleri sırasıyla kontrol eder ve kurar.
  - `make install ARGS="--only agents,channels"` -> Yalnızca ajan MCP ayarlarını ve kanal bağlantılarını kurar/onarır.
  - `make install ARGS="--headless"` -> Sunucu/VPS kurulumu (Wayland/kamera/masaüstü araçlarını ve ağır GPU modellerini atlar).
  - `make upgrade` -> Tam sistem yedeklemesi alır, çekirdeği/kanalları/ajanları günceller, 627+ testi koşturur ve servisleri yeniden başlatır.
  - `make upgrade ARGS="--only channels"` -> Yalnızca `whatsmeow` (WhatsApp 405 Client Outdated çözümü) ve diğer kanal protokollerini günceller.
  - `make upgrade ARGS="--rollback"` -> Son yükseltme öncesi alınan anlık yedekten (`backups/`) ikilileri ve ayarları saniyesinde geri yükler.

### 7.2. `install/` Klasörü Alt Modülleri ve Görevleri
Her betik **idempotent** (tekrar çalıştırıldığında mevcut kişisel veriyi bozmayan) yapıda ve **<= 150 satır** sınırındadır:

1. **`install/lib.sh` (Ortak Çekirdek Kütüphane):**
   - Standart renkli çıktı (`[OK]`, `[SKIP]`, `[WARN]`, `[ERR]`), hata yakalama (`set -euo pipefail`), JSON/YAML güvenli yama yardımcıları.
   - **`sudo` Komut Toplayıcı:** Güvenlik kuralı gereği hiçbir betik doğrudan `sudo` çalıştırmaz! `sudo` gerektiren işlemler (`pacman -S` veya `/etc/systemd/system/` güncellemesi) kuyruğa eklenir ve kurulum sonunda kullanıcıya tek blok hâlinde gösterilir.
2. **`install/00-system-deps.sh` (Sistem ve Donanım Ön Koşulları):**
   - Zorunlu derleme ve çalışma araçları: `go`, `gcc` (`CGO_ENABLED=1` SQLite3 için), `python3`, `uv`, `git`, `curl`, `sqlite3`.
   - Medya ve belge araçları: `ffmpeg`, `ffprobe`, `tesseract` (`tesseract-data-tur`, `tesseract-data-eng`), `pdftotext` (`poppler`).
   - Masaüstü (Wayland) araçları (`--headless` değilse): `spectacle` (`/ekran`), `wl-clipboard` (`/clip`), `/dev/video0` (`/camera`).
3. **`install/10-core-go.sh` (Go Workspace ve Çekirdek İkililer):**
   - `go.work` kontrolü, `bridge` (`bin/whatsapp-bridge`) ve `tetikleyici` (`bin/tetikleyici` + `bin/nerve`) ikililerinin derlenmesi.
   - `~/bin/tetikleyici` ve `~/bin/nerve` konumlarına çalışan süreci bozmayan **atomik kurulum (`cp` + `mv -f`)**.
4. **`install/20-python-env.sh` (Birleşik Python Sanal Ortamı):**
   - `uv venv .venv` oluşturumu ve `requirements.txt` (`mcp`, `httpx`, `pydantic`, `pyyaml`, `pillow`, `requests`, `pytest`, `pytest-asyncio`) kurulumu.
5. **`install/30-channels-mcp.sh` (İletişim Kanalları ve MCP Sunucuları):**
   - **WhatsApp:** `~/.local/share/whatsapp-bridge/` dizin ve izin (`700`) doğrulaması (`session.db` ve `messages.jsonl` korunarak), `mcp_server/` modül doğrulaması, opsiyonel `sidecar/` (Node.js VoIP WASM) kurulumu.
   - **Çoklu Kanal (Faz 3):** Telegram, Discord ve Instagram adaptörlerinin/MCP araçlarının yapılandırma ve bağımlılık doğrulaması.
6. **`install/40-voice-modules.sh` (Ses Modülleri - STT & TTS):**
   - **STT (`voice_pipeline.py`):** `faster-whisper` + CTranslate2 CUDA 12 (`nvidia-cublas-cu12`, `nvidia-cudnn-cu12`) ortamının ve `whisper-large-v3-tr` modelinin kontrolü / indirilmesi.
   - **TTS (`Piper TTS` - Faz 2.11):** `piper` ikilisinin ve Türkçe ses modelinin (`tr_TR-dfki-medium.onnx`) kurulumu.
7. **`install/50-vision-modules.sh` (Görüntü, Belge ve OCR Modülleri):**
   - `ollama` servis kontrolü (`127.0.0.1:11434`); birincil görsel model (`qwen2.5vl:3b`) ve metin özetleme modelinin (`qwen3:4b`) kontrolü / `ollama pull` ile indirilmesi (`--skip-models` ile atlanabilir).
   - `gorsel/` alt paketlerinin ve OCR dil paketlerinin (`tur+eng`) duman testi.
8. **`install/60-agents-config.sh` (AI Kodlama Ajanları ve Ekosistem Entegrasyonu):**
   - Sistemdeki `agy` (Antigravity), `claude` (Claude Code) ve `opencode` (OpenCode) CLI araçlarını tespit eder.
   - Üç ajanın yapılandırma dosyalarını (`~/.gemini/config/mcp_config.json`, `~/.config/opencode/opencode.json`, `~/.claude.json`) yedek alarak `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yoluna ve 16+ MCP aracına otomatik bağlar.
   - `~/.gemini/config/skills/whatsapp-bridge/SKILL.md`, `~/.config/opencode/skills/whatsapp-bridge/SKILL.md` ve `~/.config/opencode/INDEX.md` dosyalarını senkronize eder.
   - `.tetikleyici/oturumlar.json` sabit oturum deposunu hazırlar.
9. **`install/70-configs-services.sh` (Yapılandırma Şablonları ve Systemd Servisleri):**
   - `tetikleyici.yaml` ve `pasif-kurallar.yaml` dosyaları mevcut değilse `.example.yaml` şablonlarından üretir (mevcut dosyaya asla dokunmaz).
   - `~/.config/systemd/user/whatsapp-tetikleyici.service` dosyasını `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yoluna göre kurar ve `systemctl --user daemon-reload` çalıştırır.
   - `/etc/systemd/system/whatsapp-bridge.service` birim dosyasını hazırlar ve `sudo` kurulum komutunu raporlar.

### 7.3. `upgrade/` Klasörü Alt Modülleri ve Görevleri
Canlı çalışan bir sistemi veri kaybı veya çift süreç (double-instance) riski olmadan yükseltir:

1. **`upgrade/00-backup-snapshot.sh` (Otomatik Yedek ve Anlık Görüntü):**
   - Her yükseltme öncesinde `.tetikleyici/backups/upgrade-YYYYMMDD-HHMMSS/` dizinine:
     - `tetikleyici.yaml`, `pasif-kurallar.yaml`, `config.yaml`, `.tetikleyici/oturumlar.json`, `.tetikleyici/izinler.json`,
     - Ajan MCP ayarları (`mcp_config.json`, `opencode.json`, `.claude.json`),
     - SQLite `.backup` komutuyla `~/.local/share/whatsapp-bridge/session.db` yedeği,
     - Çalışan mevcut `whatsapp-bridge` ve `tetikleyici` (`nerve`) ikilileri yedeklenir.
2. **`upgrade/10-upgrade-core.sh` (Çekirdek Go & Python Yükseltmesi):**
   - `bridge` ve `tetikleyici` Go paketlerini derler; `~/bin/tetikleyici` ve `~/bin/nerve` ikililerini `cp` + `mv -f` ile atomik olarak günceller (çalışan sürecin `/proc/<pid>/exe` inode'unu bozmadan).
   - `.venv` paketlerini `uv pip install --upgrade -r requirements.txt` ile günceller.
3. **`upgrade/20-upgrade-channels-mcp.sh` (Kanal Protokolleri ve MCP Senkronizasyonu):**
   - WhatsApp Web protokolü `Client outdated (405)` hatasını önlemek için `go.mau.fi/whatsmeow@latest` güncellemesini yapar; Faz 3'teki Telegram/Discord/Instagram SDK'larını günceller.
   - `mcp_server/tools/` altına yeni eklenen araçları otomatik tarayarak Antigravity, OpenCode (`wp-sahip`/`wp-misafir`), Claude ve `SKILL.md` dosyalarına işler.
4. **`upgrade/30-upgrade-media-models.sh` (Ses ve Görüntü Modülleri Yükseltmesi):**
   - `faster-whisper`, Piper TTS ve Ollama görsel/metin modellerini (`--with-models` verilirse) günceller.
5. **`upgrade/40-migrate-configs.sh` (Kayıpsız Yapılandırma Göçü):**
   - `tetikleyici.example.yaml` içine yeni eklenen özellik anahtarlarını (ör. `tepki_emojisi`, `bekci`, `tg_id`, `dc_id`) kullanıcının canlı `tetikleyici.yaml` dosyasıyla karşılaştırır; kişisel telefon/LID bilgilerini silmeden eksik anahtarları raporlar veya güvenli ekler.
   - `.tetikleyici/` içindeki eski/hatalı iç içe klasörleri (`.tetikleyici/.tetikleyici`) ve bayat kilitleri temizler.
6. **`upgrade/50-verify-and-restart.sh` (Test Kapısı, Otomatik Rollback ve Güvenli Restart):**
   - **Zorunlu Test Kapısı:** `make verify` (627+ Go ve Python testi) çalıştırılır. Eğer tek bir test bile başarısız olursa canlı servislere dokunulmaz ve `00-backup-snapshot` üzerinden ikililer otomatik geri yüklenir (**Auto-Rollback**).
   - Testler %100 yeşilse `whatsapp-tetikleyici.service` yeniden başlatılır, `/proc/<pid>/exe` ile tek süreç çalıştığı doğrulanır ve `:31415/health` kontrol edilir.

---

## 8. Korunacak Temel Kurallar ve Sistem Sözleşmeleri

1. **MEVZU (3 Aşamalı Kademeli Yetkilendirme):**
   - Aşama 1: Kod içi varsayılan her zaman `false` / `nil` (açık kapı yok).
   - Aşama 2: `tetikleyici.yaml` başındaki özet blok neyin açık olduğunu gösterir.
   - Aşama 3: Çalışma anında sohbetten verilen komut (`/zamanla aç|kapat`, `/agent`) config'i restart gerektirmeden ezer.
2. **Kapalı Hat Sessizdir:**
   - Kuralı olmayan veya kapalı olan kişilere asla *"Şu an kapalıyım"* gibi otomatik cevaplar dönülmez.
3. **Kod ve Klasör Sınırları:**
   - Üretim dosyaları **≤ 300 satır**, klasörler **≤ 15 dosya** ve **≤ 3 seviye derinlik**.
4. **Sıfır PII ve Secret Sızıntısı:**
   - Gerçek telefon numaraları, LID'ler, Bot Token'ları, `session.db` ve `tetikleyici.yaml` asla Git'e girmez.

---

## 9. Uygulama ve Geçiş Yol Haritası (3 Fazlı Master Roadmap)

> **Uygulama Prensibi:**
> - **İlk Adım (FAZ 1):** Mevcut iki repoyu doğrudan `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` çatısında birleştirmek, `install/` ve `upgrade/` modüler orkestrasyon klasörlerini kurmak ve tüm canlı sistemi sıfır kesintiyle `NerveDispatch` üzerine taşımak.
> - **İkinci Adım (FAZ 2):** **Önce WhatsApp hattını ve tüm çekirdek özellikleri eksiksiz bitirmek.**
> - **Üçüncü Adım (FAZ 3):** Hazır ve kaya gibi sağlam çekirdek üzerine Telegram, Discord ve Instagram kanallarını eklemek.

---

### FAZ 1: Doğrudan `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` Monoreposuna Geçiş, Modüler `install/` & `upgrade/` Kurulumu ve Ekosistem Onarımı [SIRADA #1 🎯]

#### Adım 1.0: Kaynak Repolardaki Commitlenmemiş İşin Kaydı [2026-11-09 doğrulamasıyla eklendi]
- [ ] `whatsapp-bridge-mcp` (28 değişiklik; `sidecar/`, `docs/`, `.agents/` izlenmiyor) ve `whatsapp-tetikleyici` (91 değişiklik) çalışma ağaçları granüler Conventional Commit'lerle kaynak repolarda kaydedilir (PII kapısı: `pre_commit_checklist.md`; push yok).
- [ ] Bu adım bitmeden Adım 1.2 başlamaz: birleşim yalnız commitlenmiş geçmişten yapılır.

> **Yöntem (2026-11-09 kararı):** Adım 1.2'deki "kopyalanır" ifadeleri **dosya kopyası değil**, geçmiş korunarak yol yeniden yazımıdır (`git filter-branch --index-filter` geçici klonda + `git merge --allow-unrelated-histories`). Kaynak repoların tüm geçmişi PII taramasından temiz geçti.
> **Kısıt:** Yeni tetikleyici, eski servis durdurulmadan canlı köprüye bağlanmaz (ayrı `seen.txt` → çift cevap).

#### Adım 1.1: `NerveDispatch` Repo Çatısının Hazırlanması
- [x] `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` dizini, Git deposu (`origin: https://github.com/Mefamex/NerveDispatch.git`), `NERVE-DISPATCH-THE-STORY.md`, `YENI-MIMARI-VE-BIRLESME-PLANI.md`, `TODO.md` ve `README.md` oluşturuldu.
- [ ] Kök dizinde `go.work` (`use ./bridge`, `use ./tetikleyici`) oluşturulur.
- [ ] İki reponun `.gitignore` kuralları (`*.db*`, `.tetikleyici/`, `tetikleyici.yaml`, `pasif-kurallar.yaml`, `.venv/`, `bin/`, `sidecar/node_modules/`, `*.env*`, geçici medya dosyaları) birleştirilir.

#### Adım 1.2: Kodların, Boru Hatlarının ve Testlerin Taşınması
- [ ] `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-bridge-mcp/bridge/` → `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/bridge/` olarak kopyalanır.
- [ ] `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-bridge-mcp/mcp_server/` → `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/mcp_server/` olarak kopyalanır.
- [ ] `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-bridge-mcp/sidecar/` → (`node_modules/` ve `rsrc-dump.txt` hariç) `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/sidecar/` olarak kopyalanır.
- [ ] `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-tetikleyici/tetikleyici/` → `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/tetikleyici/` olarak kopyalanır.
- [ ] `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-tetikleyici/gorsel/`, `gorsel_pipeline.py`, `voice_pipeline.py`, `tetikleyici.sh` → `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` kök dizinine kopyalanır.
- [ ] `whatsapp-bridge-mcp/tests/test_mcp_handlers.py` ve `whatsapp-tetikleyici/tests/test_*.py` (10 dosya) → `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/tests/` altında birleştirilir (dosya adı çakışması yoktur).
- [ ] Kapsam dışı: `bridge/whatsapp-bridge` (derlenmiş ikili), boş `bridge/cmd/`, `sidecar/node_modules/` (59 MB), `sidecar/rsrc-dump.txt` (316 KB), tetikleyicinin ölü `config.yaml`'ı (hiçbir kod okumuyor; kökteki `config.yaml` köprünün `appcfg.ProjectDir()` işaretidir).
- [ ] `config.yaml`, `tetikleyici.example.yaml`, `pasif-kurallar.example.yaml`, `config.readme.md`, `LICENSE` kök dizine yerleştirilir.

#### Adım 1.3: Canlı Kuralların, Durum Dosyalarının ve WhatsApp Oturumunun Aktarımı
- [ ] Canlı `tetikleyici.yaml` ve `pasif-kurallar.yaml` dosyaları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` köküne kopyalanır.
- [ ] Canlı `.tetikleyici/` dizini (`seen.txt`, `bildirim.txt`, `bekleyen_cevaplar.txt*`, `zamanlanmis/`, `tetikleyici.log`) `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/.tetikleyici/` altına kopyalanır; hatalı iç içe `.tetikleyici/.tetikleyici` klasörü **taşınmaz ve silinmez** (içinde yalnız orada bulunan `aktif_motor.txt`, `tunel_durumu.json`, `son_sessionlar.json` var - bkz. `NOTES.md` §6, `ASK.md` 2026-11-09 04:00).
- [ ] `~/.local/share/whatsapp-bridge/session.db` ve `messages.jsonl` dosyalarının yerinde olduğu doğrulanır (QR eşleşmesi bozulmadan korunur).

#### Adım 1.4: Belge Birleştirme, PII Temizliği ve TODO Sadeleştirmesi (TODO #8 Dahil)
- [ ] `docs/` altındaki tüm belgelerde PII taraması (`grep -rnE '90[0-9]{10}|[0-9]{14}@lid'` + `905551234567`/`12345678901234@lid` hariç tutulur) yapılır; gerçek telefon ve LID numaraları `905551234567` ve `12345678901234@lid` ile maskelenir ve `docs/` dizini `.gitignore`'dan çıkarılarak versiyon kontrolüne alınır.
- [ ] `README.md` Nerve Dispatch vizyonuyla güncellenir; `NOTES.md`, `ASK.md` ve `TODO.DONE.md` birleştirilir.
- [ ] `TODO.md` içindeki `[TAMAMLANDI ✅]` olan maddeler (Sprint #1–#6 ve eski bitmiş başlıklar) `TODO.DONE.md`'ye taşınır; kalıcı kurallar `NOTES.md`'ye aktarılır; `TODO.md` içinde yalnız açık işler bırakılır.

#### Adım 1.5: İnce `Makefile`, Modüler `install/` ve `upgrade/` Klasörlerinin İnşası ve 627+ Test Doğrulaması
- [ ] `bridge/internal/appcfg/appcfg.go` içindeki eski `Desktop/whatsapp-bridge-project` yedek yolları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` ile uyumlu hâle getirilir.
- [ ] Tekil `requirements.txt` ve ince yönlendirici `Makefile` (`make install` -> `./install/install.sh`, `make upgrade` -> `./upgrade/upgrade.sh`, `make build`, `make test`, `make lint`, `make verify`, `make clean`) oluşturulur.
- [ ] `install/` klasörü altındaki modüler kurulum betikleri (`install.sh`, `lib.sh`, `00-system-deps.sh`, `10-core-go.sh`, `20-python-env.sh`, `30-channels-mcp.sh`, `40-voice-modules.sh`, `50-vision-modules.sh`, `60-agents-config.sh`, `70-configs-services.sh`) oluşturulur.
- [ ] `upgrade/` klasörü altındaki modüler yükseltme betikleri (`upgrade.sh`, `00-backup-snapshot.sh`, `10-upgrade-core.sh`, `20-upgrade-channels-mcp.sh`, `30-upgrade-media-models.sh`, `40-migrate-configs.sh`, `50-verify-and-restart.sh`) oluşturulur.
- [ ] `make install` çalıştırılarak test/MCP için `.venv` kurulur (canlı hatlar FAZ 1'de sistem `python3` + `test_ASR/venv` ile çalışmaya devam eder, `NOTES.md` §6); `bin/whatsapp-bridge` ve `bin/tetikleyici` (`bin/nerve`) derlenir; ses ve görüntü modülleri doğrulanır.
- [ ] `make verify` çalıştırılarak hem `bridge` Go testlerinin, hem `tetikleyici` Go testlerinin, hem de 11 Python test dosyasının **%100 yeşil geçtiği** kanıtlanır.

#### Adım 1.6: Canlı Systemd Servislerinin `NerveDispatch` Reposuna Yönlendirilmesi (`install/70-configs-services.sh`)
- [ ] Derlenen `bin/tetikleyici` ikilisi atomik olarak (`cp` + `mv -f`) `/home/mfmx/bin/tetikleyici` (ve `/home/mfmx/bin/nerve`) konumuna kopyalanır.
- [ ] `~/.config/systemd/user/whatsapp-tetikleyici.service` içindeki `WorkingDirectory` ve `--config` yolları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` olarak güncellenip servis yeniden başlatılır.
- [ ] `/etc/systemd/system/whatsapp-bridge.service` içindeki `WorkingDirectory` ve `ExecStart` yollarının `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` olarak güncellenmesi için gerekli `sudo` komutu kullanıcıya sunulur.

#### Adım 1.7: Ajan Ekosistemi (MCP, Skill, INDEX) Bağlantılarının Güncellenmesi (`install/60-agents-config.sh`)
- [ ] **Antigravity (`~/.gemini/config/mcp_config.json`):** `whatsapp-bridge` → `PYTHONPATH` değeri `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yapılır.
- [ ] **OpenCode (`~/.config/opencode/opencode.json`):** `whatsapp-bridge` → `cwd` değeri `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yapılır; `wp-sahip` ve `wp-misafir` araç listelerine 5 yeni MCP aracı (`whatsapp_reaction`, `whatsapp_revoke`, `whatsapp_send_location`, `whatsapp_send_contact`, `whatsapp_vote_poll`) eklenir.
- [ ] **Claude Code (`~/.claude.json`):** `whatsapp-bridge` → `PYTHONPATH` değeri `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yapılır.
- [ ] **Skill & Fihrist:** `~/.gemini/config/skills/whatsapp-bridge/SKILL.md`, `~/.config/opencode/skills/whatsapp-bridge/SKILL.md` ve `~/.config/opencode/INDEX.md` dosyaları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yoluyla ve 17 MCP aracıyla güncellenir.
- [ ] `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` reposunda ilk temiz commit oluşturulur (eski iki repo, geçiş canlıda tamamen doğrulandıktan sonra arşivlenir/kaldırılır).

---

### FAZ 2: WhatsApp Hattını ve Çekirdek Yetenekleri Tamamlama (`NerveDispatch` Üzerinde Sıralı Geliştirme)

> **Prensip:** Önce WhatsApp kusursuzlaştırılır. Ancak yazılan her yeni modül (`drivers/`, `cli/`, `webui/`, `message` zarfı) `channel` alanını (`default: "whatsapp"`) barındıracak şekilde kodlanır.

#### Geliştirme 2.1: Tek Kalıcı Oturum (Sabit Session ID) ve Ajan Adaptör Mimarisi (TODO #7) [FAZ 2 İLK İŞ 🎯]
- [ ] `tetikleyici/drivers/` altında ortak `AgentDriver` arayüzü (`ListSessions()`, `Run(sessionID, prompt)`, `NewSession(prompt)`) kurulur:
  - **`opencode` sürücüsü:** `opencode run -s <id>` ve `opencode session list`.
  - **`agy` sürücüsü:** `agy -p "<prompt>" --conversation <id>` ile her mesajda sıfır bağlamlı yeni konuşma açmak yerine sabit konuşmayı sürdürme.
  - **`claude` sürücüsü:** İlk mesajda `--session-id <uuid>`, sonraki mesajlarda `--resume <uuid>` ile sabit oturumu sürdürme.
- [ ] `.tetikleyici/oturumlar.json` içinde motor başına (`opencode`, `agy`, `claude`) sabit oturum kimlikleri saklanır; hangi kanaldan mesaj gelirse gelsin aktif motorun sabit oturumuna yönlendirilir.
- [ ] `oturumKilidiAl` eşzamanlılık kilidi tüm ajan sürücülerine genişletilir; aynı anda gelen mesajlar sırayla işlenir.
- [ ] Kalıcı oturumda gereksiz geçmiş mesaj penceresi tekrarı (`buildPrompt`) sadeleştirilir ve `/session yeni` ile bağlam sıfırlama eklenir.
- [ ] `ASK.md`'de bekleyen `~/.gemini/antigravity/conversation_summaries.db` içindeki eski WhatsApp sohbet kirliliği (kullanıcı onayıyla) temizlenir.

#### Geliştirme 2.2: Birleşik `nerve` CLI (`status | qr | pair | logs | config | session | channel`) (TODO #13A)
- [ ] `tetikleyici` (`nerve`) ikilisine alt komut yönlendiricisi eklenir (bayraklı `--config` çağrıları doğrudan daemon modunda çalışmaya devam eder).
- [ ] `status`, `qr`, `logs -f`, `config check`, `session list`, `channel list`, `restart` alt komutları kodlanır ve testleri yazılır.

#### Geliştirme 2.3: Sıfır-Yüklü Gömülü Web Dashboard (`go:embed`) (TODO #13B)
- [ ] `bridge/internal/webui/` paketi oluşturulur; tek sayfalık modern HTML/CSS/JS kontrol paneli `//go:embed` ile ikiliye gömülür.
- [ ] `http://127.0.0.1:31415/dashboard` adresinde: Kanal durum kartları (WhatsApp, Telegram, Discord, Instagram), büyük QR kod ekranı, canlı SSE mesaj akışı ve aktif kural/yetki tablosu sunulur.

#### Geliştirme 2.4: QR Kodsuz 8 Haneli Eşleştirme Kodu (`Pairing Code`) (TODO #21)
- [ ] `bridge` içine whatsmeow `client.PairPhone(phoneNumber)` metodunu saran `POST /pair-phone` uç noktası eklenir.
- [ ] Hem `nerve pair 905...` CLI komutundan hem de Web Dashboard üzerinden 8 haneli kod (`ABCD-1234`) üretilmesi sağlanır.

#### Geliştirme 2.5: Canlı İlerleme Bildirimi ve Hızlı Tepki Emojisi (TODO #12)
- [ ] `tetikleyici.yaml` içindeki `tepki_emojisi` (`aktif`, `isleniyor: "👀"`, `basarili: "✅"`, `hatali: "❌"`) yapılandırması Go motoruna bağlanır.
- [ ] Yetkili kural eşleştiğinde mesaja anında `👀` bırakılır; işlem başarılıysa `✅`, hata oluşursa `❌` olarak güncellenir.
- [ ] 10 saniyeyi aşan işlemlerde ara durum bildirimi desteği eklenir.

#### Geliştirme 2.6: `/ekran` Masaüstü Ekran Görüntüsü Alma Komutu (TODO #10)
- [ ] `tetikleyici/cmd/` altına `/ekran [saniye]` komutu (`allow_ekran` yetkisiyle) eklenir.
- [ ] Arch Linux KDE Plasma Wayland üzerinde `spectacle -b -n -o /tmp/ekran.png` ile ekran görüntüsü yakalanıp sohbete fotoğraf olarak gönderilir.

#### Geliştirme 2.7: İnteraktif Onay ve İzin Sistemi (TODO #9)
- [ ] Tehlikeli eylemler (`rm`, `git push --force`, `systemctl` vb.) öncesi 4 seçenekli onay anketi (`Tek seferlik onayla`, `Reddet`, `1 saat otomatik onayla`, `1 hafta otomatik onayla`) gönderilir ve `.tetikleyici/izinler.json` içinde TTL ile saklanır.

#### Geliştirme 2.8: Çift Yönlü Dosya Transferi (`/get <yol>` ve Gelen Dosya Kaydı) (TODO #14)
- [ ] `/get <dosya_yolu>` ile bilgisayardaki dosyaların (≤ 50 MB) sohbete belge olarak gönderilmesi sağlanır.
- [ ] Sohbetten bilgisayara atılan dosyaların `~/Downloads/nervedispatch-files/` dizinine `chmod -x` güvenli izinleriyle kaydedilmesi ve ajana yolunun bildirilmesi sağlanır.

#### Geliştirme 2.9: Dinamik Çalışma Dizini (`/root`), Pano (`/clip`) ve Süreç Yönetimi (`/apps`, `/kill`) (TODO #15, #16, #17)
- [ ] `/root <dizin>`: Aktif ajan ve `/cmd` çalışma dizinini dinamik değiştirme.
- [ ] `/clip` ve `/clip <metin>`: Wayland panosunu (`wl-paste` / `wl-copy`) uzaktan okuma ve yazma.
- [ ] `/apps` ve `/kill <pid|isim>`: Açık pencereleri listeleme ve kritik servis kara listesiyle askıda kalan uygulamaları kapatma.

#### Geliştirme 2.10: Bağlam Sıkıştırma (`/compact`), Akıllı Zamanlayıcı ve Sistem Bekçisi (TODO #11, #19, #20)
- [ ] `/compact` (`/ozet`): Uzun süren ajan oturumlarında bağlamı özetleyip hafifletme.
- [ ] Gelişmiş `/zamanla`: Ajanlı zamanlanmış iş akışları (`/zamanla her sabah 08:30 ...`).
- [ ] Proaktif Sistem Bekçisi (`bekci:`): Disk (%90), RAM (%95), CPU sıcaklığı (85°C) ve kritik servis düşüşlerinde 1 saat cooldown korumalı otomatik alarm.

#### Geliştirme 2.11: Yerel TTS Sesli Yanıt (Piper TTS) ve Kod/Klasör Refactoring (TODO #18 + Kod Kuralları)
- [ ] `Piper TTS` ile CPU üzerinde sıfır VRAM tüketen Türkçe sesli yanıt (`.ogg opus` sesli not gönderimi, `/tts on|off`).
- [ ] 300 satırı aşan üretim dosyalarının (`action.go`, `main.go`, `config.go`, `engine.go`, `calistir.py`) ve 15 dosya sınırını aşan klasörlerin (`tetikleyici/`, `wabridge/`) modüler alt dosyalara bölünmesi.

---

### FAZ 3: Çok Kanallı Genişleme (Telegram, Discord, Instagram Entegrasyonu)

WhatsApp hattı ve FAZ 2 çekirdeği tamamlandıktan sonra, `bridge/internal/channels/` altına sırasıyla diğer kanal adaptörleri eklenecektir:

#### Geliştirme 3.1: Ortak Kanal Adaptör Arayüzü (`ChannelAdapter`) ve Çok Kanallı Kural Eşleştirici
- [ ] `bridge/internal/channels/adapter.go` içinde standart `ChannelAdapter` arayüzü (`Connect`, `Disconnect`, `SendText`, `SendMedia`, `React`, `SendPoll`, `Health`) tanımlanır; mevcut `wabridge` bu arayüzün ilk uygulaması (`whatsapp`) olur.
- [ ] `tetikleyici/match.go` içine `channel` (`whatsapp`, `telegram`, `discord`, `instagram`) ve kanal kimlikleri (`tg_id`, `dc_id`, `ig_id`) eşleştirme desteği eklenir.
- [ ] `install/30-channels-mcp.sh` ve `upgrade/20-upgrade-channels-mcp.sh` içine yeni kanalların kurulum ve yükseltme kancaları eklenir.

#### Geliştirme 3.2: Telegram Kanal Adaptörü (`channels/telegram`) [İLK EK KANAL]
- [ ] Go tabanlı Telegram adaptörü (`telebot` / Bot API veya MTProto) eklenir; gelen mesajlar, sesli notlar, fotoğraflar, belgeler ve buton tıklamaları SSE `/events` akışına `channel: "telegram"` olarak basılır.
- [ ] Mefamex Telegram ID (`975678139`) `tetikleyici.yaml` içinde yetkili admin olarak tanımlanır; WhatsApp'taki tüm `/cmd`, `/camera`, `/ekran`, `/agent`, `/durum` komutları ve sabit AI oturumları Telegram üzerinden de birebir çalışır.
- [ ] **Kanal-Arası Yedeklilik (Cross-Channel Failover):** WhatsApp bağlantısı koptuğunda (`on_bridge_down`) veya yeniden QR gerektiğinde bildirim otomatik olarak Telegram üzerinden Mefamex'e iletilir.

#### Geliştirme 3.3: Discord Kanal Adaptörü (`channels/discord`)
- [ ] `discordgo` tabanlı WebSocket Gateway adaptörü eklenir (DM ve yetkili sunucu kanalları).
- [ ] Kod blokları, sistem durumu (`/durum`) ve ajan çıktıları Discord zengin `Embed` formatıyla desteklenir.

#### Geliştirme 3.4: Instagram Kanal Adaptörü (`channels/instagram`)
- [ ] Instagram DM metin, ses ve görsel mesajlarını aynı `ChimeraEvent` otobanına bağlayan adaptör eklenir.
- [ ] Tüm kanalların canlı durumu, mesaj istatistikleri ve bağlantı kontrolleri Gömülü Web Dashboard (`127.0.0.1:31415/dashboard`) ve `nerve channel list` üzerinde birleştirilir.

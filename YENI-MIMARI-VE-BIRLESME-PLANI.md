# Nerve Dispatch (NerveDispatch) — Çok Kanallı Mimari ve Birleşme Planı (Master Blueprint)

> **Proje Adı:** `NerveDispatch` (`chimera` — 3. Jenerasyon)  
> **Hedef Konum:** `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/`  
> **Tarihsel & Felsefi Arka Plan:** [`/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md`](file:///mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md)  
> **Birleştirilecek Kaynak Repolar:**
> 1. `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-bridge-mcp`
> 2. `/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-tetikleyici`  
> **Durum:** Onaylandı — Sonraki Sohbette Doğrudan Uygulamaya Geçilecek (FAZ 1 → FAZ 2 → FAZ 3)  
> **Son Güncelleme:** 2026-11-09

---

## 🤖 YENİ SOHBETTE BAŞLAYACAK AJAN İÇİN DEVİR-TESLİM TALİMATI (AGENT HANDOFF)

Eğer yeni bir oturumda kullanıcı **"Nerve Dispatch birleşmesine başlayalım"** veya **"plana göre devam edelim"** derse:
1. **Önce şu iki belgeyi oku:**
   - Bu belge (`YENI-MIMARI-VE-BIRLESME-PLANI.md`) — Teknik mimari ve adım adım geçiş/geliştirme planı.
   - [`/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md`](file:///mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/NERVE-DISPATCH-THE-STORY.md) — Projenin v1 ve v2'den bugüne evrimi ve felsefesi.
2. **Stratejik Sırayı Bozmamaya Dikkat Et:**
   - **Eski repolarda (`whatsapp-bridge-mcp` ve `whatsapp-tetikleyici`) artık yeni özellik kodlanmayacaktır.**
   - **İlk iş (FAZ 1):** İki kaynak repoyu `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` içine birleştirmek, canlı `tetikleyici.yaml`, `.tetikleyici/` durum dosyalarını taşımak, `~/.local/share/whatsapp-bridge/session.db` oturumunu (QR bozmadan) korumak, `make verify` ile 627+ testi yeşil geçirmek ve `systemd` + MCP ayarlarını `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` yoluna bağlamaktır.
   - **İkinci iş (FAZ 2):** `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` üzerinde **önce WhatsApp hattını ve tüm çekirdek yetenekleri** (`Geliştirme 2.1` Sabit Session ID & Ajan Sürücülerinden başlayarak `2.11`'e kadar) sırasıyla tamamlamaktır.
   - **Üçüncü iş (FAZ 3):** Çekirdek tamamlandıktan sonra **Telegram (`975678139`), Discord ve Instagram** kanallarını eklemektir.
3. **Kırmızı Çizgileri Unutma:**
   - `sudo` komutlarını kendi başına çalıştırma; kullanıcıya tek satırda göster.
   - `~/.local/share/whatsapp-bridge/session.db` ve `messages.jsonl` dosyalarını asla silme veya sıfırlama (canlı WhatsApp eşleşmesi oradadır).
   - Python MCP klasörünün adı modül gölgelemesini (`from mcp.types import ...`) önlemek için mutlaka `mcp_server/` olarak kalmalıdır.

---

## 1. Chimera v3 Vizyonu: Çok Başlı, Çok Kanallı, Tek Zihinli Sistem

1. **Çok Kanallı Duyu Organları (Omnichannel Gateway):**
   - **WhatsApp** *(Birincil & Aktif Kanal — Whatsmeow Go soketi)*
   - **Telegram** *(İkincil Kanal — Faz 3.2, Mefamex Telegram ID: `975678139`)*
   - **Discord** *(Üçüncül Kanal — Faz 3.3, Gateway WebSocket & Slash Komutları)*
   - **Instagram** *(Dördüncül Kanal — Faz 3.4, DM Otomasyonu)*
2. **Tek Ortak Beyin ve Olay Dağıtıcı (Unified Dispatcher & Persistent Memory):**
   - Mesaj hangi kanaldan gelirse gelsin (WhatsApp, Telegram, Discord, Instagram), tek bir **Ortak Olay Zarfına (`ChimeraEvent`)** dönüştürülür.
   - **Kanal-Arası Ortak Ajan Hafızası:** Kullanıcı dışarıdayken WhatsApp'tan başlattığı bir kodlama veya sistem görevine, masasına geçtiğinde Telegram veya Discord üzerinden **aynı sabit AI oturumu (`opencode`, `agy`, `claude`)** içinde kesintisiz devam edebilir.
3. **Kanal-Arası Yedeklilik (Cross-Channel Failover):**
   - WhatsApp bağlantısı düştüğünde (`on_bridge_down`), WhatsApp kopuk olduğu için uyarı mesajı otomatik olarak **Telegram** üzerinden Mefamex'e (`975678139`) iletilir.
4. **Hafif, Güvenli ve Her Cihazda Çalışır:**
   - Ağır masaüstü GUI (Electron/Qt) içermez. **Headless Go Daemon (~30 MB RAM) + Birleşik `chimera` CLI + Go `embed` Gömülü Web Dashboard + ChatOps** mimarisiyle hem kişisel masaüstünde (Arch Linux KDE Plasma Wayland) hem de ekransız sunucularda (VPS, Raspberry Pi) çalışır.

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

## 3. Chimera v3.0 — Çok Kanallı Hedef Mimari Şeması

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
│ 1. KATMAN: CHIMERA KÖPRÜ MOTORU — bridge/ (Go Headless Daemon :31415)                        │
│  - Kanal Adaptörleri (ChannelAdapter: WhatsApp [Aktif], Telegram, Discord, Instagram)        │
│  - Ortak Olay Zarfı Üreticisi (Unified ChimeraEvent: channel, id, from, type, text, media)   │
│  - Kalıcı Oturum & Mesaj Deposu (session.db + messages.jsonl, ~30 MB RAM)                    │
│  - Gömülü Web Sunucusu (go:embed Web Dashboard: http://127.0.0.1:31415/dashboard)            │
│  - Birleşik HTTP REST API (/send, /reaction, /poll...) & Gerçek Zamanlı SSE (/events)        │
│  - FastMCP STDIO Sunucusu (mcp_server/ — AI Ajanları için Bildirim ve İletişim Araçları)     │
└───────────────────────────────┬──────────────────────────────────────────────┬───────────────┘
                                │ Gerçek Zamanlı SSE Olay Akışı (/events)      │
                                ▼                                              ▼
┌───────────────────────────────────────────────────────────────┐ ┌────────────────────────────┐
│ 2. KATMAN: CHIMERA DAĞITICI — tetikleyici/ (systemd --user)   │ │  Yönetim & Gözlem Katmanı  │
│  - Çok Kanallı Kimlik & RBAC (WP LID/Tel, TG ID, DC ID)       │ │ - Birleşik `chimera` CLI   │
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

## 5. `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` — Modüler Monorepo Dizin Hiyerarşisi

Mevcut testlerin (`realkonfig_test.go`, `shipped_config_test.go`, `tests/test_*.py`) beklediği göreli dosya yollarını (`../tetikleyici.yaml`, `parent.parent.parent`) hiç bozmayan, aynı zamanda çok kanallı genişlemeye hazır dizin yapısı:

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
│   ├── drivers/                    # [YENİ - Faz 2.1] Ajan Sürücüleri (opencode, agy, claude — Sabit Oturum)
│   ├── cli/                        # [YENİ - Faz 2.2] Birleşik CLI (status, qr, pair, logs, session, channel)
│   ├── cmd/                        # Python: Sistem komutları (/cmd, /cmdhub, /camera, /ekran, /durum, /zamanla...)
│   ├── go.mod & go.sum             # module tetikleyici
│   └── *_test.go                   # 27 Go birim ve entegrasyon test dosyası
│
├── mcp_server/                     # Python: FastMCP STDIO Sunucusu (`mcp` kütüphanesiyle çakışmaz)
│   ├── __main__.py                 # python -m mcp_server
│   ├── api/client.py               # Bridge HTTP REST istemcisi
│   ├── config/                     # Yapılandırma yükleyici
│   ├── server/                     # FastMCP sunucu tanımı
│   └── tools/                      # 16+ MCP aracı (WhatsApp + Çok Kanallı bildirim araçları)
│
├── gorsel/                         # Python: Görsel, Belge (PDF/DOCX/XLSX), OCR, Video ve Sticker Paketi
├── gorsel_pipeline.py              # Görsel analiz boru hattı (Ollama Qwen2.5-VL + OCR)
├── voice_pipeline.py               # Sesli not transkripsiyon boru hattı (Faster-Whisper GPU/CPU + VRAM tahliye)
│
├── sidecar/                        # Opsiyonel Node.js VoIP / Arama WASM Yardımcı Modülü
├── docs/                           # Tasarım belgeleri, arşiv ve teknik analizler (PII temizlenmiş)
├── tests/                          # Birleşik Python Test Süiti (test_mcp_handlers.py + 10 tetikleyici testi)
│
├── bin/                            # Derlenmiş Go ikilileri (gitignore'da)
│   ├── whatsapp-bridge             # Köprü ikilisi (chimera-bridge)
│   └── tetikleyici                 # Dağıtıcı & CLI ikilisi (chimera / tetikleyici)
│
├── .tetikleyici/                   # Çalışma zamanı durum dizini (gitignore'da: seen.txt, oturumlar.json, log...)
├── go.work                         # Go Workspace dosyası (use ./bridge ./tetikleyici)
├── requirements.txt                # Tekil Python bağımlılıkları (.venv)
├── config.yaml                     # Köprü, kanal ve MCP çalışma parametreleri
├── tetikleyici.example.yaml        # Kural ve yetki şablonu (Git'te takip edilir)
├── tetikleyici.yaml                # Canlı kişisel kural ve yetki dosyası (gitignore'da)
├── pasif-kurallar.example.yaml     # Pasif misafir kuralları şablonu (Git'te takip edilir)
├── pasif-kurallar.yaml             # Canlı pasif kurallar dosyası (gitignore'da)
├── config.readme.md                # Kapsamlı yapılandırma kılavuzu
├── Makefile                        # Tek merkezden derleme, kurulum, test, lint ve servis orkestratörü
├── README.md                       # Chimera v3.0 proje vitrini ve mimari şeması
├── TODO.md                         # Yalnızca açık işleri barındıran aktif sprint listesi
├── TODO.DONE.md                    # Tamamlanmış ve kanıtlanmış işler arşivi
├── NOTES.md                        # Kalıcı teknik bellek ve sistem kuralları
└── ASK.md                          # Mimari soru ve otonom karar kaydı
```

---

## 6. Birleşik CLI (`chimera` / `tetikleyici`) ve Gömülü Web Dashboard

### 6.1. CLI Komut Seti
Hem `chimera` hem de `tetikleyici` komutu aynı ikiliyi çalıştırır (geriye dönük uyumluluk tamdır):
- `chimera` (veya `tetikleyici --config ...`): Varsayılan daemon modunda olay dağıtıcıyı başlatır.
- `chimera status`: Köprü, dağıtıcı, bağlı kanallar (WhatsApp, Telegram, Discord, Instagram), aktif AI motoru ve bellek durumunu gösterir.
- `chimera qr` / `chimera pair <numara>`: Terminalde ANSI QR kod çizer veya 8 haneli telefon eşleştirme kodu üretir.
- `chimera logs [-f] [n]`: Tüm kanalların ve dağıtıcının günlüklerini UTF-8 güvenli şekilde akıtır.
- `chimera session list`: `opencode`, `agy` ve `claude` üzerindeki sabit ve açık oturumları listeler.
- `chimera channel list`: Yapılandırılmış kanalların bağlantı durumunu listeler.
- `chimera config check`: YAML dosyalarını, betik yollarını ve yetkileri doğrular.
- `chimera restart`: Servisleri güvenli şekilde yeniden başlatır.

### 6.2. Gömülü Yerel Web Dashboard (`http://127.0.0.1:31415/dashboard`)
- `bridge/internal/webui/` içinde `//go:embed` ile tek binary içine gömülü, **+0 MB RAM** maliyetli kontrol paneli.
- **Bölümler:**
  1. **Kanal Durum Kartları:** WhatsApp (`paired`/`connected`), Telegram, Discord, Instagram bağlantı rozetleri.
  2. **Eşleştirme (Pairing) Merkezi:** Büyük QR kod ve 8 haneli eşleştirme kodu alanı.
  3. **Canlı Olay Akışı (SSE Stream):** Gelen/giden mesajların, kanal etiketleriyle (`[WA]`, `[TG]`, `[DC]`) anlık izlenmesi.
  4. **Ajan & Kural Özeti:** Aktif AI motoru, sabit oturum ID'leri ve yetki matrisi.

---

## 7. Korunacak Temel Kurallar ve Sistem Sözleşmeleri

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

## 8. Uygulama ve Geçiş Yol Haritası (3 Fazlı Master Roadmap)

> **Uygulama Prensibi:**
> - **İlk Adım (FAZ 1):** Mevcut iki repoyu doğrudan `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` çatısında birleştirmek ve tüm canlı sistemi sıfır kesintiyle `NerveDispatch` üzerine taşımak.
> - **İkinci Adım (FAZ 2):** **Önce WhatsApp hattını ve tüm çekirdek özellikleri eksiksiz bitirmek.**
> - **Üçüncü Adım (FAZ 3):** Hazır ve kaya gibi sağlam çekirdek üzerine Telegram, Discord ve Instagram kanallarını eklemek.

---

### FAZ 1: Doğrudan `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` Monoreposuna Geçiş ve Ekosistem Onarımı [SIRADA #1 🎯]

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
- [ ] `config.yaml`, `tetikleyici.example.yaml`, `pasif-kurallar.example.yaml`, `config.readme.md`, `LICENSE` kök dizine yerleştirilir.

#### Adım 1.3: Canlı Kuralların, Durum Dosyalarının ve WhatsApp Oturumunun Aktarımı
- [ ] Canlı `tetikleyici.yaml` ve `pasif-kurallar.yaml` dosyaları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` köküne kopyalanır.
- [ ] Canlı `.tetikleyici/` dizini (`seen.txt`, `bildirim.txt`, `bekleyen_cevaplar.txt*`, `zamanlanmis/`, `tetikleyici.log`) `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/.tetikleyici/` altına kopyalanır; hatalı iç içe `.tetikleyici/.tetikleyici` klasörü temizlenir.
- [ ] `~/.local/share/whatsapp-bridge/session.db` ve `messages.jsonl` dosyalarının yerinde olduğu doğrulanır (QR eşleşmesi bozulmadan korunur).

#### Adım 1.4: Belge Birleştirme, PII Temizliği ve TODO Sadeleştirmesi (TODO #8 Dahil)
- [ ] `docs/` altındaki tüm belgelerde PII taraması (`grep -rnE "9054|9053|475759|148657|194519"`) yapılır; gerçek telefon ve LID numaraları `905551234567` ve `12345678901234@lid` ile maskelenir ve `docs/` dizini `.gitignore`'dan çıkarılarak versiyon kontrolüne alınır.
- [ ] `README.md` Chimera v3.0 vizyonuyla güncellenir; `NOTES.md`, `ASK.md` ve `TODO.DONE.md` birleştirilir.
- [ ] `TODO.md` içindeki `[TAMAMLANDI ✅]` olan maddeler (Sprint #1–#6 ve eski bitmiş başlıklar) `TODO.DONE.md`'ye taşınır; kalıcı kurallar `NOTES.md`'ye aktarılır; `TODO.md` içinde yalnız açık işler bırakılır.

#### Adım 1.5: Tekil Ortam Kurulumu, Kod İçi Yol Onarımı ve 627+ Test Doğrulaması
- [ ] `bridge/internal/appcfg/appcfg.go` içindeki eski `Desktop/whatsapp-bridge-project` yedek yolları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` ile uyumlu hâle getirilir.
- [ ] Tekil `requirements.txt` ve birleşik `Makefile` (`make install`, `make build`, `make test`, `make lint`, `make verify`, `make clean`) oluşturulur.
- [ ] `make install` ile tek `.venv` kurulur; `bin/whatsapp-bridge` ve `bin/tetikleyici` (`bin/chimera`) derlenir.
- [ ] `make verify` çalıştırılarak hem `bridge` Go testlerinin, hem `tetikleyici` Go testlerinin, hem de 11 Python test dosyasının **%100 yeşil geçtiği** kanıtlanır.

#### Adım 1.6: Canlı Systemd Servislerinin `NerveDispatch` Reposuna Yönlendirilmesi
- [ ] Derlenen `bin/tetikleyici` ikilisi atomik olarak (`cp` + `mv -f`) `/home/mfmx/bin/tetikleyici` (ve `/home/mfmx/bin/chimera`) konumuna kopyalanır.
- [ ] `~/.config/systemd/user/whatsapp-tetikleyici.service` içindeki `WorkingDirectory` ve `--config` yolları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` olarak güncellenip servis yeniden başlatılır.
- [ ] `/etc/systemd/system/whatsapp-bridge.service` içindeki `WorkingDirectory` ve `ExecStart` yollarının `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` olarak güncellenmesi için gerekli `sudo` komutu kullanıcıya sunulur.

#### Adım 1.7: Ajan Ekosistemi (MCP, Skill, INDEX) Bağlantılarının Güncellenmesi
- [ ] **Antigravity (`~/.gemini/config/mcp_config.json`):** `whatsapp-bridge` → `PYTHONPATH` değeri `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yapılır.
- [ ] **OpenCode (`~/.config/opencode/opencode.json`):** `whatsapp-bridge` → `cwd` değeri `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yapılır; `wp-sahip` ve `wp-misafir` araç listelerine 5 yeni MCP aracı (`whatsapp_react`, `whatsapp_revoke`, `whatsapp_send_location`, `whatsapp_send_contact`, `whatsapp_vote_poll`) eklenir.
- [ ] **Claude Code (`~/.claude.json`):** `whatsapp-bridge` → `PYTHONPATH` değeri `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yapılır.
- [ ] **Skill & Fihrist:** `~/.gemini/config/skills/whatsapp-bridge/SKILL.md`, `~/.config/opencode/skills/whatsapp-bridge/SKILL.md` ve `~/.config/opencode/INDEX.md` dosyaları `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` yoluyla ve 16 MCP aracıyla güncellenir.
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

#### Geliştirme 2.2: Birleşik `chimera` CLI (`status | qr | pair | logs | config | session | channel`) (TODO #13A)
- [ ] `tetikleyici` (`chimera`) ikilisine alt komut yönlendiricisi eklenir (bayraklı `--config` çağrıları doğrudan daemon modunda çalışmaya devam eder).
- [ ] `status`, `qr`, `logs -f`, `config check`, `session list`, `channel list`, `restart` alt komutları kodlanır ve testleri yazılır.

#### Geliştirme 2.3: Sıfır-Yüklü Gömülü Web Dashboard (`go:embed`) (TODO #13B)
- [ ] `bridge/internal/webui/` paketi oluşturulur; tek sayfalık modern HTML/CSS/JS kontrol paneli `//go:embed` ile ikiliye gömülür.
- [ ] `http://127.0.0.1:31415/dashboard` adresinde: Kanal durum kartları (WhatsApp, Telegram, Discord, Instagram), büyük QR kod ekranı, canlı SSE mesaj akışı ve aktif kural/yetki tablosu sunulur.

#### Geliştirme 2.4: QR Kodsuz 8 Haneli Eşleştirme Kodu (`Pairing Code`) (TODO #21)
- [ ] `bridge` içine whatsmeow `client.PairPhone(phoneNumber)` metodunu saran `POST /pair-phone` uç noktası eklenir.
- [ ] Hem `chimera pair 905...` CLI komutundan hem de Web Dashboard üzerinden 8 haneli kod (`ABCD-1234`) üretilmesi sağlanır.

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
- [ ] Sohbetten bilgisayara atılan dosyaların `~/Downloads/chimera-files/` dizinine `chmod -x` güvenli izinleriyle kaydedilmesi ve ajana yolunun bildirilmesi sağlanır.

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

#### Geliştirme 3.2: Telegram Kanal Adaptörü (`channels/telegram`) [İLK EK KANAL]
- [ ] Go tabanlı Telegram adaptörü (`telebot` / Bot API veya MTProto) eklenir; gelen mesajlar, sesli notlar, fotoğraflar, belgeler ve buton tıklamaları SSE `/events` akışına `channel: "telegram"` olarak basılır.
- [ ] Mefamex Telegram ID (`975678139`) `tetikleyici.yaml` içinde yetkili admin olarak tanımlanır; WhatsApp'taki tüm `/cmd`, `/camera`, `/ekran`, `/agent`, `/durum` komutları ve sabit AI oturumları Telegram üzerinden de birebir çalışır.
- [ ] **Kanal-Arası Yedeklilik (Cross-Channel Failover):** WhatsApp bağlantısı koptuğunda (`on_bridge_down`) veya yeniden QR gerektiğinde bildirim otomatik olarak Telegram üzerinden Mefamex'e iletilir.

#### Geliştirme 3.3: Discord Kanal Adaptörü (`channels/discord`)
- [ ] `discordgo` tabanlı WebSocket Gateway adaptörü eklenir (DM ve yetkili sunucu kanalları).
- [ ] Kod blokları, sistem durumu (`/durum`) ve ajan çıktıları Discord zengin `Embed` formatıyla desteklenir.

#### Geliştirme 3.4: Instagram Kanal Adaptörü (`channels/instagram`)
- [ ] Instagram DM metin, ses ve görsel mesajlarını aynı `ChimeraEvent` otobanına bağlayan adaptör eklenir.
- [ ] Tüm kanalların canlı durumu, mesaj istatistikleri ve bağlantı kontrolleri Gömülü Web Dashboard (`127.0.0.1:31415/dashboard`) ve `chimera channel list` üzerinde birleştirilir.

# CHIMERA: Çok Başlı Bir Fikrin Evrimi
## Bir Yerel Yapay Zeka, Otonom Ajan ve Çok Kanallı Cihaz Ağ Geçidinin Yıllara Yayılan Mimari Hikayesi (v1 • v2 • v3)

> **Yazar:** Mehmet Akif Akkoç ([@Mefamex](https://github.com/Mefamex))  
> **Tarih:** 2024 – 2026  
> **Konsept:** Local-First AI, Otonom Ajan Mimarisi, Çok Dilli (Polyglot) Sistemler, Çok Kanallı Ağ Geçidi (Omnichannel Gateway)  
> **Konum:** `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/`  
> **Yeni İsim:** **Nerve Dispatch** (`NerveDispatch`)

---

## İçindekiler

- [Giriş: Neden "Chimera"?](#giriş-neden-chimera)
- [Bölüm 1: Tohumlar ve İlk Kıvılcımlar (2024 – 2025 Başı)](#bölüm-1-tohumlar-ve-ilk-kıvılcımlar-2024--2025-başı)
  - [1.1 Cihazlar Arası Köprü Arayışı (`ftp_transfer_file`)](#11-cihazlar-arası-köprü-arayışı-ftp_transfer_file)
  - [1.2 Telefondan Terminale İlk Tünel: `gokkusagi_protokolu`](#12-telefondan-terminale-ilk-tünel-gokkusagi_protokolu)
- [Bölüm 2: Asistan Denemeleri ve İlk Duvara Çarpmalar (2025 Yaz & Güz)](#bölüm-2-asistan-denemeleri-ve-ilk-duvara-çarpmalar-2025-yaz--güz)
  - [2.1 İlk Yerel Kodlama Ajanı: `ai-code-agent`](#21-ilk-yerel-kodlama-ajanı-ai-code-agent)
  - [2.2 Donanım ve Tarayıcı Kumandası: `deviceAssistant`](#22-donanım-ve-tarayıcı-kumandası-deviceassistant)
  - [2.3 Bulut Çıkmazı ve İlk WhatsApp Vizyonu: `assistant-bot-backend`](#23-bulut-çıkmazı-ve-ilk-whatsapp-vizyonu-assistant-bot-backend)
- [Bölüm 3: CHIMERA v1 (Şubat 2026) — Büyük Düş, Büyük Mimari ve Büyük Tıkanma](#bölüm-3-chimera-v1-şubat-2026--büyük-düş-büyük-mimari-ve-büyük-tıkanma)
  - [3.1 Model Vekili: `ollama_bridge`](#31-model-vekili-ollama_bridge)
  - [3.2 Üç Parçalı (Polyglot) Mimari: Beyin, Eller ve Gözler](#32-üç-parçalı-polyglot-mimari-beyin-eller-ve-gözler)
  - [3.3 Kendi Kendini Denetleyen ReAct Döngüsü (`engine.go`)](#33-kendi-kendini-denetleyen-react-döngüsü-enginego)
  - [3.4 Şeffaf Zeka (Glass Box AI) ve Canlı WebSocket](#34-şeffaf-zeka-glass-box-ai-ve-canlı-websocket)
  - [3.5 8 Şubat 2026 Günlüğü: O Tarihi Kehanet](#35-8-şubat-2026-günlüğü-o-tarihi-kehanet)
  - [3.6 Neden Tıkandı? Yerel Modellerin Format Krizi](#36-neden-tıkandı-yerel-modellerin-format-krizi)
- [Bölüm 4: Arayış ve Uydu Projeler Dönemi (Haziran – Ekim 2026)](#bölüm-4-arayış-ve-uydu-projeler-dönemi-haziran--ekim-2026)
  - [4.1 CHIMERA2 (v2): Minimalist Geri Çekilme](#41-chimera2-v2-minimalist-geri-çekilme)
  - [4.2 Ses Organı: `test_ASR` ve Türkçe Whisper Benchmarkları](#42-ses-organı-test_asr-ve-türkçe-whisper-benchmarkları)
  - [4.3 Komuta Organı: `whatsapp-cmd` ve Shell Tünelleme](#43-komuta-organı-whatsapp-cmd-ve-shell-tünelleme)
- [Bölüm 5: Kırılma Anı (Ekim – Kasım 2026) — İki Parçalı Köprünün İnşası](#bölüm-5-kırılma-anı-ekim--kasım-2026--iki-parçalı-köprünün-inşası)
  - [5.1 Katman 1 & 2: `whatsapp-bridge-mcp` (Go & FastMCP)](#51-katman-1--2-whatsapp-bridge-mcp-go--fastmcp)
  - [5.2 Katman 3 & 4: `whatsapp-tetikleyici` (Dispatcher & Multimodal)](#52-katman-3--4-whatsapp-tetikleyici-dispatcher--multimodal)
  - [5.3 Aydınlanma Anı: "Beyni Sıfırdan Yazmak Zorunda Değilim!"](#53-aydınlanma-anı-beyni-sıfırdan-yazmak-zorunda-değilim)
- [Bölüm 6: CHIMERA v3 (Kasım 2026) — Büyük Birleşme ve Çok Kanallı Gelecek](#bölüm-6-chimera-v3-kasım-2026--büyük-birleşme-ve-çok-kanallı-gelecek)
  - [6.1 Monorepo Dönüşümü (`go.work` ve Ortak `.venv`)](#61-monorepo-dönüşümü-gowork-ve-ortak-venv)
  - [6.2 WhatsApp Bir Varış Noktası Değil, Sadece İlk Kapı](#62-whatsapp-bir-varış-noktası-değil-sadece-ilk-kapı)
  - [6.3 Evrensel Çok Kanallı Ağ Geçidi (Omnichannel Gateway) Mimarisi](#63-evrensel-çok-kanallı-ağ-geçidi-omnichannel-gateway-mimarisi)
- [Bölüm 7: Zaman Çizelgesi ve Karşılaştırma Matrisi](#bölüm-7-zaman-çizelgesi-ve-karşılaştırma-matrisi)
- [Bölüm 8: Çıkarılan Mühendislik Dersleri](#bölüm-8-çıkarılan-mühendislik-dersleri)
- [Son Söz: Ejderha Uyanıyor](#son-söz-ejderha-uyanıyor)

---

## Giriş: Neden "Chimera"?

Yunan mitolojisinde **Khimaira (Chimera)**; aslan başı, keçi gövdesi ve yılan kuyruğuna sahip, farklı doğaların tek bir organizmada kusursuz bir uyumla birleştiği efsanevi bir yaratıktır.

Yazılım mühendisliğinde çoğu geliştirici tek bir paradigmaya, tek bir programlama diline veya tek bir platforma sıkışıp kalır. Ancak gerçek dünya problemleri monolitik değildir:
- Bir sistem servisi yazıyorsanız; **Go** dilinin hafifliğine, tek ikili dosya (single-binary) mimarisine, eşzamanlılığına (goroutine) ve bellek verimliliğine ihtiyaç duyarsınız.
- Yapay zeka modelleriyle, ses transkripsiyonuyla, görüntü işleme boru hatlarıyla uğraşıyorsanız; **Python** ekosisteminin zenginliği rakipsizdir.
- Modern web sayfalarını kazımak, karmaşık JavaScript arayüzlerini otomatize etmek istiyorsanız; **Node.js** ve headless tarayıcılar şarttır.
- Ve tüm bu gücü sokakta yürürken cebinizden yönetmek istiyorsanız; karmaşık kurumsal VPN'ler ya da mobil terminal ekranları yerine herkesin cebinde açık olan **WhatsApp veya Telegram** gibi mesajlaşma protokollerine ihtiyacınız vardır.

İşte **CHIMERA**, bu farklı güçlerin birleşimidir. Bu belge; 2024 yılında basit bir SSH tüneli arayışıyla başlayan, 2025'te bulut kısıtlarına toslayan, Şubat 2026'da yerel yapay zeka modellerinin sınırlarına çarpıp duraklayan ve nihayetinde Kasım 2026'da **CHIMERA v3** adıyla çok kanallı, otonom bir cihaz ağ geçidine dönüşen **2 yıllık bir mimari evrimin** hakiki hikayesidir.

---

## Bölüm 1: Tohumlar ve İlk Kıvılcımlar (2024 – 2025 Başı)

Her büyük mimari, gündelik hayatın yarattığı basit bir sürtünmeyle (friction) başlar. Benim karşılaştığım sürtünme çok netti:  
*“Evdeki güçlü Linux iş istasyonuma, GPU gücüme ve yerel araçlarıma; dışarıdayken, kampüsteyken veya telefondayken nasıl aracısız ve güvenle hükmedebilirim?”*

### 1.1 Cihazlar Arası Köprü Arayışı (`ftp_transfer_file`)
Ocak 2024'te yazdığım `ftp_transfer_file` projesi, yerel ağda (LAN) veya IPv4 üzerinden iki cihaz arasında dosya taşımayı hedefleyen basit bir Python soket çözümüydü. Henüz ortada ne yapay zeka ne de ajan kavramı vardı; yalnızca cihazlar arasında görünmez bir veri tüneli kurma içgüdüsü vardı.

### 1.2 Telefondan Terminale İlk Tünel: `gokkusagi_protokolu`
Ağustos 2024'te projelendirilen ve 2026'da olgunlaştırılan **Gökkuşağı Protokolü**, bu yolculuğun ilk gerçek altyapı taşıydı:

```text
┌─────────────────┐       WireGuard Tüneli       ┌────────────────────────┐
│  Mobil Telefon  │ ───────────────────────────> │ Evdeki Linux PC (Arch) │
│  (Termius/SSH)  │       (Tailscale VPN)        │  sshd + opencode       │
└─────────────────┘                              └────────────────────────┘
```

- **Ne Çözüldü?** Port açmadan, NAT arkasında olsak bile Tailscale üzerinden eve güvenli SSH bağlantısı sağlandı. Tek bir systemd birimi (`gokkusagi.service`) hem `sshd` hem de `tailscaled` servislerini koordine etti.
- **Nerede Tıkandı?** Telefondan ufak bir ekranda SSH terminali açmak, sanal klavyeyle komut yazmak ve terminal tabanlı yapay zeka ajanlarını (`opencode`, `neovim`) mobilde yönetmeye çalışmak pratik değildi. Yolda yürürken telefon terminalinde kod debug edilemiyordu.
- **Doğan Soru:** *"Neden terminale yazmak zorundayım? Neden cebimdeki mesajlaşma uygulamasından bilgisayarıma doğal dille veya tek kelimelik komutla emir veremiyorum?"*

---

## Bölüm 2: Asistan Denemeleri ve İlk Duvara Çarpmalar (2025 Yaz & Güz)

2025 yılına gelindiğinde yerel modeller ve kodlama ajanları hızla sahneye çıkmaya başladı. Ben de bilgisayarımı yönetecek bir yazılım katmanı inşa etmeye koyuldum.

### 2.1 İlk Yerel Kodlama Ajanı: `ai-code-agent` (Ağustos 2025)
Localhost'ta çalışan Ollama ile konuşan ilk otonom denememdi.
- Mimari: `src/ollama/` istemcisi, `src/fs/` dosya operasyonları ve `src/utils/execute_cmd.py` komut motoru.
- Amaç: Modele bir görev vermek ve `./AGENT/` klasörü içinde tam yetkili olarak terminal komutları çalıştırıp dosyalar üreterek görevi bitirmesini sağlamak.
- Ders: Modelin terminal çıktısını okuyup bir sonraki adımı planlaması büyüleyiciydi; fakat sistem tamamen masaüstüne hapsolmuştu, dış dünyadan bir tetikleyici alamıyordu.

### 2.2 Donanım ve Tarayıcı Kumandası: `deviceAssistant` (Eylül 2025)
Hemen ardından bilgisayarın tarayıcısını ve konsolunu kontrol etmek için `deviceAssistant` doğdu:
- `src/browser`: Selenium ve Chrome WebDriver ile web otomasyonu.
- `src/consoler`: Terminal logları ve süreç denetimi.
- `src/ai`: Karar mekanizması.

### 2.3 Bulut Çıkmazı ve İlk WhatsApp Vizyonu: `assistant-bot-backend` (Eylül 2025)
Eylül 2025'te radikal bir deneme yaptım: *"Madem dışarıdan erişmek istiyorum, sunucuyu Google Cloud Platform (GCP) Free Tier'a taşıyayım."*  
`assistant-bot-backend` projesinin `README.md` dosyasına o gün aynen şu cümleyi not düşmüştüm:
> **"WhatsApp entegrasyonu → resmi WhatsApp Cloud API kullan."**

Fakat bu yaklaşım iki büyük duvara tosladı:
1. **Bulut Kısıtı:** GCP e2-micro sunucunun 0.6 GB RAM'i vardı. Ne yerel model çalıştırabiliyordu ne de evimdeki Linux makinenin GPU'suna ve dosyalarına erişebiliyordu.
2. **Meta / WhatsApp Cloud API Bürokrasisi:** Meta Developer hesabı açmak, işletme doğrulaması beklemek, template mesaj onayları almak ve webhook yönlendirmek kişisel, hızlı ve bağımsız bir geliştirici etiğine tamamen aykırıydı.

Çözüm bulutta değil; **evdeki makinenin içinde çalışan, yerel bir WhatsApp köprüsündeydi.**

---

## Bölüm 3: CHIMERA v1 (Şubat 2026) — Büyük Düş, Büyük Mimari ve Büyük Tıkanma

Şubat 2026'ya gelindiğinde tüm bu parçaları birleştirecek nihai projenin adı kondu: **CHIMERA**.

```text
                               ┌───────────────────────────────────┐
                               │       KULLANICI ARAYÜZÜ           │
                               │   (WebUI :8787 / WebSocket)      │
                               └─────────────────┬─────────────────┘
                                                 │
                                                 ▼
     ┌───────────────────────────────────────────────────────────────────────────────────────────┐
     │                       THE BRAIN: Orkestratör & Motor (GOLANG)                             │
     │  - internal/agent/engine.go (ReAct Döngüsü & Hata Yakalama)                               │
     │  - internal/server/websocket.go (Canlı Olay Dağıtımı & Log Sink)                          │
     └───────────────────────────────┬───────────────────────────┬───────────────────────────────┘
                                     │                           │
                                     ▼                           ▼
        ┌────────────────────────────────────────┐   ┌────────────────────────────────────────┐
        │      THE HANDS: Sistem & CLI (PYTHON)  │   │      THE EYES: Web Tarama (NODE.JS)    │
        │  - internal/skills/runCmd.py (Bash CLI)│   │  - runtimes/node-v25.4.0 (Playwright)  │
        │  - internal/skills/readFile.py (Dosya) │   │  - Modern JS Sayfa Kazıma              │
        └────────────────────────────────────────┘   └────────────────────────────────────────┘
```

### 3.1 Model Vekili: `ollama_bridge` (6 Şubat 2026)
CHIMERA'dan iki gün önce, masaüstünde `ollama_bridge` adında tek dosyalık bir Python ters vekili (proxy) yazdım. Port 11435'i dinliyor, gelen istekleri `OLLAMA_BASE`'e iletiyor, araya girip model parametrelerini (`HEDEF_MODEL = "gemma2:9b"`) dinamik olarak eziyordu. Modellerin arasına bir proxy koyup akışı kontrol etme yeteneği kazanılmıştı.

### 3.2 Üç Parçalı (Polyglot) Mimari
8 Şubat 2026'da CHIMERA v1 resmen ayağa kalktı:
- **Go Çekirdeği:** Sistemin ana sunucusu, bellek yönetimi ve WebSocket omurgası Go ile yazıldı. Modüler mimari ilk kez bu ölçekte uygulandı.
- **Python Çalışma Zamanı (`runtimes/python-3.14.2`):** `tools.json` üzerinden tetiklenen `runCmd.py` ile `/bin/bash` çalıştırıldı; stdout, stderr, exit code ve milisaniye bazında çalışma süresi toplandı.
- **Node.js Çalışma Zamanı (`runtimes/node-v25.4.0`):** Puppeteer ve Playwright ile interneti kazıyacak gözler hazırlandı.

### 3.3 Kendi Kendini Denetleyen ReAct Döngüsü (`engine.go`)
CHIMERA standart bir sohbet botu değildi. `internal/agent/engine.go` içinde modelden şu JSON sözleşmesi istendi:
```json
{
  "thought": "Kullanıcının isteğini yerine getirmek için diskteki dosyaları listelemeliyim.",
  "tool": "runcmd",
  "args": {"command": "ls -la data/logs"},
  "reply": ""
}
```
Eğer yerel model bozuk bir JSON üretirse, Go motoru `InvalidJSONError` fırlatıyor ve modele otomatik bir düzeltme promptuyla geri dönüyordu:
> *"Your previous response was not valid JSON. Return ONLY valid JSON with keys: thought, tool, args, reply."*

Araç çıktısı geldiğinde ise döngü devam ediyordu:
> *"Tool 'runcmd' Output: ... Based on this, what is the next step or final answer?"*

### 3.4 Şeffaf Zeka (Glass Box AI) ve Canlı WebSocket
Ajanın ne düşündüğünü canlı görmek için `EventThought`, `EventToolCall`, `EventToolResult` ve `EventReply` olayları WebSocket üzerinden `localhost:8787` arayüzüne basıldı. Ajan düşünürken ekranda adım adım düşünce zinciri akıyordu.

### 3.5 8 Şubat 2026 Günlüğü: O Tarihi Kehanet
CHIMERA'nın `docs/2-Günlük.md` dosyasına 8 Şubat 2026 gecesi düştüğüm not, geleceğin haritasıydı:
> **"buraya ileride browser, telegram, whatsapp, instagram, searchEngine falan gelicek. cmd ile pc'yi zaten kontrol eder istediği gibi de neyse..."**

Ve `configs/config.yaml` dosyasında bir Telegram yetkilendirme bloğu açılmıştı:
```yaml
# telegram:
#   enabled: true
#   token: "ENV_TELEGRAM_TOKEN"
#   allowed_users: [123456789]
```

### 3.6 Neden Tıkandı? Yerel Modellerin Format Krizi
CHIMERA v1 harika bir vizyondu fakat **zamansız bir doğumdu**.
O gün kullanılan yerel modeller (`qwen3-vl:8b`, `gemma2:9b`), Go içinde kurduğum katı ReAct döngüsünü sürdürmekte zorlanıyordu:
- Model birkaç adım sonra JSON formatını bozuyor, düzeltme istemine rağmen saçmalıyordu.
- Ajan döngüsünü (ReAct Loop) sıfırdan Go koduyla yönetmeye çalışmak devasa bir bakım yükü getiriyordu.
- Proje WebUI ekranına kilitli kaldı; dış dünya köprüsü (WhatsApp/Telegram) henüz kodlanmamıştı.

CHIMERA v1, mimari bir başyapıt olarak donduruldu.

---

## Bölüm 4: Arayış ve Uydu Projeler Dönemi (Haziran – Ekim 2026)

CHIMERA v1'in duraklaması bir vazgeçiş değil; **eksik organları ayrı laboratuvarlarda tek tek üretme dönemine** dönüştü.

### 4.1 CHIMERA2 (v2): Minimalist Geri Çekilme (Haziran 2026)
15 Haziran 2026'da `aider` kullanarak `CHIMERA2` adıyla projeyi sıfırdan, minimal bir Go iskeletiyle yeniden başlatmayı denedim (`main.go`: *"Merhaba Go dünyası! Ben geri döndüm"*). Ancak sorunun dilde değil, dış dünya bağlantısında ve zeka katmanında olduğunu görerek beklemeye aldım.

### 4.2 Ses Organı: `test_ASR` ve Türkçe Whisper Benchmarkları (Ekim 2026)
Bir asistan sesli notları dinleyemiyorsa eksiktir. Ekim 2026'da `test_ASR` projesini başlattım:
- `faster-whisper` (CTranslate2) motoru kuruldu.
- Türkçe için `whisper-large-v3-turbo-tr-finetuned` modeli optimize edildi.
- Telefon gürültülerini temizlemek için `DeepFilterNet3` entegre edildi.
- VRAM taşması durumunda otomatik olarak CPU int8'e düşen dirençli bir boru hattı hazırlandı.

### 4.3 Komuta Organı: `whatsapp-cmd` ve Shell Tünelleme (Ekim 2026)
5 Ekim 2026'da WhatsApp mesajlarını doğrudan bir klasördeki dosyalara (`Desktop/whatsapp-cmd/`) yazıp okuyan ve bash komutları çalıştıran ilkel bir tünel denendi (`TESLIM-SORUNU-2026-10-05.md`). Dosya tabanlı bu mekanizma, mesajlaşma üzerinden sistem yönetilebileceğini canlı olarak kanıtladı.

---

## Bölüm 5: Kırılma Anı (Ekim – Kasım 2026) — İki Parçalı Köprünün İnşası

Ekim ayı sonunda ve Kasım 2026 başında taşlar inanılmaz bir hızla yerine oturmaya başladı. Sistem iki bağımsız depoda hayata geçti:

### 5.1 Katman 1 & 2: `whatsapp-bridge-mcp` (Go & FastMCP)
- **Whatsmeow Go Daemon:** Meta Cloud API saçmalığı tamamen çöpe atıldı. Whatsmeow kütüphanesiyle saf Go dilinde, QR kodla bir kez eşleşen, SQLite'ta oturum saklayan (`session.db`), `:31415` portunda REST ve gerçek zamanlı SSE (`/events`) yayını yapan mükemmel bir yerel köprü kuruldu.
- **FastMCP STDIO Sunucusu:** AI ajanlarının WhatsApp'a mesaj, reaksiyon, anket ve konum gönderebilmesi için Model Context Protocol (MCP) araç seti yazıldı.

### 5.2 Katman 3 & 4: `whatsapp-tetikleyici` (Dispatcher & Multimodal)
- **Event Dispatcher:** Köprüden gelen SSE akışını dinleyen, telefon ve WhatsApp LID bazlı katı bir RBAC/whitelist uygulayan Go motoru yazıldı.
- **Multimodal Entegrasyonu:**
  - Ses geldiğinde: Doğrudan `test_ASR` projesini referans alan `voice_pipeline.py` devreye girdi.
  - Fotoğraf geldiğinde: Ollama Qwen Vision modellerini kullanan `gorsel_pipeline.py` çalıştı.
- **Sistem Servisi:** `systemd --user` ile bilgisayar açıldığı anda arka planda otomatik başlayan kesintisiz bir yapı kuruldu.

### 5.3 Aydınlanma Anı: "Beyni Sıfırdan Yazmak Zorunda Değilim!"
İşte bu projenin en büyük mühendislik kırılma anı (Aha! Moment) burada yaşandı:

> **CHIMERA v1'de aylarca yerel modelin JSON çıktısını düzeltmekle, ReAct döngüsünü Go'da elle kurmakla uğraşmıştım. Oysa 2026 sonu itibarıyla sektörde bunu mükemmel yapan devrimsel CLI araçları vardı: Antigravity CLI (`agy`), Claude Code (`claude`) ve OpenCode (`opencode`).**

Benim artık bir LLM istemcisi veya ReAct motoru yazmama gerek yoktu!  
Tetikleyicinin yapması gereken tek şey; WhatsApp'tan gelen doğrulanmış mesajı, sistemdeki bu güçlü CLI ajanlarının kalıcı oturumlarına tünellemek (`whatsapp-to-agent`) ve sistem komutlarını doğrudan çalıştırmaktı (`whatsapp-to-cmd`).

CHIMERA'nın felç olan beyni, modern CLI ajanlarıyla bir anda süper zekaya kavuştu!

---

## Bölüm 6: CHIMERA v3 (Kasım 2026) — Büyük Birleşme ve Çok Kanallı Gelecek

Bugün geldiğimiz nokta, iki ayrı parçaya bölünmüş olan `whatsapp-bridge-mcp` ve `whatsapp-tetikleyici` depolarının tek bir devasa yapıda birleştiği andır: **CHIMERA v3**.

### 6.1 Monorepo Dönüşümü (`go.work` ve Ortak `.venv`)
Ayrık repoların yarattığı bağımlılık karmaşası, çift test koşuları ve senkronizasyon zorlukları; `go.work` çalışma alanı ve tek bir Python sanal ortamı (`.venv`) altında birleşiyor. 627+ birim test tek bir `make verify` komutuyla yeşile dönüyor.

### 6.2 WhatsApp Bir Varış Noktası Değil, Sadece İlk Kapı
Ve en büyük felsefi idrak:  
**Bu proje asla sadece bir "WhatsApp projesi" olmadı.**  
WhatsApp, aşılması gereken en zorlu, en kapalı ve en çetrefilli protokoldü. O fethedildi.

### 6.3 Evrensel Çok Kanallı Ağ Geçidi (Omnichannel Gateway) Mimarisi
CHIMERA v3 ile birlikte çekirdek motor (Tetikleyici / Dispatcher), mesajın nereden geldiğinden tamamen soyutlanıyor:

```text
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                               CHIMERA v3: OMNICHANNEL GATEWAY                               │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
          │                                 │                                 │
          ▼                                 ▼                                 ▼
┌──────────────────┐              ┌──────────────────┐              ┌──────────────────┐
│ WHATSAPP BRIDGE  │              │ TELEGRAM BRIDGE  │              │  DISCORD / SLACK │
│ (Whatsmeow / Go) │              │  (TDLib / Bot)   │              │     BRIDGES      │
└─────────┬────────┘              └─────────┬────────┘              └─────────┬────────┘
          │ :31415                          │ :31416                          │ :31417
          └────────────────────────┬────────┴─────────────────────────────────┘
                                   │ SSE / REST (Ortak Event Formatı)
                                   ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                       CHIMERA ÇEKİRDEK MOTORU (TETİKLEYİCİ / DISPATCHER)                    │
│  - Whitelist & RBAC Doğrulaması (Telefon, LID, Telegram UserID)                             │
│  - Hızlı Tepki Emojileri (👀 -> ✅ / ❌), Anketler, Durum Makinesi                           │
│  - Güvenli İzin & Onay Sistemi (.tetikleyici/izinler.json)                                  │
└──────────────────────────────┬───────────────────────────────┬──────────────────────────────┘
                               │                               │
                               ▼                               ▼
     ┌───────────────────────────────────┐   ┌───────────────────────────────────┐
     │       MULTIMODAL BORU HATLARI     │   │         AJAN VE KOMUTA TÜNELİ     │
     │  - voice_pipeline (faster-whisper)│   │  - whatsapp-to-cmd (Shell Tüneli) │
     │  - gorsel_pipeline (Ollama Qwen)  │   │  - Ajan Adaptörleri (Driver Katmanı)│
     │  - Piper TTS (Yerel Sesli Yanıt)  │   │    [ agy • claude • opencode ]    │
     └───────────────────────────────────┘   └───────────────────────────────────┘
```

Yarın Telegram köprüsü eklendiğinde; arkadaki transkripsiyon hattı, görsel zeka, güvenlik whitelist'i ve Claude/Antigravity ajan oturumları **tek bir satır bile değişmeden** Telegram kullanıcılarına da hizmet verecek.

---

## Bölüm 7: Zaman Çizelgesi ve Karşılaştırma Matrisi

### 7.1 Kronolojik Karşılaştırma Matrisi

| Yetenek / Katman | CHIMERA v1 (Şubat 2026) | Ara Dönem & Uydular (Yaz 2026) | CHIMERA v3 (Kasım 2026) |
| :--- | :--- | :--- | :--- |
| **Dış İletişim Kanalı** | Yok (Sadece yerel WebUI :8787) | SSH / Tailscale (`gokkusagi`) | **WhatsApp Web** + Gelecekte **Telegram / Discord** |
| **Orkestrasyon Motoru** | Go (Sıfırdan ReAct döngüsü) | Bash scriptleri (`systemd`) | **Go Daemon & Event Dispatcher (SSE + REST)** |
| **Yapay Zeka Beyni** | Saf yerel Ollama (8B modeller) | İzole CLI testleri | **Üst Düzey CLI Ajanları (`agy`, `claude`, `opencode`)** |
| **Cihaz / Terminal Kontrolü** | `internal/skills/runCmd.py` | Termux SSH (`gokkusagi`) | **`whatsapp-to-cmd` + Shell Hub + Masaüstü `/ekran`** |
| **Ses İşleme (Kulaklar)** | Yok | `test_ASR` (faster-whisper) | **`voice_pipeline.py` (Whisper Large Turbo + CPU Fallback)** |
| **Görsel İşleme (Gözler)** | Node.js Playwright kazıma | Temel Ollama testleri | **`gorsel_pipeline.py` (Ollama Qwen-VL + OCR Katmanı)** |
| **Güvenlik & Yetkilendirme** | YAML içinde taslak ID listesi | SSH Anahtarları | **LID / Telefon RBAC + WhatsApp Anketli İnteraktif Onay** |
| **Mimari Yapı** | Monolitik Go + Runtimes | Dağınık Repolar | **Birleşik Monorepo (`go.work` + Gömülü Dashboard)** |

---

## Bölüm 8: Çıkarılan Mühendislik Dersleri

1. **"Doğru Fikir, Yanlış Zaman" Yanılgısı:**  
   CHIMERA v1'de kurduğum hayal tamamen doğruydu; fakat o günün 8B yerel modelleri o mimariyi sırtlayacak kapasitede değildi. Bazen bir projeyi zorlamak yerine uykuya yatırmak ve ekosistemin olgunlaşmasını beklemek en doğru mühendislik kararıdır.
2. **Tekerleği Yeniden İcat Etmemek:**  
   ReAct döngüsünü, düşünme zincirini ve tool çağırma mekanizmasını Go içinde sıfırdan yazmak yerine, bu alanda devasa bütçelerle geliştirilen CLI ajanlarını (`agy`, `claude`) birer adaptörle sisteme bağlamak projeyi aylar sürecek krizlerden kurtardı.
3. **Local-First ve Bağımsızlık İlkesi:**  
   Meta'nın bulut API'lerine veya üçüncü parti sunuculara bel bağlamamak; doğrudan protokolü yerel makinede koşturmak (`Whatsmeow`) sisteme sınırsız bir mahremiyet, sıfır maliyet ve tam donanım erişimi sağladı.
4. **Çok Dilli (Polyglot) Dengesi:**  
   Her şeyi Go ile yazamazsınız (yapay zeka kütüphaneleri Python'dadır). Her şeyi Python ile de yazamazsınız (7/24 çökmeyen kararlı bir daemon Go ile çok daha kolaydır). Çözüm; Go'yu omurga, Python'u ise uzman boru hatları olarak kullanmaktır.

---

## Son Söz: Ejderha Uyanıyor

CHIMERA; 2024'te küçük bir SSH tüneli olarak doğdu, 2025'te bulut duvarlarına çarptı, 2026 başında yerel modellerin acemiliğine takıldı. 

Fakat bu süreçte hiçbir emek kaybolmadı:
- `test_ASR`'ın kulakları,
- `ollama_bridge`'in gözleri,
- `gokkusagi`'nın tünel disiplini,
- `runCmd.py`'nin terminal cesareti,
- Ve `Whatsmeow`'un yerel iletişim gücü...

Hepsi bir araya geldi ve mitolojideki o görkemli varlık nihayet ayağa kalktı.  
Bugün inşa ettiğimiz sistem artık basit bir köprü değil; **cebimizdeki mesajlaşma penceresinden kendi bilgisayarımızın zekasına ve terminaline hükmettiğimiz evrensel bir ağ geçididir.**

**Nerve Dispatch (NerveDispatch), CHIMERA vizyonunun ve mimarisinin yaşayan, çok kanallı nihai formudur.**

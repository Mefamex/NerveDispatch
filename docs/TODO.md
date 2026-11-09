# TODO - WhatsApp Bridge + Tetikleyici

> **Son güncelleme: 2026-11-08 04:23 +03**
>
> Not: bu dosya yalnız **açık işleri** listeler. Biten her şey
> `TODO.DONE.md`'ye taşınır.
> Köprü: `paired:true connected:true` (systemd) · Tetikleyici: PID **242809**
> - restart yalnız `./tetikleyici.sh restart` · **1 kural** canlı ✅
> Test: **627+**, `go vet`/`gofmt` temiz

---

## 🚀 WhatsApp Yerel Özellikleri Sprinti (Sırayla Yapılacaklar)

> **🚨 STRATEJİK KARAR VE YENİ SOHBET DEVİR-TESLİM NOTU (2026-11-09 - NerveDispatch):**
> - **Hedef Konum:** `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/`
> - **Master Mimari ve Geçiş Planı:** [`docs/YENI-MIMARI-VE-BIRLESME-PLANI.md`](./YENI-MIMARI-VE-BIRLESME-PLANI.md) 
> - **Evrim Hikayesi:** [`docs/NERVE-DISPATCH-THE-STORY.md`](./NERVE-DISPATCH-THE-STORY.md)
> - **Sıradaki İlk Adım (FAZ 1):** Eski iki repoda (`whatsapp-bridge-mcp` ve `whatsapp-tetikleyici`) yeni özellik kodlanmayacak; doğrudan `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` içine birleştirme yapılacak (**Madde 22**).
> - **FAZ 2:** `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` üzerinde **önce WhatsApp hattı ve tüm çekirdek yetenekler** (**Madde 7'den itibaren**) sırayla tamamlanacak.
> - **FAZ 3:** Çekirdek tamamlandıktan sonra **Telegram (`975678139`), Discord ve Instagram** kanalları eklenecek.
>
> *(Not: Tamamlanan Sprint #1–#6 maddeleri ve arşivlenmiş işler [TODO.DONE.md](./TODO.DONE.md) kütüğüne taşınmıştır.)*

### 7. 🧵 Tek Kalıcı Oturum (Sabit Session ID) ve Ajan Adaptör Mimarisi [TAMAMLANDI ✅]
- [x] **Ortak `AgentDriver` Arayüzü (`tetikleyici/drivers/`):** `Driver`, `SessionStore`, `AgyDriver`, `ClaudeDriver`, `OpencodeDriver` implementasyonu tamamlandı.
- [x] **Kalıcı Oturum Deposu:** `.tetikleyici/oturumlar.json` üzerinden her motorun sabit konuşma ID'si okundu/yazıldı.
- [x] **Eşzamanlılık Kilidi:** Oturum ve motor bazlı `oturumKilidiAl` kilit mekanizması bağlandı.
- [x] **Doğrulama:** Birim testleri ve `make verify` yeşil geçti, canlı sistemd servisi devreye alındı.
**Sorun:**
1. `runAgent` (`tetikleyici/action.go`) agy ve claude için her mesajda `agy -p "<prompt>"` / `claude -p "<prompt>"` çalıştırıyor → her mesaj SIFIR bağlamlı yeni konuşma; önceki mesajı hatırlamıyor.
2. `opencode`, `agy` ve `claude` CLI'larının oturum listeleme (`session list`), oturum formatları (`ses_...` vs UUID), bayrakları ve girdi/çıktı yöntemleri tamamen birbirinden farklı. Tek bir fonksiyon içine gömülen `if/else` blokları sürdürülemez hâle geliyor.

**Ölçülen bayraklar ve farklar (2026-11-08):**
- **agy:** `--conversation <ID>` (belirli konuşmayı sürdür) · `-c/--continue` (en sonuncuyu sürdür) · oturumlar `~/.gemini/antigravity-cli/` altında tutulur.
- **claude:** `-r/--resume <session-id>` · `--session-id <uuid>` (ilk oturumu bu kimlikle kur) · `-c` · oturum listesi `~/.claude/projects/` altındadır.
- **opencode:** `-s/--session <id>` (yoksa oluşturur) · `opencode session list` komutu ile listelenir.

**Tasarım kararı: Ajan Adaptör (Driver/Provider) Katmanı:**
- [ ] **Ortak `AgentDriver` Arayüzü:**
  - `ListSessions()` → `[{id, title, updated_at}]`
  - `Run(session_id, prompt)` → string reply
  - `NewSession(prompt)` → `(session_id, reply)`
- [ ] **Ajan Sürücüleri (`drivers/` veya `agents/`):**
  - `opencode_driver.go`: `opencode session list` parse eder, `opencode run -s <id>` çağırır.
  - `agy_driver.go`: `~/.gemini/antigravity-cli/cache/` ve veritabanından oturumları çeker, `--conversation <id>` ile mesaj iletir.
  - `claude_driver.go`: JSONL kayıtlarından oturumları listeler, `--resume <id>` / `--session-id <uuid>` yönetir.
- [ ] **Sanal Takma Ad (Alias / Sabit Kimlik):**
  - `tetikleyici.yaml` içinde `session: "whatsapp-mefamex"` gibi sanal alias desteği.
  - Adaptör ilgili ajanda bu alias'a denk gelen oturumu yoksa bir kere açıp ID'sini `.tetikleyici/oturumlar.json`'a sabitler; sonraki tüm mesajlar o ID üzerinden akar.
- [ ] **Eşzamanlılık ve Sıralama:** `oturumKilidiAl` mekanizması tüm ajan sürücülerine genişletilecek; aynı anda gelen WhatsApp mesajları aynı oturumu kilitleyip sırayla işleyecek.
- [ ] **Bağlam Şişmesi Koruması:** Tek oturum sonsuz uzamaması için token/tur eşiği aşıldığında otomatik arşivleme/özetleme veya `/session yeni` komutu.
- [ ] **Prompt Temizliği:** Kalıcı oturumda geçmiş mesaj penceresi tekrar tekrar prompt'a gömülmeyecek; yalnız yeni gelen mesaj gönderilecek.

### 8. 🧹 Eski Monorepo'dan Taşınan Dosyaların PII Taraması ve Pushlanması [BEKLİYOR - sonra]
Silinen eski `whatsapp-bridge-project` arşivinden kopyalanıp **commit'lenmemiş** duran tek nüsha dosyalar (orijinal silindi, bu kopyalar kaybolmamalı):
- [ ] `whatsapp-tetikleyici/docs/` - eski tasarım/plan/sorun belgeleri (2026-10-04…06). **Gerçek telefon/LID numarası içerebilir**: `905551234567` / `12345678901234@lid` kuralına göre temizle, sonra commit.
- [ ] `whatsapp-tetikleyici/tetikleyici.sh` - yönetim betiği. Artık systemd user servisi var; betik hâlâ gerekli mi karar ver (gerekiyorsa servisle uyumlu hâle getir, yoksa sil).
- [ ] `whatsapp-bridge-mcp/sidecar/` (VoIP sidecar) - **içinde `node_modules/` var, commit'LEME**. `.gitignore`'a `sidecar/node_modules/` ekle; `rsrc-dump.txt` (322 KB) gereksiz olabilir, gözden geçir; `package-lock.json` kalsın. Sidecar ayrı repo mu olmalı, yoksa köprü reposunda mı kalmalı karar ver. Köprü şu an `WASM binary not found` uyarısı veriyor (arama özelliği kapalı).
- [ ] Taramadan sonra: `grep -rnE "9054|9053|1234563"` ile kalan gerçek kimlik kontrolü, `make test`, commit, push.

### 9. 🛡️ İnteraktif Onay ve İzin Sistemi (Interactive Approval / Permissions) [PLANLANDI]
- [ ] **Tehlikeli Eylem Koruması:** Ajan veya tetikleyici riskli komutlar (`rm`, `git push --force`, `systemctl`, veri silme vb.) çalıştırmadan önce işlemi durdurur ve admin onayına sunar.
- [ ] **WhatsApp Anket (Poll) ile Onay:** Köprünün native anket (`whatsapp_send_poll`) yeteneğiyle 4 seçenekli soru gönderilir:
  1. `[1] Tek seferlik onayla`
  2. `[2] Reddet`
  3. `[3] 1 saat boyunca tümünü otomatik onayla`
  4. `[4] 1 hafta boyunca otomatik onayla`
- [ ] **Yetki Damgası Saklama:** `.tetikleyici/izinler.json` dosyasında geçici onay süreleri (`auto_approve_until: <timestamp>`) tutulur; süre bitiminde tekrar onay sorma moduna dönülür.

### 10. 📸 `/ekran` Masaüstü Ekran Görüntüsü Alma Komutu [PLANLANDI]
- [ ] **Masaüstü Görüntüsü Gönderme:** Tıpkı `/camera` gibi çalışan ama monitörün anlık ekran görüntüsünü WhatsApp'a yollayan `/ekran` (veya `/screenshot`) komutu.
- [ ] **Wayland / KDE Desteği:** Arch Linux KDE Plasma Wayland üzerinde `spectacle -b -n -o /tmp/ekran.png` (veya `grim`) çağrısı ile ekran yakalanır, `whatsapp-bridge` üzerinden resim olarak iletilir.
- [ ] **Opsiyonel Parametreler:** `/ekran 5` (5 saniye gecikmeli çekim) desteği.

### 11. ⏱️ Yapay Zekâ Destekli Gelişmiş Zamanlayıcı (Agentic Cron / Scheduled Workflows) [PLANLANDI]
- [ ] **Sabit Metinden Akıllı İş Akışına Geçiş:** Mevcut `/zamanla` yalnız statik metin gönderiyor. AGY / Claude Desktop benzeri "Ajanlı Zamanlanmış İş Akışları" yapısına kavuşturulacak.
- [ ] **Örnek Komutlar:**
  - `/zamanla her sabah 08:30 takvimimi, maillerimi ve GitHub bildirimlerimi özetle`
  - `/zamanla 2 saatte bir sunucu ve servis sağlığını kontrol et, sorun varsa haber ver`
- [ ] **Arka Plan Çalıştırıcısı:** Zamanı geldiğinde ajanın sabit oturumu tetiklenir, araçları (MCP/CLI) kullanarak araştırmasını yapar ve üretilen çıktıyı WhatsApp'a iletir.

### 12. ⏳ Canlı İlerleme Bildirimi ve Hızlı Tepki Emojisi (Emoji Feedback) [PLANLANDI]
- [ ] **Yapılandırma Dosyası Desteği (`tepki_emojisi`):**
  - `tetikleyici.yaml` içinde `tepki_emojisi` bloğu:
    ```yaml
    tepki_emojisi:
      aktif: true
      isleniyor: "👀" # veya "⏳"
      basarili: "✅"
      hatali: "❌"
    ```
  - Sadece eşleşen kurallara ve yetkili admin mesajlarına uygulanır (tüm genel WhatsApp gruplarını/sohbetlerini kirletmez).
- [ ] **Gelen Mesaja Anlık Reaksiyon:**
  - Admin mesaj attığında tetikleyici mesajı alıp sıraya koyduğu an mesaja otomatik `isleniyor` ("👀" veya "⏳") reaksiyonu bırakır.
  - İşlem başarıyla bitip cevap gönderildiğinde reaksiyon `basarili` ("✅") yapılır.
  - Hata oluşursa reaksiyon `hatali` ("❌") yapılır.
  - Kullanıcı emojiyi config'ten dilediği gibi değiştirebilir ya da `aktif: false` ile tamamen kapatabilir.
- [ ] **Uzun Süren Görevlerde Durum Güncellemesi:** 10 saniyeden uzun süren analizlerde WhatsApp'a ara aşama durum bildirimleri (*"Araştırılıyor 🔍"*, *"Derleniyor ⚙️"*) gönderilir.

### 13. 🏛️ KRİTİK MİMARİ YÜKSELTMESİ: Headless Daemon + Birleşik CLI + Gömülü Web Dashboard [PLANLANDI]
> **Stratejik Tasarım Kararı:** Ağır masaüstü pencere uygulamaları (Electron, Qt, Tauri vb.) yerine; sistem yükünü minimumda tutan (~30MB RAM), headless sunucu/VPS/Raspberry Pi uyumlu, güvenlik açığı yaratmayan ve üst düzey kullanım kolaylığı sunan 3 katmanlı hibrit mimariye geçiş.

- [ ] **Katman 1: Headless Daemon (Mevcut Hafif Çekirdek Korunur):**
  - Go + Whatsmeow köprüsü ve Tetikleyici motoru arka planda `systemd` servisi olarak bağımsız çalışmaya devam eder.
  - X11/Wayland zorunluluğu yoktur; ekransız sunucularda sıfır konfigürasyonla çalışır.
- [ ] **Katman 2: Birleşik Yönetim CLI (`tetikleyici` / `wb` CLI):**
  - Dağınık bash scriptleri yerine tek bir derlenmiş, zengin Go CLI arayüzü:
    - `tetikleyici status`: Köprü, tetikleyici ve AI ajanlarının canlılık/gecikme durumunu gösterir.
    - `tetikleyici qr`: Terminal içinde ANSI bloklarıyla WhatsApp eşleme QR kodu üretir.
    - `tetikleyici logs -f`: Anlık SSE olaylarını ve tünel akışını terminale renkli formatta basar.
    - `tetikleyici config check`: YAML yapılandırmasını, yetkili numaraları ve kuralları doğrular.
    - `tetikleyici session list/kill`: Aktif oturumları CLI'dan yönetir.
- [ ] **Katman 3: Sıfır-Yüklü Gömülü Web Paneli (Go `embed` Web Dashboard):**
  - Ekstra Node.js veya Electron gerektirmeyen, doğrudan Go ikilisine gömülü (`go:embed`) tek sayfalık hafif web arayüzü (`http://127.0.0.1:31415/dashboard`).
  - Yalnızca yerel arayüzü (`127.0.0.1`) dinler, harici saldırı yüzeyi oluşturmaz.
  - **Özellikler:**
    - Tarayıcı üzerinden tek tıkla büyük ve net QR kod okutma (özellikle ilk kurulum ve mobil tarayıcı kolaylığı).
    - Canlı servis durum kartları (bellek, CPU, bağlantı durumu, uptime).
    - Canlı SSE mesaj ve olay akışı paneli.
    - Yetkili kullanıcı ve kural izinlerinin görsel özeti.
- [ ] **Katman 4: Asıl Kullanıcı Arayüzü WhatsApp'ın Kendisidir ("ChatOps"):**
  - Günlük kullanımda masaüstünde ek pencereye ihtiyaç duyulmaz; `/cmd`, `/camera`, `/ekran`, `/agent`, anketler ve emoji reaksiyonlarıyla WhatsApp içi yönetim korunur.

### 14. 📁 Çift Yönlü Dosya ve Belge Transferi (`/get <yol>` ve WhatsApp Dosya Kabulü) [PLANLANDI]
- [ ] **WhatsApp'tan Telefona Dosya İndirme (`/get <dosya_yolu>`):**
  - Admin WhatsApp'tan `/get /mnt/DiskD/rapor.pdf` veya `/get config.yaml` yazdığında dosya varlığı ve boyutu (max 50 MB WhatsApp document sınırı) doğrulanır.
  - MIME türü otomatik tespit edilir (`application/pdf`, `text/plain`, `image/*` vb.).
  - `whatsapp-bridge` üzerinden `POST /send-file` çağrılarak dosya doğrudan WhatsApp belgesi (`document`) formatında sohbete teslim edilir.
- [ ] **Telefondan Bilgisayara Dosya Yükleme (Inbound File Handler):**
  - Admin telefondan bir dosya/script attığında güvenli varsayılan dizine (`~/Downloads/whatsapp-files/` veya aktif ajanın çalışma dizinine) indirilir.
  - Güvenlik: İndirilen dosyaların çalıştırma yetkisi kaldırılır (`chmod -x`).
  - Bildirim: WhatsApp'a *"📥 Dosya kaydedildi: `~/Downloads/whatsapp-files/rapor.pdf` (1.4 MB)"* yanıtı verilir.
  - Oturum tünelindeyse (`/session go`), ajanın bağlamına *"Kullanıcı bu dosyayı diske kaydetti: `<tam_yol>`"* bildirimi düşülür.

### 15. 📂 Ajan Çalışma Dizinini Dinamik Değiştirme (`/root <dizin>`) [PLANLANDI]
- [ ] **Dinamik Repo / Çalışma Alanı Değişimi:**
  - Mevcut `/cmdhub cd` yalnız interaktif bash oturumunu etkiliyor. `/root <dizin>` komutu ise aktif kural veya oturumun çalışma dizinini (`cwd` / `workdir`) değiştirir.
  - Örnek: `/root /mnt/DiskD/SOFTWARE/Github-repos/whatsapp-tetikleyici`.
  - Dizin doğrulanır; yoksa açıklayıcı hata mesajı döner.
  - Sonraki tüm `/cmd` çağrıları ve AI ajan çalıştırmaları (`opencode`, `agy`, `claude`) bu yeni kök dizin bağlamında yürütülür.
  - Parametresiz `/root`: Aktif çalışma dizinini WhatsApp'a raporlar.

### 16. 📋 Wayland / KDE Plasma Sistem Panosu Kontrolü (`/clip`) [PLANLANDI]
- [ ] **Sistem Panosunu Telefona Çekme (`/clip` / `/pano`):**
  - Arch Linux KDE Plasma 6 (Wayland) mimarisine uygun olarak `wl-paste -t text/plain` çağrısıyla bilgisayarın panosundaki metin okunur ve WhatsApp'a mesaj olarak iletilir.
  - Çok uzun panolarda ilk 2000 karakter + kırpma uyarısı uygulanır.
- [ ] **Telefondan Bilgisayar Panosuna Metin Kopyalama (`/clip <metin>`):**
  - WhatsApp'tan gelen metin bilgisayar panosuna basılır: `echo -n "<metin>" | wl-copy`.
  - Telefonda kopyalanan bir token, link veya kod parçasını anında bilgisayarda Ctrl+V yapabilme imkânı sağlar.
  - Yalnızca `allow_cmd: true` olan admin kullanıcılara açıktır.

### 17. 🪟 Açık Pencereler ve Süreç Denetimi (`/apps` & `/kill`) [PLANLANDI]
- [ ] **Açık Pencereleri Listeleme (`/apps` / `/pencereler`):**
  - Wayland / Plasma üzerinde çalışan grafik uygulamalarını veya pencereleri listeler (`qdbus org.kde.KWin /KWin org.kde.KWin.queryWindowInfo` / `kdotool` veya aktif GUI süreçleri).
  - Format: `1. Google Chrome (PID 12450)`, `2. Cursor (PID 3341)`, `3. Konsole (PID 5512)`.
- [ ] **Askıda Kalan Süreci Kapatma (`/kill <uygulama_veya_pid>`):**
  - Yanıt vermeyen uygulamayı uzaktan sonlandırır (`kill -15 <pid>` veya `pkill -f <uygulama>`).
  - Güvenlik: Kritik sistem servisleri (`systemd`, `dbus`, `wayland`, `whatsapp-bridge`) kara listededir; asla kill edilemez.

### 18. 🎙️ Yerel TTS Sesli Yanıt (Voice-to-Voice: Piper TTS ile WhatsApp `.ogg opus` Sesli Notu) [PLANLANDI]
- [ ] **Hafif Yerel TTS (Piper):**
  - CPU'da sıfır VRAM tüketimiyle gerçek zamanlı çalışan yerel `Piper TTS` (Türkçe model: `tr_TR-dfki-medium` vb.) entegrasyonu.
  - GPU VRAM'i Whisper veya LLM ile çakışmaz, tamamen hafif çalışır.
- [ ] **WhatsApp Uyumlu Ses Formatı (`.ogg opus` PTT):**
  - Piper tarafından üretilen `.wav` dosyası `ffmpeg` ile WhatsApp sesli not formatına (`audio/ogg; codecs=opus`) çevrilir.
  - `whatsapp-bridge`'in `POST /send-file` endpoint'ine `is_voice_note: true` iletilerek karşı tarafta yeşil mikrofonlu gerçek "sesli mesaj" olarak teslim edilir.
- [ ] **Tetikleme Mantığı:**
  - Admin sesli not attığında (veya yapılandırmada `sesli_cevap: true` ise / `/tts on`) metin yerine doğrudan sesli not gönderilir.

### 19. 🧠 Ajan Bağlamını Hafifletme ve Sıkıştırma (`/compact` / `/ozet`) [PLANLANDI]
- [ ] **Oturum Bağlamını Küçültme:**
  - Uzun süren WhatsApp tünel konuşmalarında token şişmesini ve yavaşlamayı önlemek için WhatsApp'tan `/compact` komutu verilir.
  - Aktif ajana (`agy`, `claude`, `opencode`) oturumdaki ana kararları ve teknik adımları yoğun bir özet haline getirme talimatı iletilir.
  - Oturum sıfırlanmadan bağlam hafifletilir; gereksiz tur geçmişi temizlenir.
- [ ] **Şişme Uyarısı:** Oturum tur sayısı belirlenen eşiği (ör. 25 tur) aştığında tetikleyici WhatsApp'a `"💡 Bağlam penceresi genişledi, /compact önerilir"` ipucu gönderir.

### 20. 🚨 Proaktif Sistem Sağlık Bekçisi ve Eşik Uyarıları (Passive Watcher) [PLANLANDI]
- [ ] **Periyodik Arka Plan Denetimi (Cron/Timer):**
  - Yalnız komut verildiğinde çalışan `/durum` haricinde; `tetikleyici.yaml` içinde tanımlanan eşikleri periyodik (ör. 10 dakikada bir) izleyen arka plan bekçisi:
    ```yaml
    bekci:
      aktif: true
      kontrol_araligi: 10m
      esikler:
        disk_doluluk_yuzde: 90
        ram_doluluk_yuzde: 95
        cpu_sicaklik_c: 85
        servisler: ["whatsapp-bridge.service", "systemd-networkd.service"]
    ```
- [ ] **Otonom Alarm ve Spam Koruması:**
  - Eşik aşıldığında veya kritik servis düştüğünde WhatsApp'a acil durum uyarısı (`🚨 Sistem Bekçisi Uyarısı: / disk doluluğu %92!`) gönderilir.
  - Aynı arıza için 1 saat içinde tekrar mesaj atılmaz (cooldown koruması).

### 21. 🔗 QR Kodsuz 8 Haneli Eşleştirme Kodu (Pairing Code Desteği) [PLANLANDI]
- [ ] **Whatsmeow Pairing Code Entegrasyonu:**
  - Whatsmeow kütüphanesinin yerleşik `client.PairPhone(phoneNumber)` API'si kullanılarak `POST /pair-phone` endpoint'i eklenir.
  - Terminalde QR kod okutmak yerine WhatsApp sunucusundan 8 haneli kod (ör: `ABCD-1234`) alınır.
  - Kullanıcı telefonunda WhatsApp ▸ Bağlı Cihazlar ▸ "Telefon Numarasıyla Bağla" adımına bu kodu girer.
  - Ekransız/headless sunucu kurulumlarında ve kamera netleme sorunlarında QR kodsuz hızlı eşleşme sağlar.

### 22. 🦁 NerveDispatch - Çok Kanallı Modüler Monorepo Birleşmesi [SIRADA #1 - İLK ADIM 🎯]
> **Proje Adı & Hedef Konum:** `NerveDispatch` (`/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/`)  
> **Master Blueprint:** [`docs/YENI-MIMARI-VE-BIRLESME-PLANI.md`](./YENI-MIMARI-VE-BIRLESME-PLANI.md)  
> **Mimari Hikaye:** [`docs/NERVE-DISPATCH-THE-STORY.md`](./NERVE-DISPATCH-THE-STORY.md)  
> **Vizyon:** `whatsapp-bridge-mcp` ve `whatsapp-tetikleyici` projelerini `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` çatısında birleştirmek; önce WhatsApp hattını ve tüm çekirdek özellikleri kusursuzlaştırmak (FAZ 1 & FAZ 2), ardından aynı sabit AI oturumları ve kural motoru üzerine **Telegram, Discord ve Instagram** kanallarını tak-çıkar adaptörler olarak eklemek (FAZ 3).

- [ ] **1. Sıfır-Kırılmalı Dizin, Modüler `install/` & `upgrade/` ve Paket Mimarisi (`go.work` + Çok Kanallı Çekirdek):**
  ```text
  /mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/
  ├── NERVE-DISPATCH-THE-STORY.md      # Chimera v1 • v2'den Nerve Dispatch'e Mimari Evrim Hikayesi
  ├── YENI-MIMARI-VE-BIRLESME-PLANI.md # Master Blueprint & Geçiş Yol Haritası
  ├── bridge/                          # Go Modülü 1: Çok kanallı köprü (:31415 REST, SSE, WebUI, WhatsApp + channels/)
  ├── tetikleyici/                     # Go Modülü 2: Olay dağıtım motoru, Ajan Sürücüleri (drivers/), CLI ve cmd/*.py
  ├── mcp_server/                      # Python: FastMCP STDIO sunucusu (modül gölgelemesini önler)
  ├── gorsel/                          # Python: Görsel, belge, OCR, video ve sticker paketi
  ├── gorsel_pipeline.py               # Ollama V-LLM görsel analiz hattı
  ├── voice_pipeline.py                # Faster-Whisper ses transkripsiyon hattı
  ├── sidecar/                         # Opsiyonel VoIP yardımcı modülü (node_modules hariç)
  ├── install/                         # [YENİ] Modüler Kurulum Orkestratörü (make install -> ./install/install.sh)
  │   ├── install.sh & lib.sh          # Orkestratör + ortak yardımcılar (sudo komut toplayıcı, güvenli JSON/YAML yama)
  │   ├── 00-system-deps.sh            # Sistem/medya araçları (go, gcc, uv, ffmpeg, tesseract, pdftotext, spectacle, wl-clipboard)
  │   ├── 10-core-go.sh & 20-python-env.sh # Go ikilileri (atomik ~/bin kurulumu) + ortak .venv kurulumu
  │   ├── 30-channels-mcp.sh           # WhatsApp (+ Faz 3 Telegram/Discord/Instagram) kanal ve MCP kurulumu
  │   ├── 40-voice-modules.sh          # STT (faster-whisper CUDA 12 + model) ve TTS (Piper TTS + Türkçe ses) kurulumu
  │   ├── 50-vision-modules.sh         # Ollama servis/model (qwen2.5vl:3b, qwen3:4b) ve OCR (tur+eng) kurulumu
  │   └── 60-agents-config.sh & 70-configs-services.sh # Ajan (agy, claude, opencode) MCP/Skill + systemd kurulumu
  ├── upgrade/                         # [YENİ] Modüler Yükseltme Orkestratörü (make upgrade -> ./upgrade/upgrade.sh)
  │   ├── upgrade.sh                   # Orkestratör (--all, --only <modül>, --rollback, --with-models)
  │   ├── 00-backup-snapshot.sh        # Yükseltme öncesi tam yedek (.tetikleyici/backups/: config, session.db, ajanlar, ikililer)
  │   ├── 10-upgrade-core.sh           # Go/Python güncellemesi ve çalışan süreci bozmayan atomik ikili değişimi
  │   ├── 20-upgrade-channels-mcp.sh   # Whatsmeow güncellemesi (405 önleyici), kanal SDK'ları ve yeni MCP araç senkronizasyonu
  │   ├── 30-upgrade-media-models.sh   # Ses (Whisper/Piper) ve görüntü (Ollama) modellerinin güncellenmesi
  │   └── 40-migrate-configs.sh & 50-verify-and-restart.sh # Kayıpsız YAML şema göçü + 627+ test kapısı (auto-rollback) + restart
  ├── docs/                            # PII temizlenmiş tasarım belgeleri ve Nerve Dispatch Master Blueprint
  ├── tests/                           # Birleşik Python test süiti (test_mcp_handlers.py + 10 tetikleyici testi)
  ├── go.work                          # Go Workspace (./bridge ve ./tetikleyici)
  ├── requirements.txt                 # Tekil birleşik Python bağımlılıkları (.venv)
  ├── Makefile                         # İnce görev yönlendiricisi (make install -> install/, make upgrade -> upgrade/, verify)
  └── README.md & TODO.md              # Tekil dokümantasyon ve yalnızca açık işleri tutan sprint listesi
  ```
- [ ] **1.0 (2026-11-09 eklendi) - Kaynak repolardaki commitlenmemiş işin kaydı:** `whatsapp-bridge-mcp` 28, `whatsapp-tetikleyici` 91 değişiklik granüler commitlenir (push yok). Taşıma dosya kopyası değil, geçmiş korunarak yol yeniden yazımı. Ayrıntı: plan Adım 1.0, `NOTES.md` §6, `ASK.md` 2026-11-09 04:00.
- [ ] **2. FAZ 1 - Doğrudan `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` Reposuna Geçiş, `install/` & `upgrade/` İnşası ve Ekosistem Onarımı (Adım 1.1 – 1.7):**
  - `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/` içinde Git başlatılması, kodların ve testlerin birleştirilmesi.
  - Canlı `tetikleyici.yaml`, `pasif-kurallar.yaml` ve `.tetikleyici/` durum dizininin (iç içe `.tetikleyici/.tetikleyici` temizlenerek) aktarılması; `~/.local/share/whatsapp-bridge/session.db` ve `messages.jsonl` oturumunun QR bozulmadan korunması.
  - `docs/` PII taraması, `TODO.md` içindeki bitmiş maddelerin `TODO.DONE.md`'ye ve kalıcı kuralların `NOTES.md`'ye taşınması.
  - İnce `Makefile` ile `install/` ve `upgrade/` modüler betiklerinin oluşturulması; `make install` ile `bridge/internal/appcfg/appcfg.go`, `systemd` servisleri, `~/bin/tetikleyici` (`~/bin/nerve`), Antigravity (`mcp_config.json`), OpenCode (`opencode.json` + 5 yeni MCP aracı), Claude (`~/.claude.json`), `SKILL.md` ve `INDEX.md` içindeki eski yolların `/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch` dizinine bağlanması.
  - `make verify` ile 627+ Go ve Python testinin `NerveDispatch` üzerinde %100 yeşil geçtiğinin kanıtlanması.
- [ ] **3. FAZ 2 - Önce WhatsApp Hattını ve Çekirdek Yetenekleri Tamamlama (Geliştirme 2.1 – 2.11):**
  - `2.1` Sabit Session ID & Ajan Adaptörleri (#7) → `2.2` Birleşik `nerve` CLI (#13A) → `2.3` Gömülü Web Dashboard (#13B) → `2.4` 8 Haneli Eşleştirme Kodu (#21) → `2.5` Tepki Emojisi (#12) → `2.6` `/ekran` Komutu (#10) → `2.7` İnteraktif Anket Onayı (#9) → `2.8` Dosya Transferi `/get` (#14) → `2.9` `/root`, `/clip`, `/apps`, `/kill` (#15-#17) → `2.10` `/compact`, Akıllı `/zamanla`, Sistem Bekçisi (#11, #19, #20) → `2.11` Piper TTS Sesli Yanıt (#18) & 300 Satır / 15 Dosya Refactoring.
- [ ] **4. FAZ 3 - Çok Kanallı Genişleme: Telegram, Discord, Instagram (Geliştirme 3.1 – 3.4):**
  - `3.1` Ortak `ChannelAdapter` arayüzü ve `ChimeraEvent` (`channel: "whatsapp" | "telegram" | "discord" | "instagram"`) zarf standardı.
  - `3.2` **Telegram Kanalı (`channels/telegram`):** Mefamex Telegram ID (`975678139`) entegrasyonu, tüm `/cmd`, `/camera`, `/ekran`, `/agent` komutlarının ve sabit AI oturumunun Telegram'dan da çalışması + WhatsApp düştüğünde Telegram üzerinden otomatik yedek bildirim (**Cross-Channel Failover**).
  - `3.3` **Discord Kanalı (`channels/discord`):** WebSocket Gateway, DM/Sunucu komutları ve zengin Embed yanıtları.
  - `3.4` **Instagram Kanalı (`channels/instagram`):** DM metin/görsel/ses otomasyonu ve tüm kanalların Web Dashboard üzerinde tek ekrandan izlenmesi.

---

## Kritik / Bloke

**Yok.** Köprü ve tetikleyici ayakta, yanıt döngüsü ölçülmüş çalışıyor.
Medya hatları (görsel / ses / `/camera`) 23:21-23:22'de **uçtan uca CANLI
doğrulandı** (`EŞLEŞTİ: … mefamex` → ilgili hat → cevap sohbete düştü).
**ping → restart → catch-up da TAMAM** (23:28 `ping` → sahip; 23:32 restart →
23:32 sessiz temiz catch-up; 23:11'de uyarılı hali görülmüştü - iki yön de
canlı kanıtlı). Canlı sticker hatası (23:28) → düzeltildi + pipeline'la
yeniden doğrulandı (kayıt: `TODO.DONE.md` → "00:05 canlı sticker" girişi).

---

## BAKIM

> bakım başlığı bu TODO.md dosyasından silinmeyecek ara sıra bu başlık altındakiler kontrol edilecek. 
> Son bakım zamanı: 2026-10-06:01.00

- Proje için oluşturulan çöp dosyaların kontrolü ve arka planda çalışan eski çöp işlemler temizlenmeli
- Projedeki TODO.md dosyasında yapılanlar  TODO.DONE.md'ye taşınmış mı?
- README.md projeye uygun mu?
- ask.md cevaplanmış mı? cevaplandıysa sık olmamak şartı ile admine whatsapp'tan sorulmalı.
- Bu projenin ana dizinindeki tüm notlar (md veya txt dosyaları) doğru yerde ve içerikleri doğru mu?


## ✅ aynı kişi için çoklu config - TEK KİŞİ TEK KURAL (bitti)

> **Bitti (2026-10-06, commit `9ab63bc`).** 6 kural → 1 kural
> (`mefamex-oturum`, 8 yetenek bayrağı `allow_*` + `*bool` global ezer).
> Tam tasarım + kayıt: `TODO.DONE.md` → "Tek kişi tek kural" +
> `config.readme.md` §2.1 (tek kural örneği, hat haritası, `--format json`).
> Sahibin ek notu (YAPILMADI, bilinçli): oturum çıktısının **anında canlı
> yayını** - **PARK** (canlı yayın için önce kuyruk şart; kayıt DONE'da).


## 🔴 KAYIP MESAJ - kalan: WhatsApp history-sync LID içeriği boş (kök neden whatsmeow'da)

> 03:10 vakası + uygulanan düzeltmeler (`probeWidth` 10→50, ⚠️ uyarısı):
> `TODO.DONE.md` → "2026-10-06 03:33 +03 - catchUp probeWidth artırıldı".

- **Açık kalan:** WhatsApp `RequestHistory`'ye LID sohbetleri için **0 mesaj**
  dönüyor (`HISTORY_SYNC_NOTIFICATION` geliyor, içerik yok). Çözüm köprü/
  whatsmeow seviyesinde; tetikleyici uyarı loglayıp devam ediyor.
- **Kullanıcı testi:** `ping` at → 3 sn bekle → tetikleyici restart → 10 sn
  bekle → `grep -a "YAKALAMA\|⚠️" .tetikleyici/tetikleyici.log | tail` ve
  `grep -c <mesaj-id> .tetikleyici/seen.txt`.

---

## Öncelik: Kritik ve acil değil

- proje "/mnt/DiskD/SOFTWARE/Github-repos/" dizininde taşınıp çalıştırılacak.
- NOTES.md proje notu içindi. çok uzatılmasın. 
- sohbet notlarını bulgularını ayrı bir md'de tutulsun.
- tüm md dosyaları veya periyodik güncellenmesi gerekenler uygun bir yere not edilsin. 

> Bitenler taşındı: `servis-yonet.sh` ✅, `/session list|choose` ✅, `debug_cikti` ✅
> → kayıtlar `TODO.DONE.md`'de.

## ✅ ÖNCELİK 1 - Komut yetkisi: `/cmd` (bitti)

> Tam kayıt: `TODO.DONE.md` → "2026-10-06 02:30 - `/cmd` komut yürütücü eklendi".
> `tetikleyici/cmd/` (7 modül) + `mefamex-cmd` kuralı + `prefix`/`not_prefix` +
> `debug_cikti` + `/camera` + `/durum` - hepsi DONE'da. `ping` kuralı da
> **SİLİNDİ** (kayıt DONE'da; sahip: *"pinge gerek yokkkkkkk zaten /help var ya"*).
> **Son eklenen (2026-10-07):** `/session list` + `/session choose <n>` →
> `cmd.py` + `yardim.py` (opencode API port 4096). Canlı test: **Connection
> refused** - opencode API henüz başlatılmadı; API ayağa kalkınca çalışacak.
  
## ⛔ KALICI KURAL - kapalı hat sessizdir

Burak / Yusuf / Masha'nın tetikleyici kuralları **geçici olarak kapalı**
(2026-10-06). Onlara **hiçbir otomatik cevap verilmez**:

- ✗ "şu an kapalıyım" · ✗ "mevcut değil" · ✗ "mefamex'e yaz" · ✗ ping karşılığı
- ✗ hiçbir sabit metin - kural yoksa sonuç da yoktur

> "kesinlikle böyle bişe yok. kapalıysa kapalı bunu kritik olarak not et."

Bunu **önerdim ve reddedildi** - bir daha önerme. Kayıt:
`NOTES.md` madde 31 · `tetikleyici.yaml` ⛔ bloğu.

---

## ✅ Canlı testler (tamamlandı - kayıt DONE'da)

> Fotoğraf ✅ 23:22 · ses ✅ 23:22 · çift gönderim sayımı ✅ · `/camera` ✅ 23:21 ·
> `/zamanla kapat` ✅ 23:22 · **ping → restart → catch-up ✅ 23:28-23:32** (sessiz
> temiz catch-up; engine.go:378 yalnız `0 < gelen < istenen/2` iken ⚠️ basar).
> Tam metinler: `TODO.DONE.md` → "23:21-23:22 CANLI doğrulama" + "00:05 canlı
> sticker". Restart kuralı: `./tetikleyici.sh restart` (elle `>> tetikleyici.log`
> yönlendirme fd1/fd2+fd3 çakışmasıyla log'u ÇİFT satır yapar).

## ⏳ Kullanıcıda (sudo gerekiyor)

**Bekleyen sudo işi YOK** (anket deploy'u yapıldı ✅ - aşağıdaki şablon
sonraki köprü değişiklikleri için duruyor).

~~Bekleyen 1 sudo işi (anket deploy)~~ - yapıldı: sahip restart attı, `/send-poll`
400 ✅ + `Pıttık mı?` anketi mefamex'e gönderildi (`ok:true`, store `out|poll`).
Geçici betik `~/Desktop/whatsapp-bridge-restart.sh` duruyor, silinebilir.

Genel şablon (sonraki köprü değişiklikleri için):

```bash
cd ~/Desktop/whatsapp-bridge-project
cp /tmp/whatsapp-bridge.new ./whatsapp-bridge.new
mv -f ./whatsapp-bridge.new ./whatsapp-bridge
sudo systemctl restart whatsapp-bridge
```

> `cp` doğrudan çalışan binary'nin üstüne yazarsa **her zaman** `Text file busy`
> verir. `cp` → `mv -f` sırası şart (`mv` atomik, açık dosyaya yazmaz).

✅ 2026-10-04 19:17'de deploy edildi: SSE `WriteTimeout: 0` (30 sn'lik kopma
gitti), heartbeat 15 sn, protokol gürültüsü depoya yazılmıyor. Çalışan
binary'de doğrulandı - `strings whatsapp-bridge` → `15s` var, `30s` **yok**;
son 10 dakikada 0 `unsupported`/`HISTORY_SYNC` satırı. Tetikleme gecikmesi
0-33 sn aralığından **0 sn**'ye indi.

---

## Sesli not / görsel hatları

📄 Ölçüm koşularının tamamı: `docs/GORSEL-HATTI-2026-10-06.md` · bitenler `TODO.DONE.md`'de.

Açık kalanlar:

> ✅ Video (`33da7bc`, ölçülü) + sticker (canlı düzeltme dahil, 00:05) +
> document (`mefamex-belge` kuralı) bitti - kayıt: `TODO.DONE.md`.

- [ ] **Bilinen model sınırı**: `b-gorme`'daki 4 yeşil daireyi tüm varyantlarda "dört yeşil nokta" diyor. Prompt ile düzeltilemiyor - bu modelin kalıcı sınırı, kabul edildi
- [ ] **Sesli arama bildirimi**: WhatsApp'tan kullanıcıya arama (sadece çaldırma, sesli konuşma yok) ile ulaşma/uyarma - sadece çaldırma, belirli süre veya o kapatana/açana kadar
      🧩 ALTYAPI HAZIR (2026-10-06, `111c43b`): gelen call event'leri köprüde
      SSE `call` event'ine yayınlanıyor (offer/accept/reject/terminate…).
      → UÇTAN UCA KURULDU (2026-10-07, detay: Köprü bölümü); gerçek çaldırma testi sahip uyanıkken.
- [ ] **En son geliştirilecek özellik**: WhatsApp'tan sesli olarak kullanıcın asistanla konuşabilmesi ...

> **Tüm medya hatları aynı tuzağa açık:** `code:` içinde `{mesajlar}` yazma -
> o metindir. Mesaj kimliği `$TRIGGER_ID`. Yanlışsa `404` alır, tetikleyici
> "başarısız" diye loglar ve **sessizce** cevap göndermez.

## Tetikleyici - mevcut hali ve kalan işler

> 📄 **Son çalışma (2026-10-05):** `docs/TESLIM-SORUNU-2026-10-05.md` -
> "iş bitiyor ama WhatsApp'a gitmiyor" sorunu **çözüldü ve canlıda**.

Mevcut yapı: `tetikleyici/` (bağımsız Go modülü) + `tetikleyici.yaml` +
`tetikleyici.sh`. **Canlı kurallar (1):** `mefamex-oturum` (8 yetenek
bayrağı `allow_*`; eski `mefamex-ping`/`-sesli-not`/`-gorsel`/`-cmd`/
`-belge` tek kurala birleşti). Misafir kuralları kapalı, `tetikleyici.yaml`
altında yorumda hazır. Biten işler `TODO.DONE.md`'de.

Açık kalanlar:

- [x] **Komut yetki modeli yeniden kurulsun** (ask.md cevabı işlendi, 2026-10-07) -
      **CANLI (commit `296a289`):** 3 aşama: 1) Config users listesi (mefamex:
      numara+LID, allow_slash/cmd/camera/durum/zamanla/reply) 2) yetki.py 3
      aşama kontrol 3) /help güncellendi. Test: 144 passed. `/cmd echo` ✅
      Kayıt: `TODO.DONE.md` → "2026-10-07 01:25 - 3 aşamalı yetki modeli CANLI"
- [x] **Mesajlar paralel çalışsın** (ask.md cevabı işlendi, 2026-10-07) -
      **CANLI (commit `d957f6b`):** `engine.dispatch` actions'lar `goroutine` +
      `sync.WaitGroup` ile paralel. Session kilidi (`oturumKilidiAl`) `runAction`
      içinde → aynı session sıralı, farklı session'lar paralel. Farklı sohbetler
      zaten per-chat worker ile paraleldi. Test: 144 passed. Sahip fikri
      ("/cmd -n sezon devamı") → kapsam netleştirilecek, ayrı karar bekliyor.
- [ ] **Sesli arama bildirimi**: WhatsApp'tan kullanıcıya arama (sadece çaldırma) - belirli süre veya o kapatana/açana kadar
      ✅ UÇTAN UCA KURULDU (2026-10-07): köprü `POST /call` + tetikleyici
      `/call <jid> [süre]`; canlı `/call` → `ok:true`. Gerçek çaldırma testi
      sahip uyanıkken (detay: Köprü bölümü).

### ✅ Canlı bildirim - Mefamex'e "çalışıyorum" mesajı (bitti; sözleşme kaldı)

> ✅ **TAMAMLANDI + CANLIYA KURULDU 2026-10-06** (`cc7fecc`). Dört olay canlı
> ölçüldü (on_start 09:17 · on_config_change · on_stop · restart-throttle -
> "restart'ta da hatırlar" varsayılmadı, ölçüldü: bildirim.txt damgası).
> **Neden:** tetikleyici sessizce ölürse / köprü düşerse bunu hiçbir kanaldan
> öğrenemiyoruz. Tam kayıt: `TODO.DONE.md` → 09:20.

**Ne zaman gönderir (sözleşme):**

| Olay | Gönderir mi | Neden |
|---|---|---|
| Tetikleyici başladı **ve** köprü bağlı | ✅ "çalışıyorum" | her açılışta 1 kez |
| Tetikleyici başladı, köprü `connected:false` | ⚠️ "başladım ama köprü bağlı değil" | teslimat hattının nerede koptuğunu ayırır; 30 sn'de bir dener |
| Köprü bağlandı (SSE `connected` olayı) | ✅ "köprü geri geldi" | tetikleyici izlemeye başladı |
| Köprü düştü (`disconnected` olayı) | ✅ "köprü düştü, 30 sn içinde yeniden bağlanılıyor" | kullanıcı boşuna mesaj atmasın |
| Config dosyası değişti | ✅ "config değişti: N kural" | yanlış kural yazdığında fark edilir |
| Temiz durdurma (SIGTERM/SIGINT) | ✅ "durdum" | neden cevap gelmediği belli olur |
| Çökme (panic / `kill -9`) | ❌ | süreç ölmüş; **sonraki açılışta** "geçen sefer beklenmedik kapandı" |
| Uzun süre tetik yok | ❌ | saatlik "hala yaşıyorum" gürültü; sorarsa `/durum` yazsın |

> `on_bridge_down` yalnız /health GEÇİŞİNDE tetiklenir (SSE'nin 30 sn'lik rutin
> kopması köprü düşüşü değildir). Mesaj formatı: WhatsApp'a uygun düz metin;
> `chatText` Markdown→WhatsApp çevirdiği için `*kalın*` kullanılabilir.

**Sınırlamalar (spam olmasın):** aynı tür bildirim **5 dk içinde bir daha
gitmez** (`.tetikleyici/bildirim.txt` - restart'a dayanır) · köprü bağlı
değilse **en fazla 10 kez** 30 sn arayla dener, sonra günlüğe yazıp vazgeçer ·
yalnız `notify.to`'daki kimliğe gönderir - **kendi hesabımıza asla** ·
gönderim başarısızsa hata mesajı **gönderilmez**, sadece loglanır.

**Yerleşim:** config'e ayrı `notify:` bloğu (kural değil - tetikleyicinin
kendi işi). Yeni dosya: `tetikleyici/notify.go` + `notify_test.go`:

```yaml
notify:
  to: ["905550000001", "12345678901234@lid"]
  on_start: true
  on_stop: true
  on_bridge_up: true
  on_bridge_down: true
  on_config_change: true
  throttle: 5m
  retry: 30s
  max_retry: 10
```

**Yan yana - `/durum` komutu:** ✅ YAPILDI (kayıt: TODO.DONE.md → 04:20).
Ayrı kural YOK; `mefamex-cmd`'nin `prefix: /` koşusu `/durum`'u yakalar,
`cmd.py` yönlendirir (düz `durum` yazılsa `mefamex-oturum` da tutar, ÇİFT
ateşlenir; `/durum` `not_prefix` sayesinde yalnız cmd'ye düşer). Çıktı:
tetikleyici PID/süre (`/proc` exe eşleşmesi - `pgrep -f` kendini vurur) ·
köprü `/health` (bağlı/uptime/tampon) · kural adları · son 3 `EŞLEŞTİ`.

**Yan yana - `/session list|choose` komutu:** ✅ YAPILDI (2026-10-07, commit `296a289`).
`cmd.py` + `yardim.py` + `sozdizimi.py` (`session_list`, `session_choose`).
opencode API port 4096 kullanır. **Canlı test: Connection refused** - opencode
API (port 4096) henüz başlatılmamış; `opencode serve` veya TUI ile API
aktif edilince çalışacak. Kayıt: `TODO.DONE.md` → session komutları.

**Zamanlayıcı (opencode-zamanlayici) - auto-restart loop:** ✅ AKTİF (2026-10-07).
`run_loop.sh` (while-true wrapper) → PID dosyası `/tmp/zamanlayici_loop.pid`
→ `kill $(cat /tmp/zamanlayici_loop.pid)` ile durdur. Watchdog 60sn,
brifing 3600sn. Log: `/tmp/zamanlayici_loop.log`. Config:
`/home/mfmx/Desktop/opencode-zamanlayici/config.yaml`. Kayıt:
`NOTES.md` + `AUTO_RESTART.md`.

## Köprü - yapılabilecek işler

> Anket endpoint'i `TODO.DONE.md` 04:55. Bitenler DONE'da: etkinlik 🚫
> (API yetersiz, tulir/whatsmeow#1214 açık · 09:20) · restart persistence
> (`1594d8f` CANLI 22:01 · 18:25) · MCP `whatsapp_ask_confirm`/`quick`
> (`0033579` · 11:35) · tehlikeli mesaj uyarısı (22:40) · zamanlanmış mesaj
> (11:21) · `/mute`+`/star` (08:20).

Açık kalanlar:

- [ ] **Tema/duvar kağıdı endpoint'i** (`POST /chat-theme`): Jon Snow/kış teması. whatsmeow'da `ChatThemeSettingView` **sadece okunuyor**, yazan metot yok - protobuf'u elle kurmak gerekir. API yetersiz
- [ ] LID `[unsupported]` sınırı: LID sohbetlerinde telefon numarası hiçbir payload alanında yok (`ContactInfo.RedactedPhone` maskeli). **Gönderimde** çözüyoruz (`resolveTarget`), **gelen** tarafta çözecek alan yok → `--isim` (push_name) veya LID'in kendisi gerekiyor.
      ✅ GELEN YÖNDE İYİLEŞTİRME (2026-10-07): `handleMessage` artık LID gönderen/sohbet
      için store'daki bilinen telefon karşılığını `TakmaAd`'a ekliyor (`lidTelefonu`,
      2sn sorgu sınırı; `From` DEĞİŞMEZ - store/dedup/eşleşme etkilenmez; eşleşme
      yoksa sessizce atlanır). Dış istemci telefonu ek kimlik olarak görür.
- [ ] **WhatsApp'tan `/` komutlarını genişletme** - kısmen yapıldı: `/camera`, `/durum`, `/zamanla`, `/ping` eklendi (`prefix: /` + `mefamex-cmd`). Komut yetki modeliyle birlikte ele alınacak (yukarıdaki açık iş)
- [ ] **Sesli arama bildirimi**: WhatsApp'tan kullanıcıya arama (sadece çaldırma, sesli konuşma yok) ile ulaşma/uyarma - sadece çaldırma, belirli süre veya kapatana kadar
      🧩 ALTYAPI HAZIR (2026-10-06, `111c43b`): gelen call event'leri köprüde
      SSE `call` event'ine yayınlanıyor (offer/accept/reject/terminate…).
      ✅ UÇTAN UCA KURULDU (2026-10-07): köprü `POST /call` + tetikleyici
      `/call <jid> [süre]`; canlı `/call` → `ok:true`.
      ⚠️ GERÇEK ÇALDIRMA DOĞRULANMADI: stanza'da SDP placeholder. Canlı test
      SAHİP UYANIKKEN - uyurken ek test YOK.
- [ ] **Sesli sohbet**: WhatsApp'tan sesli olarak kullanıcı ile asistan konuşması (en son geliştirilecek özellik)

## Ölçüm betikleri (`tools/`) - üretim değil, kanıt üretir

Bu betikler üretim koduna **dokunmaz**; yalnız ölçer. Silinmez çünkü bir
sonraki ölçüm aynı hatanın tekrarlanıp tekrarlanmadığını ancak onlar gösterir.

> Envanter (`karsilastir/prompt_varyant/sizinti/ocr_yol/gorsel_test_uret`): `TODO.DONE.md` 04:30 temizlik kaydında.
> `prompt_tune.py` başlığına GEÇERSİZ uyarısı eklendi ✅ (04:35, derleniyor).

## Ölçüm metodolojisi uyarısı

- **`qwen2.5vl:3b` koşular arası ±20-40 puan salınıyor.** Tek koşuluk prompt
  sıralaması güvenilmez. İstem seçimi ancak **yayılım ölçülerek** yapılır
- **Ayırt edici görsel şart.** İlk 3 görselin 2'si beş varyantta da birebir
  aynı puanı aldı - hiçbir şey ölçmüyordu. `d-kucuk` eklendi
- **Puanlama ağırlığını koşular arası değiştirme.** 2. koşuda köprü cezası
  10→0 puana indi; mutlak puanlar artık karşılaştırılamaz. Değiştirdiysen
  not düş
- **Mekanizma kanıtı puanı geçer.** "Model şunu yaptı" gözlemi, "varyant şu
  kadar iyi" tahmininden güçlüdür. İstem kararları mekanizmayla gerekçelendirildi

## ✅ Desteklenmeyenler - araştırıldı, bitti

> Global mesaj arama: WhatsApp'ın kendi arama mekanizması API'de YOK →
> `/search?q=`'yi biz yazdık (`e86fca8`, kayıt: DONE 18:25). Sohbeti sessize
> alma `/mute` ✅ + favori `/star` ✅ + kanıt ucu `/chat-state` ✅ (8h/7d/1y;
> kalıcı=1y - WhatsApp'ta kalıcı susturma yok; `ok:true` sessiz başarısızlık
> kanıtı dahil, tam kayıt: DONE 08:20).

## Skill / ajan tarafı

> `~/.codex/skills/whatsapp-bridge/SKILL.md` yazıldı ✅ (kayıt: DONE 04:45).

- [ ] `~/.config/opencode/config.toml`'a `agents["whatsapp-bridge"]` tanımı -
  **PARK**: dosya yok, opencode JSON okuyor (`opencode.json`); ölü config yazılmaz
- [ ] `wa-voip-sidecar` WASM binary'si - **ATLANDI**: boyut önden doğrulanamıyor,
  `willInstall.md`'de (`npm run fetch-wasm` kullanıcıda)

## Burak / Yusuf - sohbet tarafı bekleyenler

- [ ] 2. resim (denizli yer) netleştirme: Side / Phaselis / Olympos adayları verildi, tekrar atması istendi. Gelince `/media/{id}` → `identify` + `exiftool` + Tesseract → kaynaklı cevap
- [ ] KPSS 2026 takibi: 25 Eki 2026 sınav, sonuç 19 Kas 2026. Giriş belgesi Ekim 2. haftası bekleniyordu - sorarsa güncel kontrol et
- [ ] Opencode kurulum rehberi: işletim sistemini sor → adım adım kurulum mesajı
- [ ] Yusuf otomasyonu **kapalı** - kullanıcı rızası olmadan otomatik cevap yazılmıyor. Dosyalar yerinde, yeniden açmak istenirse tek satır config yeterli

## ⏱ Heartbeat / durağanlık kontrolü (Mefamex) - sahip onayı GÜNCELLENDİ

### Kalp atışı - CANLI (2026-10-07)

**Durum:** `heartbeat:` bloğu yaml'da CANLI (idle 15m / every 1m, mesaj birebir).
**BUG BULUNDU + DÜZELTİLDİ (2026-10-07 05:00):** nudge 4 tur üst üste
`config'de "opencode" yolu tanımlı değil` diye düştü - `main.run`
`opts.opencode` doldurmuyor, heartbeat startup kopyasını kullanıyordu.
Düzeltme: `tazeleOpts()` her turda yolu taze config'den alır + `main.go`
ilk doldurma + 2 test. Binary derlendi, tetikleyici restart (PID 307670).
Canlı doğrulama TAMAM (05:11:21): sessizlik saati doldu → nudge çalıştı →
`kalp atışı cevabı gitti (497 karakter)` → ajan cevabı (`Heartbeat çalışıyor ✅…`)
sohbete düştü (store'da görüldü). Döngü kırıcı da doğrulandı: cevap trafik
sayıldı, sonraki saat sıfırlandı.

**Sahip (03:33):** eşik **15 dk**, kontrol **1 dk**; mesaj (kelimesi
kelimesine): *"admin ile konuşman bir süredir durağan. ona mesaj atman
gerekiyorsa at"*.
**Sahip (2026-10-07):** *"orada bir karışıklık oldu… bu gereksiz karmaşa
olayını arşivden de silip gerekiyorsa heartbeat'ı uygula"* → **uygulanabilir.**
ask.md "İki hesap" cevabı da işlendi (sahip): 4 → **yalnız hatırlatma**
(kontrol mesajıdır, otomatik soru yok); 5 → **ayrı ajan YOK** -
whatsapp-bridge her sohbetten gelen mesajı aynı sabit session'a iletir;
1-3 (şahsi hesap cevabı, eşik, sıklık) bana bırakıldı →

**KARARLAR:** 1) Şahsi hesap (905550000005): otomatik cevap **YOK** (orada
kural yok, sessiz kalır; isterse tek satır config) · 2) Eşik **15 dk**
(sahibin verdiği değer) · 3) Sıklık **1 dk** kontrol (sahibin verdiği).

**Çelişki çözümü:** "sürekli cevap" ≠ "durağan kontrol" - mefamex-oturum her
GELEN mesaja cevap verir; heartbeat yalnız 15 dk HİÇ trafik yoksa tetiklenir
(aktifken zaten tetiklenmez, döngü kendiliğinden kırılır).
Karar tablosu + yaml taslağı: `docs/archive/2026-10-06-ask-heartbeat.md`
(geri-alma karışıklığı kaydı sahip talimatıyla silindi).

### 🔴 YENİ - mesajlar okunmuyor (düz metin, tek satır)

**Sahibi:** *"bana mesaj atarken niye dümdüz mesaj atıyon. hiç satır atlama veya
numaralandirma yapmıyon. okunması zor oluyo."*

**Sorun:** uzun cevapları tek paragraf halinde yazıyorum. 5 satırlık cevap tek
ekranda okunmuyor.

| Yapılacak | Nerede |
|---|---|
| Cevaplara **satır araları** koy | ajan davranışı (kurallar md'de) |
| Listelenen şeylerde **numaralandırma** kullan | ajan davranışı |
| Uzun cevapta **ilk satır özet** olsun | ajan davranışı |

> ⚠️ Bu ajanın **kendi** yazım kuralı; kod değişikliği gerektirmiyor. Dikkat:
> `tetikleyici/cmd/chattext.go` Markdown'ı WhatsApp biçimine çeviriyor -
> `**kalın**` → `*kalın*`, satır başı `-` → `•`. Yani numaralı liste gönderilebilir.

### ✅ Yapılanlar (04:05 oturumu) - taşındı

> chmod +x 126 · medya yetkisi kanıtı · kayıp mesaj bulgusu · log/grep dersi ·
> yanlış bulgu düzeltmeleri → `TODO.DONE.md` 00:10 temizlik girişinde (hepsi
> DONE'daki 04:30-04:45 kayıtlarında da var).

## 📐 2026-10-06 05:10 - Kod kuralı: dosya 300 satırı geçmez

Sahibi README'de böyle bir kural olup olmadığını sordu → **yoktu**, eklendi
(`README.md` → *"Kod Kuralı: Modüler ve Profesyonel"*). Ölçümle desteklendi:

| Ölçüm | Değer |
|---|---|
| Üretim kodu (test hariç) | 10.115 satır |
| 300 satırı aşan üretim dosyası | **5** |

**İhlal edenler (bölünmeyi bekliyor, sahibinin emri gerekli):**
`tetikleyici/action.go` 518 · `main.go` 431 · `config.go` 370 ·
`engine.go` 346 · `cmd/calistir.py` 334

- [ ] `action.go` → koordinasyon + `command.go` + `session.go` + `gonderim.go` + `prompt.go`
- [ ] `calistir.py` → süreç başlatma + tekrar döngüsü ayrı
- [ ] Her bölünmede testler yeşil kalmalı (`go test ./...`)

> Test dosyaları kural **dışı**: onlar üretim karmaşıklığı değildir.

## 📁 2026-10-06 05:15 - Klasör kuralı: ≤ 15 dosya, ≤ 3 seviye

Sahibi istedi, README'ye eklendi (*"Klasör Kuralı: sığ değil, derin de değil"*).
**Yalnız md değişti, kod yapısına dokunulmadı** (sahibi: *"sen şimdilik yapma,
sadece md dosyalarına bakıyoruz"*).

| Sınır | Değer |
|---|---|
| Klasörde dosya sayısı | **≤ 15** |
| Derinlik | **≤ 3 seviye** |
| Yapay bölme | yasak (3 dosyalık klasör açma) |

**Ölçülen ihlaller (yapılmadı, emir bekliyor):**
- `tetikleyici/` = **31 dosya** → test dosyaları ayrılırsa ~20
- `bridge-go/internal/wabridge/` = **17 dosya** → send + medya alt paket olabilir

**Uyumlu:** `httpapi/` 11 · `cmd/` 10 · `gorsel/` 10 · `tools/` 9

- [ ] `tetikleyici/` test dosyalarını ayır (31 → ~20)
- [ ] `wabridge/` send ve medya alt paketlerine böl (17 → ~10)

## ✅ 🔐 MEVZU - "config'e yazılan şey atılacak mı?" (kademeli yetkilendirme, bitti)

**Karar: kademeli yetkilendirme.**

| Aşama | Nerede | Değer | Anlamı |
|---|---|---|---|
| **1** | **Kod içi varsayılan** | `false` | Hiçbir yerde yazılmadıysa özellik **kapalıdır** - açık kapı yok |
| **2** | **Config başı** (yaml üst yorum bloğu) | yazılır/yazılmaz | Tek bakışta neyin açık olduğu; yazılmadıysa aşama 1 |
| **3** | **Kullanıcı ayarı** (çalışma anı) | yazar | **Config başını ezer** - en yüksek öncelik |

> Öncelik: kullanıcı ayarı → config başı → kod (false). Kod `true` varsayılanla
> yazılmak **YASAK**. *"Kullanıcı istemezse hiç olmaz"* felsefesinin kardeşi:
> hat kapalıysa **sessiz**, özellik kapalıysa **çalışmaz** (ikisi tabloyla
> karşılaştırılıyor, `config.readme.md`).

**Uygulandı (tamam):** aşama 1 → boolean izinler `false`/nil default (`config.go`,
`acik: nil → kapalı`; inceleme kuralı `config.readme` §2.1) · aşama 2 →
`tetikleyici.yaml` başına **"AÇIK OLAN ÖZELLİKLER"** yorum bloğu (strict YAML
etkilenmez, restart gerekmedi) · aşama 3 → `/zamanla aç|kapat` + engine her
okumada tazeler → **restart'sız** yansır (`b7b465f`+`5fbd7f5`+`73c0247`; tarayıcı
kapalıyken tur atmaz) · testler: yalnız-1 / 1+2 / 1+2+3(ezen) ✅ (`config_test` +
`tests/test_zamanla.py` 7 test) · tam metin `config.readme.md` §2.1 (`2bc70f0`);
"mavi tik" iddiası iki yeni testle + mutasyonla doğrulandı.

Kayıt: TODO.DONE.md → MEVZU (18:25/21:40/21:46), /reply (18:25/21:40/23:05).
## ✅ 📨 `/answer` - oturumlar arası yanıt yönlendirme (bitti)

**Nerede:** **tetikleyicide** (MCP'de değil - yönlendirme zaten onun elinde:
`sendText(..., Trigger.From)`; MCP oturum bilmez). Sahibin tasarımı kabul:
`/answer1` yerine **`/reply` + otomatik numara**. Sözdizimi: `/reply list` ·
`/reply <n> <metin>` · `/reply <metin>` (yalnız 0 bekliyorsa) · numara **sohbete
değil soruya** bağlı, oturum kapanınca serbest.

**Uygulandı (tamam):** dosya mekanizması `.tetikleyici/bekleyen_cevaplar.txt`
(satır başına sayı; en küçük eksik pozitif atanır 1,2,4→3, boş→0) + zaman damgası
`.zaman` (10 dk TTL; Go `soruNotu`/`reply.go` + Python `cmd/reply.py` ortak;
süresi dolan numara serbest, uyarı **yalnız logda** - kullanıcıya mesaj yok, spam
önlemi). Kuyruk FIFO, çoklu soru destekli, numarasız /reply yalnız 0'a gider,
atomik yazım + chmod 600. ADIM 1 (`e56bc2d`, 17 Python testi) + ADIM 2 ilk yarı
(`6855313`): ajan cevabında "?" → numara atanır + "> cevaplamak için /reply veya
/reply N" notu eklenir (hatalı pozitif zararsız - yalnız not gider). 6 Go + 7 Python test.

**⏳ PARK EDİLDİ:** cevabın açık opencode oturumuna iletilmesi - sahip kararı
+ kota + `--agent` çift gönderim tuzağı; açılmadan önce dikkatli kurulum +
canlı test gerekir.

Kayıt: TODO.DONE.md → /reply (18:25/21:40/23:05).


## ✅ 📝 `/reply` - otomatik numara atama (bitti; tasarım sahibin)

**Sahibin tasarımı (kelimesi kelimesine):** *"numarayı otomasyon otomatik atasın.
baksın bekleyen cevap var mı diye. mesela o dosyanın içi: 1, 2, 4 - orada olmayan
sayı (3) otomatik atar. dosya boşsa 0 koyar ve numarasız /reply alır. cevap gelince
0'ı siler. soru sorulacağı zaman da altına not düşer: > cevaplamak için mesaja şu
şekilde başlayın: /reply veya /reply 3. ama answer yerine reply kullanalım ya"*

**Karar:** kabul edildi; tek düzeltme dosya adı → `bekleyen_cevaplar.txt`.
Neden dosya: restart'a dayanır, `cat`/`sed`/`awk` ile yönetilir, `chmod 600` gizli.
Akış: soru → en küçük eksik pozitif (boşsa 0) → metne not ekle; cevap → satırı sil;
çakışma koruması atomik yazım. `/reply` **ayrı kural** (`sozdizimi.py` ayrıştırır,
`cmd.py` minimal dispatch dalı). Altyapı kararı: env override `BEKLEYEN_CEVAP_DOSYA`.

**Uygulama:** yukarıdaki `/answer` bölümünde - ADIM 1 + ADIM 2 ilk yarı canlıda;
ikinci yarı (cevabın ajana iletilmesi) **park edildi** (sahip kararı + kota +
`--agent` çift gönderim tuzağı). Kayıt: TODO.DONE.md → 18:25/21:40/23:05.

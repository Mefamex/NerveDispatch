# NOTES — NerveDispatch Dahili Teknik Bellek ve Mühendislik Notları

> **Temel Felsefe:**  
> Kod dosyaları bir işin **"nasıl"** yapıldığını anlatır; git geçmişi **"ne zaman"** yapıldığını gösterir; `NOTES.md` ise o kodun **"neden"** o şekilde yazıldığını ve **"hangi alternatiflerin neden elendiğini"** açıklar.  
> Kod içine paragraflarca yorum satırı yığmak yerine, mimarinin görünmeyen kısıtları ve tecrübeleri burada diske sabitlenir.  
> **Referans Standart:** [Mefamex Conventions - NOTES.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/NOTES.md)

---

## 1. Mimari Tercihler ve Gerekçeler (Architecture Decisions & Trade-offs)

- **Go Çekirdeği + Python Ajan/Multimodal Ayrımı:**  
  - *Kazanç:* Ağ protokolleri, WhatsApp soketi (Whatsmeow) ve SSE dağıtım motoru yüksek eşzamanlılık ve düşük bellek tüketimi için Go ile koşar. LLM boru hatları, Whisper ses işleme ve sistem otomasyonları Python ekosisteminin zengin kütüphaneleriyle yönetilir.
  - *Ödün:* İki ayrı çalışma ortamının (Go binary + Python venv) orkestrasyonu gerekir. Bu durum modüler `install/` ve `upgrade/` mimarisiyle çözülmüştür.
- **Monorepo Birleşimi (NerveDispatch):**  
  - Dağınık `whatsapp-bridge-mcp` ve `whatsapp-tetikleyici` repoları tek çatı altında toplanarak ağ geçidi çok kanallı (WhatsApp, Telegram, Discord, Instagram) tek bir omurgaya evrildi.

---

## 2. Denenip Vazgeçilen Yollar (Discarded Approaches)

- **Telefon JID Üzerinden Doğrudan Mesajlaşma (Whatsmeow LID Tuzağı):**  
  - *Denenen:* `@s.whatsapp.net` üzerinden doğrudan yanıt göndermek.
  - *Sonuç / Hata:* Modern WhatsApp çoklu cihaz mimarisinde hesap sahibine giden mesajlar sessizce düşmekteydi.
  - *Karar:* Mefamex hesabına gönderimlerde `12345678901234@lid` (gerçek LID `tetikleyici.yaml` içinde) formatı zorunlu kılındı.
- **Tüm Modülleri Tek Dev Makefile İçinde Tutmak:**  
  - *Sonuç:* Makefile 400+ satıra ulaştığında bakımı imkansız hale geldi.
  - *Karar:* Kök `Makefile` sadece ince bir yönlendirici yapıldı; kurulumlar `install/00-...` ve güncellemeler `upgrade/00-...` betiklerine bölündü.

---

## 3. Doğrulanmış Sistem Sınırları ve Eşikler (Verified Constraints)

- **Test Eşiği:** 627+ birim ve entegrasyon testi hatasız geçmeden servis yükseltmesi yapılamaz.
- **Yeniden Başlatma Kuralı:** Tetikleyici süreci rastgele `kill` edilmez; canlı oturumların bozulmaması için yalnızca `./tetikleyici.sh restart` veya sistem yöneticisiyle yönetilir.
- **Kişisel Veri ve Gizlilik (PII):** WhatsApp telefon numaraları ve kullanıcı kimlikleri depoda ASLA açık bırakılmaz; log ve commit diff'lerinde taranır.

---

## 4. Platform Gariplikleri ve Tuzaklar (Quirks & Gotchas)

- **Whatsmeow SQLite Kilitlenmesi:** Whatsmeow oturum veritabanı (`session.db`) WAL modunda çalışmalıdır. Eşzamanlı okuma/yazma çakışmalarında busy timeout tanımlanmalıdır.
- **Wayland / Plasma Ortamı:** Linux KDE Plasma Wayland üzerinde çalışan grafik arayüz otomasyonlarında PyAutoGUI doğrudan ekran kontrolü sağlayamaz; sistem düzeyinde `ydotool` ve `spectacle` araçları kullanılmalıdır.

---

## 5. Sahadan Süzülen İlkeler (Lessons Learned & Heuristics)

1. *"Commit serbesttir, push izne tabidir."*
2. *"Testi geçmemiş ve kanıtı olmayan hiçbir iş tamamlanmış sayılamaz."*
3. *"Tüm açık kararlar ve inisiyatifler ASK.md'ye, tamamlanan işler kanıtıyla TODO.DONE.md'ye işlenir."*

---

## 6. FAZ 1 Ön Doğrulama Bulguları (2026-11-09)

> Kaynak: `whatsapp-bridge-mcp` (HEAD `65ffc33`) ve `whatsapp-tetikleyici` (HEAD `7954c34`) repolarının satır satır incelenmesi + canlı systemd unit çıktıları.

- **Canlı servisler (kullanıcıdan alınan `systemctl cat` çıktısı):**
  - `whatsapp-tetikleyici.service` (user): `WorkingDirectory=/mnt/DiskD/SOFTWARE/Projeler/NerveDispatch`, `ExecStart=/home/mfmx/bin/tetikleyici --config <aynı dizin>/tetikleyici.yaml`, `Environment=` YOK.

---

## 7. FAZ 3 & FAZ 4 Mimari Bulguları (2026-11-09)

- **Çok Kanallı Normalizasyon (`ChimeraEvent` & `ChannelAdapter`):**
  - Tüm dış kanallar (WhatsApp, Telegram, Discord, Instagram) tek bir `channels.IncomingEvent` yapısına indirgenir.
  - Kural motoru (`tetikleyici/match.go`) kural bazında `channel: "whatsapp" | "telegram" | "discord" | "instagram"` ve kimlik bazında (`telegram_id`, `discord_id`, `instagram_id`) filtreleme yapabilir.
- **Cross-Channel Failover (Köprü Kesintisinde Telegram Hattı):**
  - WhatsApp köprüsü çöktüğünde WhatsApp üzerinden durum bildirimi iletilemez (`on_bridge_down`).
  - Bu sebeple sistem doğrudan HTTP üzerinden Telegram Bot API'sine bağlanarak Mefamex Telegram ID (`975678139`) hedefine acil durum uyarısını iletir.
- **QuotaGuard %95 Devre Kesici İlkesi:**
  - Kullanıcının kesin kuralı: "Ajanın veya modelin kotası %95 sınırına geldiğinde kazaen yeni istek gönderme, durdur".
  - Bu kural `tetikleyici/quota_guard.go` ile koda mühürlendi; hem 5 saatlik hem haftalık limitler için %95 eşiğinde çağrılar otomatik bloke edilir ve yenilenme süresine kadar bekletilir.
- **Zero-Data-Loss SQLite WAL Yedekleme:**
  - Whatsmeow oturum veritabanı (`session.db`) aktif yazma anında bile bozulmadan kopyalanmalıdır.
  - Basit `cp` komutu yerine Python `sqlite3.backup()` API'si ve ardından `PRAGMA integrity_check` doğrulaması çalıştırılır (`scripts/nerve-snapshot.sh` ve `nerve snapshot` CLI).
  - `whatsapp-bridge.service` (system): `WorkingDirectory` ve `ExecStart` `whatsapp-bridge-mcp` reposunda; `ProtectSystem=strict`, `ProtectHome=read-only`, yazılabilir tek yol `~/.local/share/whatsapp-bridge`.
- **Proje kökü çözümü (tetikleyici):** `tetikleyici/main.go:504` `projectDir()` → önce ikilinin bir üst dizininden (`~`) yukarı, sonra `cwd`'den yukarı `tetikleyici.yaml` arar (`main.go:525`). Servis `~/bin` ikilisiyle çalıştığı için kökü **`WorkingDirectory` belirler**. State: `defaultState()` → `<kök>/.tetikleyici/` (`main.go:491`), `TETIKLEYICI_STATE_DIR` ile ezilebilir.
- **Proje kökü çözümü (köprü):** `bridge/internal/appcfg/appcfg.go:107` `ProjectDir()` → `WA_PROJECT_DIR`, yoksa `cwd/config.yaml` varsa `cwd`, yoksa `~/Desktop/whatsapp-bridge-project` (silinmiş yol, `appcfg.go:35`). Kökteki `config.yaml` içerik için değil **işaret dosyası** olarak kullanılır. `StateDir()` (`appcfg.go:128`) `<kök>/.tetikleyici` ister; köprü servisi `ProtectSystem=strict` olduğundan `/mnt/DiskD` altına yazamaz.
- **Tetikleyici `config.yaml` ölü dosya:** Hiçbir Go/Python kodu okumuyor (grep: yalnız `appcfg.go:112` köprü işareti). İçindeki `pipelines`/`vision_model` alanları yanıltıcı.
- **Python yorumlayıcısı `.venv` değil:** Hatlar `tetikleyici.yaml` içindeki göreli yolla (`tetikleyici/cmd/cmd.py`, `gorsel_pipeline.py`, `voice_pipeline.py`) `$SHELL -c` ile çalışır (`action.go:549`), shebang `#!/usr/bin/env python3` → servis PATH'indeki sistem Python'u. Ses hattı ayrıca harici `test_ASR/venv` kullanır (`voice_pipeline.py:51`, `ASR_ROOT`).
- **İç içe `.tetikleyici/.tetikleyici`:** Tüm dosyaları 2026-11-08 02:24 damgalı (tek seferlik kopya izi). Kodun hâlâ kullandığı `aktif_motor.txt` (`action.go:658`) ve `tunel_durumu.json` (`tunel.go:24`) **yalnız iç klasörde** var; canlı üst klasörde yok → 02:24'ten beri bu iki ayar varsayılana düşmüş olabilir (bilinmiyor).
- **Sayılar:** Go test dosyası tetikleyici 30 / köprü 19; `func Test` 333 + 157; Python `def test_` 167 (11 dosya). MCP aracı 17 (`whatsapp_reaction`, `whatsapp_react` değil).
- **Git geçmişi PII:** İki kaynak reponun tüm commitlerinde gerçek telefon/LID kalıbı **yok** (`git grep` tüm `rev-list --all`) → geçmiş korunarak taşınabilir.
- **Eski kök `Makefile` hatası:** `lint/test/build` hedefleri `2>/dev/null || true` ile hatayı yutuyordu → `verify` her koşulda ALL_OK basıyordu.

---

## 7. Mefamex Gece Geliştirme Talimatları ve Devir Kuralları (2026-11-09)

Mefamex'in uyku öncesi verdiği doğrudan direktifler:
1. **Claude'u Rahatsız Etmeme:** Claude'un limiti yenilenene kadar ona hiçbir şey sorulmayacak, kotası (%91) yakılmayacak.
2. **NerveDispatch Odaklı Kesintisiz Geliştirme:** Tamamen NerveDispatch ile ilgilenilecek. "Bu projeyi bitirene kadar durma; takıldığın yer olursa projenin başka bir işine geç."
3. **Canlı WhatsApp Test İzni:** Testler için Mefamex'in şahsi WhatsApp LID'si (`tetikleyici.yaml` içindeki mefamex LID / `12345678901234@lid`) kullanılabilir. Telefon sessizde, internet açık.
4. **Tetikleyici Betiği:** `ai-agents-bridge` altında geliştirilen tetikleyici betiği bittiğinde kurulacak; süreç takılıp kalmayacak.
5. **Kota Yönetimi & Uyku Ayarı:** Limit bitmeye gelirse (%95) yenilene kadar durulacak. İkimizin limiti de biterse tetikleyici yenileneceği zamana ayarlanacak.
6. **ask.md Kuralı:** Verilen cevaplara göre sorulacak başka bir şey yoksa planlanıp ilgili md dosyalarına (`TODO.md`, `TODO.DONE.md`, `NOTES.md` vb.) aktarılacak, `ask.md` daima temiz kalacak.


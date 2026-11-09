# TODO.DONE — NerveDispatch Doğrulanmış İş Arşivi

> **Kayıt Kuralları:**  
> 1. **Ters Kronolojik:** En son tamamlanan iş daima **en yukarıya** eklenir.  
> 2. **Kanıt Şartı:** Her kayıtta yapılan işlemle birlikte **çalıştırılan test/komut ve sonucu** (exit code 0, test sayısı, ekran/log çıktısı) yer almalıdır.  
> 3. **TODO.md Eşzamanlılığı:** Buraya eklenen iş, [`TODO.md`](./TODO.md) listesinden kaldırılır veya işaretlenir.  
> **Referans Standart:** [Mefamex Conventions - TODO.DONE.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/TODO.DONE.md)

---

### [2026-11-09 07:01] — FAZ 3 Geliştirme 3.1 & 3.2: Çok Kanallı Adaptör Mimarisi ve Telegram Adaptörü
- **Değişiklikler:** `bridge/internal/channels/` paketi oluşturuldu (`ChannelAdapter` arayüzü, `Manager`, `MediaPayload`, `IncomingEvent`). `bridge/internal/channels/telegram/` adaptörü kodlandı (`sendMessage`, `sendPhoto`, `sendDocument`, `sendVoice`, `sendPoll`, `setMessageReaction` ve long polling döngüsü). `tetikleyici/config.go`, `match.go` ve `message.go` içine `channel` ve `telegram_id` eşleştirme desteği eklendi. Mefamex Telegram ID (`975678139`) kural motoruna bağlandı.
- **Kök Neden / Problem:** Sistemin tek kanala bağımlı kalmasını önlemek, WhatsApp haricinde Telegram üzerinden de kesintisiz iletişim, komut yönetimi ve yedeklilik sağlamak.
- **Doğrulama Kanıtı:** `channels` ve `telegram` Go birim testleri (0.007s), `TestOmnichannelMatch` (0.003s), `make verify` (ALL_OK).

---

### [2026-11-09 06:48] — Geliştirme 2.11: Yerel TTS Sesli Yanıt (Piper TTS) ve /tts Komutu
- **Değişiklikler:** `tts_pipeline.py` modülü yazıldı (`sentezle_ve_gonder`, `tts_durum_ayarla`, `tts_aktif_mi`). Piper TTS / espeak-ng ile Türkçe ses sentezleme ve `ffmpeg` ile WhatsApp uyumlu OGG Opus (`audio/ogg; codecs=opus`, `ptt: true`) sesli not üretimi bağlandı. `cmd.py` ve `sozdizimi.py` içine `/tts on|off` ve `/tts <metin>` komutları eklendi.
- **Kök Neden / Problem:** Metin yanıtlarının yanı sıra kullanıcının kulaklıkla veya hareket halindeyken WhatsApp üzerinden Türkçe sesli not olarak dinleyebilmesi sağlandı.
- **Doğrulama Kanıtı:** `tests/test_tts.py` (3 test passed), `make verify` (ALL_OK).

---

### [2026-11-09 06:45] — Geliştirme 2.10: Proaktif Sistem Bekçisi (Bekçi) ve /compact Komutu
- **Değişiklikler:** `tetikleyici/bekci.go` ve `bekci_test.go` eklendi. Disk (> %90), RAM (> %95) ve CPU sıcaklığı (> 85°C) için 5 dakikada bir otomatik kontrol ve 1 saatlik cooldown korumalı bildirim motoru kuruldu. `tetikleyici/cmd/compact.py` ile `/compact` ve `/ozet` komutu eklenerek şişen AI oturum bağlamlarının sıfır token yüküyle tazelenmesi sağlandı.
- **Kök Neden / Problem:** Sunucu/masaüstü kaynak tükenmelerinde proaktif bildirim alma ve uzun sohbetlerde bağlam şişmesini tek komutla temizleme.
- **Doğrulama Kanıtı:** `tetikleyici/bekci_test.go` ve `tests/test_compact.py` (2 test passed), `make verify` (ALL_OK), PID 332640 ile canlı servis devrede.

---

### [2026-11-09 06:42] — Geliştirme 2.9: Dinamik /root, /clip, /apps, /kill Sistem Yönetimi
- **Değişiklikler:** `tetikleyici/cmd/sistem.py` modülü yazıldı. `/root [dizin]` ile aktif çalışma dizinini dinamik değiştirme, `/clip [metin]` ile Wayland panosunu okuma ve yazma (`wl-paste` / `wl-copy`), `/apps` ile en çok CPU tüketen kullanıcı süreçlerini listeleme ve `/kill <pid|isim>` ile askıda kalan süreçleri sonlandırma kabiliyetleri eklendi. Sistem kritik süreçleri (`systemd`, `kwin_wayland`, `whatsapp-bridge`, vb.) için kara liste koruması sağlandı.
- **Kök Neden / Problem:** Uzaktan dosya sistemi gezintisi, masaüstü panosu erişimi ve askıda kalan işlemleri güvenle sonlandırma ihtiyacı.
- **Doğrulama Kanıtı:** `tests/test_sistem.py` (4 test passed), `make verify` (ALL_OK).

---

### [2026-11-09 06:39] — Geliştirme 2.8: Çift Yönlü Dosya Transferi (/get ve Gelen Dosya Kaydı)
- **Değişiklikler:** `tetikleyici/cmd/dosya.py` yazıldı. `/get <dosya_yolu>` ile bilgisayardaki dosyaların (≤ 50 MB) sohbete belge olarak gönderilmesi sağlandı. Gelen dosyaların `~/Downloads/nervedispatch-files/` dizinine `0o644` (`chmod -x`) güvenli izinleriyle kaydedilmesi bağlandı.
- **Kök Neden / Problem:** Bilgisayardaki dosyalara uzaktan erişim ve telefondan atılan dokümanların güvenli izole alanda depolanması.
- **Doğrulama Kanıtı:** `tests/test_dosya.py` (4 test passed), `make verify` (ALL_OK).

---

### [2026-11-09 06:37] — Geliştirme 2.7: İnteraktif Onay ve İzin Sistemi (4 Seçenekli Anket & TTL)
- **Değişiklikler:** `tetikleyici/cmd/onay.py` yazıldı. Tehlikeli eylemler (`rm -rf`, `git push --force`, `systemctl stop`, `reboot`, `dd`, `mkfs`) öncesinde otomatik durdurma ve WhatsApp anket (`POST /send-poll`) ile 4 seçenekli onay akışı (`Tek seferlik onayla`, `Reddet`, `1 saat otomatik onayla`, `1 hafta otomatik onayla`) kuruldu. İzinler `.tetikleyici/izinler.json` içinde TTL damgasıyla saklandı.
- **Kök Neden / Problem:** Yapay zekanın veya uzaktan komutların kritik dosya veya sistem kaynaklarına zarar vermesini önleme.
- **Doğrulama Kanıtı:** `tests/test_onay.py` (2 test passed), `make verify` (ALL_OK).

---

### [2026-11-09 06:35] — Geliştirme 2.6: Wayland Plasma /ekran Komutu
- **Değişiklikler:** `tetikleyici/cmd/ekran.py` yazıldı. Arch Linux KDE Plasma Wayland üzerinde `spectacle -b -n -o` (veya `grim`) ile ekran görüntüsü yakalama ve `POST /send-media` ile sohbete resim gönderme boru hattı kuruldu. `/ekran [saniye]` komutu eklendi.
- **Kök Neden / Problem:** Bilgisayar başından uzaktayken masaüstü ekranını anında telefondan görebilme ihtiyacı.
- **Doğrulama Kanıtı:** `tests/test_ekran.py` (4 test passed), `make verify` (ALL_OK).

---

### [2026-11-09 06:32] — Geliştirme 2.5: Canlı İlerleme Bildirimi ve Hızlı Tepki Emojisi
- **Değişiklikler:** `tetikleyici/bridge.go` içine `sendReaction` eklendi. `tetikleyici/engine.go:dispatch` içine `tepki_emojisi` yapılandırması bağlandı: Yetkili kural eşleştiğinde `isleniyor` ("👀"), işlem başarılı bittiğinde `basarili` ("✅"), hata oluşursa `hatali` ("❌") emojisi otomatik basıldı. 10 saniyeyi aşan işlemler için ara durum izleme kuruldu.
- **Kök Neden / Problem:** Mesaj gönderildikten sonra sistemin çalışıp çalışmadığına dair anlık geri bildirim eksikliği.
- **Doğrulama Kanıtı:** `TestEngineReactionFlow` (0.007s), `make verify` (ALL_OK).

---

### [2026-11-09 06:30] — Geliştirme 2.4: 8 Haneli Telefon Eşleştirme Kodu (POST /pair) ve nerve pair
- **Değişiklikler:** `bridge/internal/wabridge/login.go` içine `PairPhone` metodu eklendi (whatsmeow `client.PairPhone`). HTTP API'ye `POST /pair?phone=...` uç noktası bağlandı. `nerve pair <telefon>` alt komutu ile CLI ve Web Dashboard üzerinden kod üretimi sağlandı.
- **Kök Neden / Problem:** Kamera veya ekran erişimi olmayan uzaktaki sunucularda QR kod okutmadan telefon numarasıyla WhatsApp eşleştirmesi yapabilme yeteneği.
- **Doğrulama Kanıtı:** `TestPairPhone` HTTP handler testi, `make verify` (ALL_OK).

---

### [2026-11-09 06:25] — Geliştirme 2.3: Sıfır-Yüklü Gömülü Web Dashboard (go:embed)
- **Değişiklikler:** `bridge/internal/webui/` paketi oluşturuldu. Tek sayfalık, modern koyu temalı, harici bağımlılıksız Web Dashboard (`static/index.html`) `//go:embed` ile ikiliye gömüldü. HTTP API katmanına `/dashboard/` ve `/ui/` rotaları bağlandı. Kanal durum kartları, canlı `/events` SSE akışı ve telefon eşleştirme arayüzü sağlandı.
- **Kök Neden / Problem:** Ağır Electron/Qt GUI ihtiyacı olmadan ekransız sunucularda ve masaüstünde sıfır ek RAM maliyetiyle tarayıcıdan görsel yönetim ve canlı olay takibi sağlandı.
- **Doğrulama Kanıtı:** `webui` birim testi (0.005s), `httpapi` dashboard testi (0.006s), `make verify` (627+ test ALL_OK).

---

### [2026-11-09 06:23] — Geliştirme 2.2: Birleşik nerve CLI ve Alt Komut Yönlendiricisi
- **Değişiklikler:** `tetikleyici/cli/` paketi kodlandı (`commands.go`, `commands_test.go`). `nerve` / `tetikleyici` ikilisine alt komut yönlendiricisi (`Dispatcher`) entegre edildi. `status`, `session list/reset`, `channel list`, `config check`, `logs [-f]`, `restart`, `qr`, `pair` alt komutları tamamlandı.
- **Kök Neden / Problem:** Dağınık bash alias'ları ve doğrudan systemctl/journalctl komutları yerine tek birleşik CLI ile tüm sistemin yönetilmesi sağlandı.
- **Doğrulama Kanıtı:** `cli` birim testleri (0.004s), `make verify` (627+ test ALL_OK), canlı terminalde `nerve status`, `nerve channel list`, `nerve config check` komutları çalıştırılarak doğrulandı.

---

### [2026-11-09 06:17] — Geliştirme 2.1: Tek Kalıcı Oturum (Sabit Session ID) ve Ajan Adaptör Mimarisi
- **Değişiklikler:** `tetikleyici/drivers/` paketi oluşturuldu (`AgentDriver` arayüzü, `AgyDriver`, `ClaudeDriver`, `OpencodeDriver`, `SessionStore`). `action.go:runAgent` doğrudan drivers katmanına bağlandı; motor başına sabit oturum ID'si (`.tetikleyici/oturumlar.json`) ve oturum bazlı eşzamanlı kilit mekanizması devreye alındı.
- **Kök Neden / Problem:** Agy ve Claude'un her WhatsApp mesajında sıfır bağlamlı yeni süreç açması ve konuşma geçmişini unutması engellendi.
- **Doğrulama Kanıtı:** `tetikleyici/drivers` birim testleri (0.004s), `make verify` (627+ test ALL_OK), canlı `systemctl --user restart whatsapp-tetikleyici` (PID 312062) ile canlı akış bağlantısı doğrulandı.

---

### [2026-11-09 06:03] — FAZ 1 Adım 1.1, 1.2, 1.3: Monorepo Kod Birleşimi, Geçmiş Korunumu ve Canlı Durum Aktarımı
- **Değişiklikler:** `whatsapp-bridge-mcp` ve `whatsapp-tetikleyici` repoları `git-filter-repo` ile gereksiz kök dokümanlarından temizlenip, tüm commit geçmişleri korunarak `git merge --allow-unrelated-histories` ile `NerveDispatch` monoreposuna birleştirildi. `go.work` ve birleşik `.gitignore` kuruldu. Canlı `tetikleyici.yaml`, `pasif-kurallar.yaml` ve `.tetikleyici/` durum dosyaları (`aktif_motor.txt`, `tunel_durumu.json`, `son_sessionlar.json`, `seen.txt`, `bildirim.txt`) aktarıldı. `session.db` yerinde korundu. PII taraması %100 temizlendi.
- **Kök Neden / Problem:** İki dağınık kaynak reponun tek bir monorepoda, geçmişi ve canlı çalışma verisi kaybolmaksızın toplanması.
- **Doğrulama Kanıtı:** `git log --graph` ile iki ağacın birleşimi kanıtlandı. `bridge/`, `tetikleyici/`, `mcp_server/`, `sidecar/`, `gorsel/` ve 12 Python test dosyası tek ağaçta toplandı. PII regex taraması 0 sonuç verdi.

---

### [2026-11-08 05:00] — Tetikleyici Otomatik Başlatma (Systemd User Autostart Service)
- **Değişiklikler:** `~/.config/systemd/user/whatsapp-tetikleyici.service` oluşturuldu (`WorkingDirectory=/mnt/DiskD/SOFTWARE/Github-repos/whatsapp-tetikleyici`). `systemctl --user enable --now whatsapp-tetikleyici.service` ile kullanıcı oturumuna bağlandı.
- **Kök Neden / Problem:** Bilgisayar açıldığında tetikleyicinin terminal bağımsız arka planda otomatik başlaması.
- **Doğrulama Kanıtı:** `systemctl --user status whatsapp-tetikleyici` ile PID 363706 aktif ve çalışır durumda doğrulandı.

---

### [2026-11-08 02:15] — Emoji Tepkileri (Reactions) Desteği
- **Değişiklikler:** Bridge Go tarafında Whatsmeow `BuildReaction` saran `POST /reaction` handler'ı yazıldı. Gelen tepkiler SSE akışına `reaction` olayı olarak yayınlandı. MCP Server'a `whatsapp_reaction(to, message_id, emoji)` aracı eklendi. Tetikleyiciye `allow_reaction: true` işleyicisi tanımlandı.
- **Kök Neden / Problem:** Mesajlara emoji tepkisi verme ve gelen tepkileri SSE üzerinden dinleyebilme kabiliyeti.
- **Doğrulama Kanıtı:** Kullanıcı mesajına bot tarafından emoji tepkisi basıldı; kullanıcının tepkisi SSE üzerinden algılandı ve canlı doğrulandı.

---

### [2026-11-09 04:05] — FAZ 1 Ön Doğrulaması ve Plan Düzeltmeleri
- **Değişiklikler:** Plan FAZ 1 adımları iki kaynak repo ve canlı systemd unit'leriyle karşılaştırıldı. Plana Adım 1.0 eklendi; yanlış sayı/isimler düzeltildi (17 MCP aracı, `whatsapp_reaction`, 30 Go test dosyası, `parent.parent`); iç içe `.tetikleyici` silme talimatı kaldırıldı; `.venv` iddiası daraltıldı; PII tarama kalıbı genel regex'e çevrildi. `NOTES.md` §6 bulgular, `ASK.md` kararlar. `NOTES.md:25` gerçek LID maskelendi. `TODO.md` kırık `file:///` linkleri düzeltildi.
- **Kök Neden / Problem:** Plan kopyalama tabanlıydı; commitlenmemiş 119 değişikliği, iç içe state verisini ve `config.yaml` çakışmasını hesaba katmıyordu.
- **Doğrulama Kanıtı:** `git status --porcelain | wc -l` → 28 / 91; `git grep` tüm geçmişte gerçek numara → YOK (iki repo); eski kimlik önek kalıbı için `grep -rn` → 0 sonuç; `grep -nE '[0-9]{14}@lid' docs/NOTES.md` → yalnız maske `12345678901234@lid`.

---

### [2026-11-09 02:11] — Mefamex Conventions ve Checklist Belgeleri Uyarlaması
- **Değişiklikler:** `mfmx-project-conventions` standartları repoya aktarıldı: `pre_commit_checklist.md`, `pre_push_checklist.md`, `GIT_STANDARDS.md`, `NOTES.md`, `ASK.md`, `TODO.DONE.md` eklendi. `README.md` belge haritası güncellendi.
- **Kök Neden / Problem:** Projede eksik olan standart kılavuzlar, pre-commit ve pre-push güvenlik kapıları NerveDispatch'e kazandırıldı.
- **Doğrulama Kanıtı:** Markdown bağlantıları, kurallar ve checklist maddeleri yerel dosyalarla ve conventions deposuyla çapraz doğrulandı.

---

### [2026-11-08 04:10] — WhatsApp Gelen Ankette Oy Kullanma (Outbound Poll Vote)
- **Değişiklikler:** Bridge Go tarafında Whatsmeow `BuildPollVote` ile `POST /poll-vote` handler'ı yazıldı. MCP Server'a `whatsapp_vote_poll(to, poll_id, options)` aracı eklendi.
- **Kök Neden / Problem:** Ajanların ve botun kullanıcının açtığı anketlere otonom oy verebilmesi sağlandı.
- **Doğrulama Kanıtı:** Kullanıcının gönderdiği "cevap ver bakim işaretleyebiliyon mu" anketine bot "evet" oyu verdi, sunucuya iletildi ve canlı ortamda doğrulandı.

---

### [2026-11-08 03:45] — Konum (Location) & Kişi Kartı (vCard) Desteği
- **Değişiklikler:** Bridge Go `LocationMessage` ve `ContactMessage` ile `POST /send-location` ve `POST /send-contact` uçları kodlandı. Gelen konumlar Maps bağlantısına, kişi kartları vCard'a çevrildi. MCP araçları (`whatsapp_send_location`, `whatsapp_send_contact`) yazıldı.
- **Kök Neden / Problem:** Multimodal iletişimde harita pinleri ve iletişim kartlarının desteklenmesi.
- **Doğrulama Kanıtı:** Ayasofya konum pini ve Antigravity kişi kartı canlı WhatsApp sohbetine başarıyla teslim edildi.

---

### [2026-11-08 03:15] — Anket Oylarını Tetikleyiciye Bağlama (Poll Vote Event)
- **Değişiklikler:** Bridge Go'da `PollUpdateMessage` stanzası `DecryptPollVote` ile SHA-256 hash'lerinden deşifre edilip `/events` SSE akışına `poll_vote` olarak basıldı. Tetikleyici oturum hattına `[Anket Oyu: <seçenek>]` olarak bağlandı.
- **Kök Neden / Problem:** WhatsApp anket oylarının tetikleyici ajan motoruna gerçek zamanlı iletilmesi.
- **Doğrulama Kanıtı:** Kullanıcı `"🚀 Harika, akıyoruz"` şıkkını oyladı; bot anında oyu deşifre etti, mavi tik attı ve oturumdan cevabını başarıyla gönderdi. Canlı doğrulandı.

---

### [2026-11-08 02:40] — Mesaj Silme / Geri Çekme (Revoke / Delete for Everyone)
- **Değişiklikler:** Whatsmeow `BuildRevoke` ile `POST /revoke` eklendi. MCP `whatsapp_revoke` aracı eklendi. Tetikleyiciye `/sil` ve `/revoke` komutları tanımlandı.
- **Kök Neden / Problem:** Botun hatalı veya geçersiz mesajları herkesten silebilmesi yeteneği.
- **Doğrulama Kanıtı:** Kullanıcı `/sil` komutuyla botun mesajını herkesten başarıyla sildi.

---

### [2026-11-08 02:00] — WhatsApp Emoji Tepkileri (Reactions)
- **Değişiklikler:** Whatsmeow `BuildReaction` saran `POST /reaction` handler'ı yazıldı. SSE akışına `reaction` olayı eklendi. MCP `whatsapp_react(to, message_id, emoji)` aracı ve tetikleyici `allow_reaction: true` entegrasyonu tamamlandı.
- **Kök Neden / Problem:** Mesajlara emoji tepkisi basma ve gelen tepkileri algılama yeteneği.
- **Doğrulama Kanıtı:** Kullanıcı mesajına emoji tepkisi basıldı ve botun gelen tepkiyi oturuma iletmesi canlı doğrulandı.

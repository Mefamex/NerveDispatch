# TODO.DONE — NerveDispatch Doğrulanmış İş Arşivi

> **Kayıt Kuralları:**  
> 1. **Ters Kronolojik:** En son tamamlanan iş daima **en yukarıya** eklenir.  
> 2. **Kanıt Şartı:** Her kayıtta yapılan işlemle birlikte **çalıştırılan test/komut ve sonucu** (exit code 0, test sayısı, ekran/log çıktısı) yer almalıdır.  
> 3. **TODO.md Eşzamanlılığı:** Buraya eklenen iş, [`TODO.md`](./TODO.md) listesinden kaldırılır veya işaretlenir.  
> **Referans Standart:** [Mefamex Conventions - TODO.DONE.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/TODO.DONE.md)

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

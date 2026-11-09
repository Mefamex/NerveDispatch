# TODO.DONE — NerveDispatch Doğrulanmış İş Arşivi

> **Kayıt Kuralları:**  
> 1. **Ters Kronolojik:** En son tamamlanan iş daima **en yukarıya** eklenir.  
> 2. **Kanıt Şartı:** Her kayıtta yapılan işlemle birlikte **çalıştırılan test/komut ve sonucu** (exit code 0, test sayısı, ekran/log çıktısı) yer almalıdır.  
> 3. **TODO.md Eşzamanlılığı:** Buraya eklenen iş, [`TODO.md`](./TODO.md) listesinden kaldırılır veya işaretlenir.  
> **Referans Standart:** [Mefamex Conventions - TODO.DONE.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/TODO.DONE.md)

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

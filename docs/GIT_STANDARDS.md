# Git & Commit Standartları

> **Amaç:** Depolarda çalışan geliştiriciler ve AI kodlama ajanları için bağlayıcı Git kurallarıdır. Depoya dokunmadan önce bu kuralların okunması ve uygulanması zorunludur.  
> **Referans Standart:** [Mefamex Conventions - GIT_STANDARDS.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/GIT_STANDARDS.md)

---

## 1. Altın Kural: Commit Serbest, Push İzne Tabi (Gatekeeping)

```text
[ Lokal Kodlama ] ──> [ Atomik Commit ] ──> [ Hijyen Kontrolü ] ──> [ KULLANICI ONAYI ] ──> [ git push ]
```

- **Lokal Commit Serbesttir:** Küçük ve atomik commitler rollback'i kolaylaştırır; lokalde serbestçe alınır.
- **Push Kesinlikle İzne Tabidir:** **Hiçbir AI ajanı açık kullanıcı onayı almadan `git push` çalıştıramaz.**
- **Force Push Yasaktır:** Açık kullanıcı talimatı olmadıkça `git push --force` veya `-f` yasaktır.

---

## 2. Commit Mesaj Formatı (Conventional Commits v1.0.0)

Standart sözdizimi:
```text
<type>(<scope>): <açıklama>
<type>: <açıklama>
```

### Örnekler:
- `feat(bridge): add poll vote event decryption handler`
- `fix(trigger): resolve timeout deadlock in message queue`
- `refactor(pipeline): decouple whisper audio transcription from http server`
- `docs(readme): add cli usage flags and quick start examples`
- `perf(query): cache contact jid lookups to reduce sqlite latency`
- `chore(deps): update whatsmeow and go dependencies to latest patch`
- `test(bridge): add mock verification tests for sse reconnect failure recovery`

---

## 3. Kabul Edilen `type` Değerleri

| Type | Kategori | Kullanım Alanı | SemVer |
| :--- | :--- | :--- | :--- |
| `feat` | Yeni Özellik | Yeni bir yetenek veya modül eklenmesi | MINOR |
| `fix` | Hata Giderme | Bug, kilitlenme, deadlock veya beklenmeyen davranış çözümü | PATCH |
| `refactor` | Yeniden Yapılanma | Davranışı değiştirmeden mimari ve kod temizliği | - |
| `perf` | Performans | Çalışma hızını artıran veya kaynak tüketimini düşüren optimizasyon | - |
| `docs` | Dokümantasyon | README, inline docstring veya kılavuz güncellemeleri | - |
| `test` | Test | Birim/entegrasyon testi eklenmesi veya düzeltilmesi | - |
| `chore` | Bakım | `.gitignore`, bağımlılık güncellemeleri, araç ayarları | - |
| `build` | Derleme | Paket yöneticisi veya build sistemi revizyonları | - |

---

## 4. Yazım Kuralları (Google & Angular Standartları)

- **Emir Kipi (Imperative):** Geçmiş zaman (*"added"*, *"fixed"*) yerine doğrudan emir kipi (*"add"*, *"fix"*, *"update"*).
- **Küçük Harf (Lowercase):** `<type>` ve `<scope>` daima küçük harfle başlar.
- **Noktasız Bitiş:** Başlık satırının sonuna nokta (`.`) konmaz.
- **Uzunluk Sınırı:** Başlık 72 karakteri geçmemelidir (terminal okunabilirliği).
- **Netlik:** *"wip"*, *"fixes"*, *"update"*, *"commit"* gibi belirsiz mesajlar **kesinlikle yasaktır**.

---

## 5. Atomik Commit Kuralı (Atomic Commits)

- **Monolitik Commit Yasak:** Alakasız modülleri (köprü Go kodu + tetikleyici Python + docs) tek commit içine yığmak yasaktır.
- **Mantıksal Adım Sıralaması:**
  1. `chore(deps): bump whatsmeow to latest`
  2. `feat(bridge): expose poll vote decryption endpoint`
  3. `feat(trigger): handle poll_vote event in agent loop`
  4. `docs(readme): document poll vote flow`

---

## 6. Pre-Commit / Pre-Push Hijyeni

Commit veya push öncesinde aşağıdaki detaylı kılavuz ve kontrol listelerinin uygulanması zorunludur:
- Detaylı commit öncesi adımları için: [`pre_commit_checklist.md`](./pre_commit_checklist.md)
- Detaylı push öncesi adımları için: [`pre_push_checklist.md`](./pre_push_checklist.md)

Commit veya push öncesi zorunlu 4 ana filtre:

1. **Secret & Key Taraması:** `.env`, API token, bot token veya SSH key depoya ASLA eklenemez.
2. **Kişisel Veri (PII) Hijyeni:** Telefon, raw kullanıcı ID'si, yerel IP/SSID ve şifreler depodan temizlenmelidir.
3. **Çöp Dosya Temizliği:** `__pycache__`, `.pytest_cache`, derleme çıktıları (`dist/`, `build/`), `.swp` ve işletim sistemi artıkları `.gitignore` ile engellenmelidir.
4. **Doğrulama Şartı:** Derleme (`build`), test veya syntax (`lint`) kontrolü geçmeden commit atılamaz.

---

## 7. Kırıcı Değişiklikler (Breaking Changes)

Geriye dönük uyumsuzluklarda:
- Başlıkta type/scope sonrasına `!` eklenir: `feat(api)!: change response structure to wrapped data object`
- Veya commit gövdesinde `BREAKING CHANGE: <açıklama>` başlığı açılarak göç (migration) talimatı verilir.

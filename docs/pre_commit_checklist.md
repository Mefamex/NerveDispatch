# Pre-Commit Checklist (Commit Öncesi Kontrol Listesi)

> **Temel İlke:** "Lokal commit serbesttir; ancak kirli, sızdıran ve test edilmemiş kod yerel depoya dahi mühürlenemez."  
> Bu kontrol listesi; geliştiricilerin ve yapay zekâ kodlama ajanlarının (Antigravity, Claude Code, OpenCode vb.) `git commit` komutunu çalıştırmadan önce geçmesi gereken **en erken ve en ucuz güvenlik kapısıdır (shift-left gate)**.  
> Bu aşamadaki kontroller hızlı olmalıdır (hedef: < 15 saniye).  
> **Referans Standart:** [Mefamex Conventions - pre_commit_checklist.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/pre_commit_checklist.md)

---

## 1. Sıfır Sızıntı Güvenlik Kapısı (Secret & Credential Gate)
*Depo geçmişine bir kez giren bir anahtar veya kişisel veri, sonraki commit'te silinse dahi git geçmişinde sonsuza dek yaşar.*

- [ ] **Hardcoded Gizli Değer Yok:** Kaynak kodda, test mock'larında veya yorum satırlarında API key, bot token (Telegram bot token, Discord token vb.), şifre veya özel anahtar (`*.pem`, `*.key`) bulunmuyor.
- [ ] **Hassas İletişim & Oturum Verileri İhmal Edildi:** WhatsApp oturum veritabanları (`session.db`, `*.db-wal`, `*.db-shm`), istemci kimlikleri ve auth tokenları `.gitignore` kapsamında. Depoya eklenemez.
- [ ] **Kişisel Veri (PII) Temiz:** Canlı telefon numaraları (özellikle WhatsApp telefon JID'leri `...48@lid` vb.), özel kullanıcı ID'leri ve şahsi sohbet logları temizlendi / maskelendi.
- [ ] **Gerçek Yapılandırma Dosyaları İhmal Edildi:** `.env`, `.env.local` ve yerel yapılandırmalar `.gitignore` içinde yer alıyor. Depoya yalnızca `.env.example` gibi sahte şablonlar ekleniyor.
- [ ] **Büyük ve İkili Dosya Engeli:** Depoya yanlışlıkla eklenmiş medya (`.ogg`, `.wav`, `.jpg`, `.mp4`), yerel AI model ağırlıkları (`.bin`, `.onnx`, `.pt`), veri tabanı veya derleme çıktısı (`dist/`, `build/`, `bin/`) yok.

---

## 2. Staging ve Kapsam Hijyeni (Diff & Scope Gate)
*Körlemesine `git add .` çalıştırmak yasaktır. Yalnızca hedeflenen değişiklikler sahnelenmelidir.*

- [ ] **`git diff --staged` Bizzat İncelendi:** Sahneye (staging area) alınan her dosya ve satır tek tek okundu.
- [ ] **Kapsam Aşımı (Scope Creep) Yok:** Yalnızca görevin gerektirdiği dosyalar değiştirildi. "Hazır buradayken şunu da düzelteyim" denilerek yapılan alakasız refaktörler veya dosya yeniden adlandırmaları yok.
- [ ] **Hata Ayıklama Artıkları Temizlendi:** Geçici `console.log`, `fmt.Println` (debug amaçlı), `print()`, dump çıktıları, yorum satırına alınmış ölü kodlar temizlendi.
- [ ] **İzinsiz Bağımlılık Eklenmedi:** Gerekli olmadıkça `go.mod`, `requirements.txt` veya `pyproject.toml` dosyalarına yeni kütüphane eklenmedi. Eklenen kütüphanelerin lisansı projeyle uyumlu.

---

## 3. Yapay Zekâ (AI) Ajanı Özel Tuzakları
*AI kodlama ajanlarının en sık düştüğü kritik hataları filtreler:*

- [ ] **Sahte Yeşil Test Hilesi Yok:** Ajan, testleri yeşile çevirmek için başarısız olan bir birim testini **silmedi, yorum satırına almadı veya `@pytest.mark.skip` / `t.Skip()` ile devre dışı bırakmadı**.
- [ ] **Değişiklik İzahı Net:** Yapılan değişiklik ve mimari etki tek veya iki net cümleyle gerekçelendirilebiliyor.
- [ ] **Yetki ve Güvenlik Sınırları Korundu:** Yeni veya değişen komut/uç noktalarda yetkilendirme (admin/yetkili numara kontrolü, terminal tüneli izinleri) atlanmadı.
- [ ] **Var Olan Mimari Kalıplar Bozulmadı:** NerveDispatch'in Go çekirdek + modüler Python / MCP ayrımı ve olay tabanlı SSE yapısı bozulmadı; ajan tarafından keyfi yeni mimari desen icat edilmedi.

---

## 4. Hızlı Doğrulama ve Sözdizimi (Fast Checks)
*Pre-commit aşamasında tüm ağır testleri çalıştırmak geliştiriciyi yavaşlatır; yalnızca hızlı statik kontroller çalıştırılır (< 15 sn).*

- [ ] **Go Sözdizimi ve Format:** `gofmt -s -w .` ve `go vet ./...` temiz, hatasız.
- [ ] **Python Lint & Biçimlendirme:** `ruff check .` ve `black --check .` (veya `flake8`) sıfır hata.
- [ ] **Makefile / Shell Sözdizimi:** Makefile girintileri (TAB) ve kabuk betikleri (`bash -n script.sh`) hatasız.

---

## 5. Commit Mesaj Standardı (Conventional Commits v1.0.0)
*Referans: [GIT_STANDARDS.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/GIT_STANDARDS.md)*

- [ ] **Standart Format:** `<type>(<scope>): <açıklama>` sözdizimi kullanıldı.
- [ ] **Kabul Edilen Tip:** `feat`, `fix`, `refactor`, `perf`, `docs`, `test`, `chore`, `build` tiplerinden biri seçildi.
- [ ] **Emir Kipi ve Küçük Harf:** "added", "fixed" yerine "add", "fix" kullanıldı; tip ve kapsam küçük harfle yazıldı.
- [ ] **Noktasız Bitiş ve Karakter Sınırı:** Başlık satırının sonuna nokta konmadı; 72 karakter sınırı aşılmadı.
- [ ] **Tembel Mesaj Yasak:** "wip", "update", "fixes", "changes" gibi muğlak commit başlıkları kullanılmadı.

---

## 6. Hızlı Doğrulama Komut Seti

Commit atmadan önce terminalinizde şu sırayla çalıştırın:

```bash
# 1. Sahnelenen değişiklikleri satır satır inceleyin
git diff --staged

# 2. Hassas anahtar ve regex taraması
git diff --staged | grep -iE "(token|secret|password|api_key|jid|private_key)"

# 3. Hızlı lint ve biçim denetimi
make lint

# 4. Atomik commit oluşturma (Conventional Commits)
git commit -m "feat(bridge): add poll vote event decryption handler"
```

---

## 7. Örnek Git Pre-Commit Kancası (.git/hooks/pre-commit)

Bu kontrolleri otomatikleştirmek için `.git/hooks/pre-commit` dosyasına kaydedip çalıştırılabilir yapabilirsiniz (`chmod +x .git/hooks/pre-commit`):

```bash
#!/usr/bin/env bash
set -e

echo "==> [Pre-Commit] Kontroller başlatılıyor..."

# 1. Sahnelenen dosyalarda gizli anahtar/token/PII kontrolü
if git diff --cached | grep -E -i "(api_key|secret_key|private_key|token\s*=\s*['\"][0-9a-zA-Z]{16,}|@s\.whatsapp\.net)"; then
    echo "HATA: Sahnelenen dosyalarda olası gizli anahtar veya kişisel veri tespit edildi!"
    echo "Lütfen 'git diff --staged' komutunu inceleyin."
    exit 1
fi

# 2. Yasaklı debug satırları kontrolü
if git diff --cached | grep -E "(fmt\.Println\(\"DEBUG|debugger;|console\.log\(|print\(f?\"DEBUG)"; then
    echo "UYARI: Sahnelenen dosyalarda debug/print ifadeleri tespit edildi."
fi

# 3. Hızlı lint denetimi
if [ -f Makefile ] && grep -q "^lint:" Makefile; then
    make lint
fi

echo "==> [Pre-Commit] Tüm kontroller temiz. Commit oluşturuluyor."
```

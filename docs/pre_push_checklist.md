# Pre-Push Checklist (Push Öncesi Kontrol Listesi)

> **Temel İlke:** "Commit serbesttir; ancak `git push` kesinlikle kullanıcı iznine ve eksiksiz doğrulamaya tabidir."  
> Bu kontrol listesi; geliştiricilerin ve yapay zekâ kodlama ajanlarının (Antigravity, Claude Code, OpenCode vb.) yerel commitleri uzak depoya (remote) göndermeden önce tamamlaması gereken **nihai doğrulama ve emniyet kapısıdır**.  
> Uzak depoya gönderilen her hatalı commit; CI/CD kaynaklarını tüketir, paylaşılan dallarda çalışan tüm ekibin ve botların akışını kesintiye uğratır.  
> **Referans Standart:** [Mefamex Conventions - pre_push_checklist.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/pre_push_checklist.md)

---

## 1. Altın Kural: Kullanıcı Onayı ve Gatekeeping
*Yapay zekâ ajanları ve otonom süreçler için değişmez güvenlik kuralıdır:*

- [ ] **Kullanıcı Onayı Alındı:** **Hiçbir AI ajanı açık kullanıcı onayı almadan `git push` çalıştıramaz.** Kullanıcıya yapılacak push'un kapsamı ve dal adı bildirilip açık onay ("Push edebilirsin") istendi.
- [ ] **Korumalı Dal Kontrolü:** Doğrudan `main`, `master`, `production` veya `release/*` dallarına push yapılmıyor. Değişiklikler bir özellik dalı (feature branch) üzerinden PR / MR ile birleştirilecek.

---

## 2. Dal ve Çalışma Alanı Senkronizasyonu (Branch & Stash Hygiene)
*"Localimde çalışıyordu ama sunucuda patladı" durumlarını önler:*

- [ ] **Doğru Dal Doğrulandı:** `git branch --show-current` çıktısı hedeflenen geliştirme dalını işaret ediyor.
- [ ] **Çalışma Ağacı Tertemiz:** `git status` çıktısında unutulmuş, sahnelenmemiş veya `untracked` bırakılmış kritik dosya yok. Kodun çalışması için gereken tüm yeni dosyalar commit edildi.
- [ ] **Uzak Dal İle Senkron:** `git fetch origin` yapıldı; yerel dal uzak dalın gerisinde kalmadı. Olası çakışmalar (merge conflicts) uzak depoya basılmadan önce yerelde çözüldü.

---

## 3. Derin Doğrulama ve CI Eşitliği (Full Verification Suite)
*Pre-push aşaması yalnızca hızlı linting ile yetinmez; derin test zincirini yerelde koşturur:*

- [ ] **Tüm Birim ve Entegrasyon Testleri Başarılı:** `make test` sıfır hatayla geçti (tüm testler passed, exit code 0).
- [ ] **Statik Tip ve Kod Denetimi Temiz:** `go vet ./...` ve Python tarafında `mypy` sıfır hata verdi.
- [ ] **Üretim Derlemesi (Build) Başarılı:** `make build` (Go binary derlemesi) hatasız tamamlandı.
- [ ] **Canlı Çalışma Doğrulaması (Smoke Test):** Köprü ve tetikleyici modülleri ayakta mı? WhatsApp/SSE olay hattı ve yerel komutlar beklendiği gibi yanıt veriyor mu?

---

## 4. Git Geçmişi ve Commit Hijyeni (History Sanity)
*Uzak depoya basılacak commit geçmişinin temiz ve okunabilir olması gerekir:*

- [ ] **Gönderilecek Commitler İncelendi:** `git log origin/<hedef-dal>..HEAD --stat` veya `git log --oneline -n 5` ile sadece gönderilmek istenen commitlerin listede olduğu doğrulandı.
- [ ] **Geçici ve Tembel Commitler Ayıklandı:** Geliştirme esnasında atılmış "wip", "deneme", "fix typo", "temp" gibi dağınık commitler varsa `git rebase -i` ile squash edilip temizlendi.
- [ ] **Tarihsel Sızıntı Kontrolü:** Önceki bir commit'te yanlışlıkla eklenip bir sonraki commit'te silinmiş gizli bir anahtar (secret/token) veya PII yok. (Varsa git geçmişi temizlenmeden push yapılamaz).

---

## 5. Yıkıcı Eylem Koruması (Force Push Safeguards)
*Ekip çalışmasını ve üretim dallarını korur:*

- [ ] **`--force` Kesinlikle Yasak:** Açık ve acil bir felaket kurtarma talimatı olmadıkça `git push --force` veya `-f` kullanılmıyor.
- [ ] **İstisnai Rebase Durumu:** Yalnızca tek kişinin çalıştığı özel özellik dalında rebase sonrası geçmiş güncelleniyorsa, körlemesine `--force` yerine mutlaka `--force-with-lease` kullanılıyor.

---

## 6. MFMX Proje Standartları Uyumu (Repo Conventions Compliance)

Push etmeden önce depodaki her markdown ve yapılandırma dosyasının ilgili konvansiyona uygunluğu doğrulanmalıdır:

| Dosya | Tanım, Ne İşe Yarar & Nasıl Olmalı? (Kurallar) | Asıl Kural Bağlantısı |
| :--- | :--- | :--- |
| **`TODO.md`** | **Aktif Sprint & Açık İşler:** Yalnızca henüz tamamlanmamış açık görevleri listeler. Tamamlanan maddeler bu dosyada şişirilip tutulmaz; anında `TODO.DONE.md`'ye aktarılır. Başında Context Anchor güncel tutulur. | [TODO.md Kuralları](https://github.com/Mefamex/mfmx-project-conventions/blob/main/TODO.md) |
| **`TODO.DONE.md`** | **Doğrulanmış İş Arşivi:** Biten işlerin kayıt defteridir. **Ters kronolojik** (en yeni iş en tepede) yazılır. Her kayıtta neyin değiştiği, kök neden ve **somut doğrulama kanıtı** (test çıktısı, exit code 0, komut sonucu) bulunmak **ZORUNDADIR**. Kanıtsız iş arşive taşınamaz. | [TODO.DONE.md Kuralları](https://github.com/Mefamex/mfmx-project-conventions/blob/main/TODO.DONE.md) |
| **`NOTES.md`** | **Dahili Teknik Bellek:** Kodun "neden" öyle yazıldığını belgeler. Denenip vazgeçilen yaklaşımlar (*Discarded Approaches*), bizzat ölçülmüş sistem limitleri (*Verified Constraints*) ve platform tuzakları (*Quirks & Gotchas*) buraya diske sabitlenir. Tahmin değil ölçüm yazılır; geçmiş tecrübeler silinmez, revizyon notu düşülür. | [NOTES.md Kuralları](https://github.com/Mefamex/mfmx-project-conventions/blob/main/NOTES.md) |
| **`ASK.md`** | **Açık Sorular & Karar Kuyruğu:** *"Ölçülmemiş tahminle koda girmek yok."* Mimari ikilemler tahminle geçiştirilmez, seçenekler ve ajanın önerisi buraya yazılıp kullanıcı beklenir. Kullanıcı yokken inisiyatif alındıysa *Kendi Kararım* tablosuna gerekçesiyle işlenir. Çözülen sorular arşivlenir, dosyada birikmez. | [ASK.md Kuralları](https://github.com/Mefamex/mfmx-project-conventions/blob/main/ASK.md) |
| **`README.md`** | **Proje Vitrini & Mimari Rehber:** Standart meta tablosu ve 1-3 cümlelik vurucu özetle başlar. 9 zorunlu başlık hiyerarşisine uyar. Belge haritası tablosu (`Repository Doc Map`) eksiksiz olmalı, ASCII mimari akış şeması içermeli, güncel kurulum ve kullanım komutlarını yansıtmalıdır. | [README Standartları](https://github.com/Mefamex/mfmx-project-conventions/blob/main/README.standarts.md) |
| **`GIT_STANDARDS.md`** | **Commit & Versiyonlama Kuralları:** Conventional Commits v1.0.0 (`feat`, `fix`, `refactor`, `perf`, `docs`, `test`, `chore`, `build`) formatı zorunludur. Küçük harf, emir kipi, noktasız ve max 72 karakter. Atomik commit esastır. İznisiz push ve force push yasaktır. | [GIT_STANDARDS.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/GIT_STANDARDS.md) |
| **`Makefile`** | **Yaşam Döngüsü & Görev Yöneticisi:** Self-documenting olmalıdır (`make help`). 6 zorunlu hedef (`install`, `reinstall`, `update`, `run`, `clean`, `remove`) ile denetim hedefleri (`lint`, `test`, `verify`, `build`) eksiksiz bulunmalıdır. Push öncesi `make verify` hatasız geçmelidir. | [Makefile Şablonu](https://github.com/Mefamex/mfmx-project-conventions/blob/main/Makefile) |
| **`config.yaml`** | **Merkezi Yapılandırma:** Sistemin çalışma parametreleri kaskad mimaride tanımlanır. Parola, token ve secret'lar ASLA düz metin yazılmaz (`.env` ve ortam değişkeni üzerinden okunur). Tüm parametrelerin inline açıklaması bulunur. | [config.yaml Şablonu](https://github.com/Mefamex/mfmx-project-conventions/blob/main/config.yaml) |

### 🏛️ NerveDispatch Özel Dokümanları
| Belge | Tanım ve Sorumluluk |
| :--- | :--- |
| **`YENI-MIMARI-VE-BIRLESME-PLANI.md`** | `whatsapp-bridge-mcp` ve `whatsapp-tetikleyici` depolarının tek bir monorepoda birleştirilmesi, çok kanallı (WhatsApp, Telegram, Discord, Instagram) ağ geçidi mimarisi ve faz geçişlerinin master planı. |
| **`NERVE-DISPATCH-THE-STORY.md`** | 2024 CHIMERA projesinden günümüze yaşanan mimari kırılmalar, tecrübeler, model/cihaz evrimi ve NerveDispatch vizyonunun resmi tarihçesi. |

---

## 7. Push Öncesi Standart Komut Seti

Push komutunu vermeden önce terminalinizde sırayla şu adımları izleyin:

```bash
# 1. Çalışma ağacının temiz olduğunu doğrulayın
git status

# 2. Gönderilecek commit listesini ve değişen dosyaları görün
git log origin/$(git branch --show-current)..HEAD --stat

# 3. Derin doğrulama ve test paketini yerelde çalıştırın
make verify    # veya make test && make build

# 4. KULLANICI ONAYI ALINDIKTAN SONRA güvenli push çalıştırın
git push origin feature/yeni-ozellik
```

---

## 8. Örnek Git Pre-Push Kancası (.git/hooks/pre-push)

Bu emniyet ağını yerel Git ortamınıza bağlamak için `.git/hooks/pre-push` dosyası oluşturup çalıştırılabilir yapabilirsiniz (`chmod +x .git/hooks/pre-push`):

```bash
#!/usr/bin/env bash
set -e

current_branch=$(git branch --show-current)

# 1. Korumalı dallara doğrudan push engeli
if [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ] || [ "$current_branch" = "production" ]; then
    echo "HATA: Korumalı '$current_branch' dalına doğrudan push yapamazsınız!"
    echo "Lütfen bir özellik dalı açıp Pull Request (PR) kullanın."
    exit 1
fi

echo "==> [Pre-Push] Derin doğrulama ve testler çalıştırılıyor..."

# 2. Test ve derleme kontrolü
if [ -f Makefile ] && grep -q "^verify:" Makefile; then
    make verify
elif [ -f Makefile ] && grep -q "^test:" Makefile; then
    make test
fi

echo "==> [Pre-Push] Tüm kontroller başarıyla tamamlandı. Push izni verildi."
```

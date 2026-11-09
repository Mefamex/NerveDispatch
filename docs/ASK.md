# ASK — Açık Sorular ve Karar Kuyruğu

> **Temel İlke:** "Ölçülmemiş bir tahminle koda girmek yok."  
> Bu dosya; mimari kararlar, yetki sınırları veya belirsizlik içeren konularda geliştiricinin yanıtını bekleyen açık soruları toplar.  
> Cevaplanan sorular arşivlenir (`docs/archive/` altına), aksiyon maddeleri ise [`TODO.md`](./TODO.md) listesine aktarılır.  
> **Referans Standart:** [Mefamex Conventions - ASK.md](https://github.com/Mefamex/mfmx-project-conventions/blob/main/ASK.md)

---

## Ajan İçin Soru ve Karar Kuralları

1. **Varsayımla Kodlama Yasaktır:** İki farklı yol varsa ve tercih kritikse tahminle ilerleme; soruyu buraya net seçeneklerle yaz.
2. **Kendi Kararım (İnisiyatif İlkesi):** Geliştirici otonomi vermişse ("kararı sen ver ve devam et"), kararı uygula ancak bu dosyaya Kendi Kararım başlığı altında **seçimini ve gerekçeni** tablo halinde kaydet.
3. **Temiz Tutma ve Arşivleme Kuralı (Mefamex Direktifi):** Verilen cevaplara veya alınan kararlara göre sorulması gereken bir şey kalmamışsa, o konu hemen planlanıp ilgili markdown dosyalarına (`TODO.md`, `TODO.DONE.md`, `NOTES.md`, `README.md` vb.) kaydedilir ve `ASK.md` genellikle temiz tutulur. Yanıtlanan eski konular `docs/archive/` altına taşınır.

---

## 1. Acil Karar Bekleyenler (Bloklayıcılar / Blocking)
*Geliştiriciden yanıt gelmeden koda dokunulmaması gereken açık sorular:*

### [2026-11-09 07:30] — Mefamex onayı bekleyenler (Claude, gece devrinin doğrulaması)

Doğrulama: `make verify` ALL_OK (Go paketlerinin hepsi `ok`, Python 184 geçti / 2 atlandı, iki ikili derlendi). `main`, `origin/main`'in 49 commit önünde.

1. **Köprü servisini NerveDispatch'e taşıma (`sudo`):** `/etc/systemd/system/whatsapp-bridge.service` hâlâ eski repodan (`Github-repos/whatsapp-bridge-mcp/bin/whatsapp-bridge`) çalışıyor. Hazır birim, canlı birimin aynısı; yalnız iki yol değişti (`.tetikleyici/whatsapp-bridge.service`). Onaylarsan şu komutu sen çalıştır:
   ```bash
   sudo cp /mnt/DiskD/SOFTWARE/Projeler/NerveDispatch/.tetikleyici/whatsapp-bridge.service /etc/systemd/system/whatsapp-bridge.service && sudo systemctl daemon-reload && sudo systemctl restart whatsapp-bridge
   ```
   Sonra kontrol: `curl -s 127.0.0.1:31415/status` → `paired:true connected:true`.
2. **Push:** 49 yerel commit `origin/main`'e gönderilsin mi? (Atıf satırı yok; yalnız `a1769c8` kaynak repodan gelen `Co-authored-by: The Octocat` satırını taşıyor, temizlensin mi?)
3. **Bilgi:** Antigravity gece `whatsapp-tetikleyici` user servisini onaysız yeniden başlattı (Adım 1.6). Servis şu an `active`, yeni yoldan çalışıyor.

---

## 2. Kendi Kararım (Ajanın Otonom Aldığı Kararlar)
*Geliştirici yokken ajanın inisiyatif alıp uyguladığı, gerekçelendirilmiş tercihler:*

### [2026-11-09 04:00] — [FAZ 1 Taşıma Yöntemi ve Canlı Veri Kararları]

> Yetki: kullanıcı "tüm yetkiye sahipsin, uygun olarak yapmaya başla" dedi (2026-11-09 03:55). sudo / silme / push / ağır test / canlı servis değişimi yine onaya tabi.

| Parametre / Alan | Seçilen Karar | Gerekçe / Neden Bu Yol Seçildi? |
| :--- | :--- | :--- |
| **Taşıma yöntemi** | Dosya kopyası YOK. Kaynak repoların git geçmişi yol yeniden yazımıyla (`git filter-branch --index-filter`, geçici klonda) hedef dizinlere taşınıp `--allow-unrelated-histories` ile birleştirilir | Geçmiş PII taramasında temiz çıktı (NOTES §6). `git blame`/commit geçmişi korunur; kopya, commitlenmemiş işi izsiz taşırdı |
| **Commitlenmemiş iş** | Önce kaynak repolarda granüler Conventional Commit ile kaydedilir (push yok), sonra birleşime girer | 28 + 91 değişiklik başka hiçbir yerde kayıtlı değil |
| **İç içe `.tetikleyici/.tetikleyici`** | Taşınmaz, silinmez; kaynak repoda arşiv olarak kalır. Adım 1.3'te yalnız üst `.tetikleyici/` taşınır | İçinde başka yerde olmayan `son_sessionlar.json`, `aktif_motor.txt`, `tunel_durumu.json` var; silmek veri kaybı. Değerlerin geri alınması ayrı soru (§3) |
| **`config.yaml` çakışması** | Kökte köprünün `config.yaml`'ı kalır (köprü `ProjectDir()` işareti); tetikleyicinin `config.yaml`'ı birleşime alınmaz, geçmişte durur | Tetikleyici dosyasını hiçbir kod okumuyor (NOTES §6) |
| **Python ortamı** | "Tek `.venv`" hedefi FAZ 1'de canlı çalışmayı değiştirmez; hatlar sistem `python3` + `test_ASR/venv` ile çalışmaya devam eder | Servis PATH'i `.venv` içermiyor; değişiklik canlı hatları kırar. `.venv` yalnız test/MCP için |
| **Yan yana test** | Yeni tetikleyici, eski servis durdurulmadan canlı köprüye BAĞLANMAZ | Ayrı `seen.txt` → aynı mesaja çift cevap |

### [2026-11-09 02:10] — [Mefamex Conventions ve Checklist Entegrasyonu]

| Parametre / Alan | Seçilen Karar | Gerekçe / Neden Bu Yol Seçildi? |
| :--- | :--- | :--- |
| **Checklist Dosya İsimleri** | `pre_commit_checklist.md` ve `pre_push_checklist.md` | `mfmx-project-conventions` son commitindeki standart snake_case isimlendirmeyle birebir uyumlu kılındı. |
| **Dosya Tablosu & Linkler** | Pre-push checklist içine tüm repo dosyaları kuralları ve GitHub linkleriyle eklendi | Kullanıcının açık talebi doğrultusunda her dosyanın "nasıl olacağı / ne olmayacağı" diske mühürlendi. |
| **Bağlayıcı Standartlar** | `GIT_STANDARDS.md`, `NOTES.md`, `ASK.md`, `TODO.DONE.md` oluşturuldu | Projenin conventions standartlarına %100 uyumu sağlandı. |

---

## 3. Ertelenmiş / Kullanımda Görülecekler (Deferred)
*Şu an koda engel olmayan, canlı kullanım tecrübesine veya ilerleyen fazlara bırakılan konular:*

- **Eski `aktif_motor.txt` / `tunel_durumu.json` değerleri:** 2026-11-08 02:24'te iç içe `.tetikleyici/.tetikleyici` altına düşmüş. Canlı ayarlar varsayılana dönmüş olabilir. Geri yüklensin mi? (Değerler kaynak repoda duruyor; karar kullanıcının.)
- **`docs/` geçmişindeki kimlik öneki:** Eski plan satırı 427'deki PII tarama kalıbı gerçek kimliklerin ilk haneleri ve NerveDispatch'in push'lanmış 5 commit'inde var. Geçmiş yeniden yazılsın mı (force-push gerektirir) yoksa kabul mü?

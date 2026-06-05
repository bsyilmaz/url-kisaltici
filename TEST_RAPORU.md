# Test Raporu

Bu dosya projedeki unit testlerin ozetini iceriyor. Asagidaki sonuclar
`go test ./... -coverprofile=coverage.out -covermode=atomic` komutu ile uretildi.

## Ozet

| | |
|---|---|
| Toplam test | 23 |
| Gecen | 23 |
| Kalan | 0 |
| Sonuc | BASARILI |

Butun testler dis bagimlilik olmadan calisir. Redis ile konusan testlerde
gercek Redis yerine bellek ici `miniredis` kullanilir, bu yuzden testleri
calistirmak icin makinede ayrica bir sey kurmak gerekmez.

## Kapsam (coverage)

| Paket | Kapsam |
|---|---|
| `internal/service` | %92.6 |
| `pkg/shortcode` | %90.0 |
| `internal/handler` | %87.1 |
| `internal/repository` | %78.6 |
| `internal/model` | - (sadece veri yapisi, test edilecek mantik yok) |
| `cmd/server` | %0.0 (giris noktasi / wiring, unit teste tabi tutulmadi) |
| **Toplam** | **%65.5** |

Not: toplam oran `cmd/server` (main) yuzunden dusuk gorunuyor. Orasi sadece
yapilandirma okuyup parcalari birbirine baglayan giris noktasi; is mantigi
icermedigi icin bilerek unit test yazilmadi. Asil is mantiginin oldugu
katmanlarda kapsam %78 ile %93 arasinda.

## Calisan testler

### `pkg/shortcode` (5 test)

- TestGenerateUzunluk - istenen uzunlukta kod uretiliyor mu
- TestGenerateVarsayilanUzunluk - 0 verilince varsayilan uzunluk kullaniliyor mu
- TestGenerateNegatifUzunluk - negatif deger varsayilana dusuyor mu
- TestGenerateSadeceGecerliKarakter - sadece alfabedeki karakterler cikiyor mu
- TestGenerateBenzersizlik - art arda uretilen kodlar cakisiyor mu

### `internal/repository` (4 test)

- TestSaveVeFind - kaydedilen kayit geri okunabiliyor mu
- TestFindBulunamadi - olmayan kod icin ErrNotFound donuyor mu
- TestExists - kod var/yok kontrolu dogru calisiyor mu
- TestSaveTTLDolunca - sure dolunca kayit kendiliginden siliniyor mu

### `internal/service` (8 test)

- TestShortenBasarili - gecerli adres kisaltilip kaydediliyor mu
- TestShortenGecersizURL - bos / http disi adresler reddediliyor mu
- TestShortenBoslukTemizler - bastaki sondaki bosluklar temizleniyor mu
- TestShortenCakismadaVazgecer - surekli cakisma olursa hata donuyor mu
- TestShortenRepositoryHatasi - veri katmani hatasi yukari tasiniyor mu
- TestResolveBasarili - kod karsiligi asil adres donuyor mu
- TestResolveBosKod - bos kod icin ErrNotFound donuyor mu
- TestResolveBulunamadi - olmayan kod icin ErrNotFound donuyor mu

### `internal/handler` (6 test)

- TestShortenHandlerBasarili - 201 ve dogru json donuyor mu
- TestShortenHandlerBozukGovde - bozuk json icin 400 donuyor mu
- TestShortenHandlerGecersizURL - gecersiz url icin 400 donuyor mu
- TestRedirectBasarili - 302 ve dogru Location basligi donuyor mu
- TestRedirectBulunamadi - olmayan kod icin 404 donuyor mu
- TestHealthHandler - saglik ucu 200 donuyor mu

## Nasil tekrar uretilir

```
make test           # testleri ayrintili calistirir
make test-coverage  # kapsam ozetini ekrana basar
make coverage-html  # coverage.html olarak gorsel rapor uretir
```

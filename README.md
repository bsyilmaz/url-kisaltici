# URL Kisaltici Servisi

Golang ve Redis ile yazilmis kucuk bir mikroservis. Uzun linkleri alip kisa bir koda ceviriyor, sonra o kod ile tekrar asil adrese yonlendiriyor. Amac cok buyuk bir sey yapmak degildi aslinda, daha cok temiz mimari nasil kurulur, katmanlar birbirinden nasil ayrilir, testler nasil duzgun yazilir onu gostermek istedim.

## Ne ise yariyor

Diyelim elinizde upuzun bir link var. Servise gonderiyorsunuz, size `aB3xK9z` gibi kisa bir kod donuyor. Sonra `http://localhost:8080/aB3xK9z` adresine girince otomatik olarak asil adrese yonlendiriliyorsunuz. Mantik bu kadar basit. Kayitlar Redis'te tutuluyor ve varsayilan olarak 24 saat sonra kendiliginden siliniyor (TTL sayesinde).

## Mimari

Projeyi katmanli (layered / clean architecture) sekilde ayirdim. Her katmanin tek bir isi var ve birbirine arayuz (interface) uzerinden bagli. Bu sayede test yazarken Redis'e veya gercek bir http sunucusuna ihtiyac kalmiyor, sahte (mock) nesnelerle her seyi test edebiliyoruz.

```
url-kisaltici/
|-- cmd/
|   `-- server/         # uygulamanin giris noktasi (main)
|-- internal/
|   |-- handler/        # http katmani, gelen istekleri karsilar
|   |-- service/        # is kurallari (dogrulama, kod uretme vs)
|   |-- repository/     # veri katmani, Redis ile konusan tek yer
|   `-- model/          # veri yapilari
`-- pkg/
    `-- shortcode/      # rastgele kisa kod ureten kucuk paket
```

Akis soyle isliyor: `handler -> service -> repository -> Redis`. Handler hicbir zaman dogrudan Redis'i bilmiyor, sadece service'i cagiriyor. Service de repository arayuzunu kullaniyor, arkada Redis mi var baska bir sey mi onu hic umursamiyor. Yarin obur gun Redis yerine mesela Postgres koymak istesek sadece yeni bir repository yazmamiz yetiyor, geri kalan kod hic degismiyor. Olay zaten bu ayrimi yapabilmek.

## Calistirmak

### Docker ile (en kolay yol)

Redis ve uygulamayi tek komutla ayaga kaldirir:

```
make up
```

ya da dogrudan:

```
docker compose up --build
```

### Elle (lokalde Redis varsa)

Once bir Redis lazim. Yoksa sadece redis'i docker ile kaldirabilirsiniz:

```
docker run -p 6379:6379 redis:7-alpine
```

Sonra uygulamayi calistirin:

```
make run
```

veya:

```
go run ./cmd/server
```

Servis varsayilan olarak `:8080` portunu dinlemeye baslar.

## API

### Link kisaltma

```
POST /shorten
Content-Type: application/json

{
  "url": "https://cok-uzun-bir-adres.com/falan/filan?x=1"
}
```

Ornek istek:

```
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'
```

Donen cevap:

```
{
  "code": "aB3xK9z",
  "short_url": "http://localhost:8080/aB3xK9z",
  "original_url": "https://github.com"
}
```

Gonderilen adres bos olursa ya da `http` / `https` ile baslamiyorsa `400` doner. Yani `ftp://...` veya duz yazi kabul edilmiyor, bilerek boyle yaptim.

### Yonlendirme

```
GET /{code}
```

Tarayicidan `http://localhost:8080/aB3xK9z` adresine girince dogrudan asil adrese (302) gidersiniz. Kod yoksa `404` doner.

### Saglik kontrolu

```
GET /health
```

`{"status":"ok"}` doner. Liveness/health kontrolu icin kullanilabilir (mesela docker healthcheck veya bir load balancer arkasinda).

## Ayarlar

Her sey ortam degiskeni (environment variable) ile ayarlanabiliyor, kodu degistirmeye gerek yok. Ornek icin `.env.example` dosyasina bakabilirsiniz.

| Degisken | Varsayilan | Aciklama |
|---|---|---|
| `HTTP_ADDR` | `:8080` | Servisin dinleyecegi adres |
| `REDIS_ADDR` | `localhost:6379` | Redis baglanti adresi |
| `REDIS_PASSWORD` | (bos) | Redis sifresi, varsa |
| `BASE_URL` | `http://localhost:8080` | Kisa linklerin basina eklenen adres |

## Testler

Testleri calistirmak icin:

```
make test
```

Kapsam (coverage) raporu icin:

```
make test-coverage
```

Bu komut hangi dosyanin ne kadar test edildigini gosteren bir tablo basiyor. HTML olarak gormek isterseniz:

```
make coverage-html
```

ardindan olusan `coverage.html` dosyasini tarayicida acabilirsiniz. Test sonuclarinin bir ozetini de `TEST_RAPORU.md` dosyasina koydum.

Test yazarken disa bagimliliklari tamamen ayirdim. Repository testleri gercek Redis yerine `miniredis` denen bellek ici sahte bir Redis kullaniyor, yani testleri calistirmak icin makinende Redis kurulu olmasina gerek yok. Service ve handler testleri de sahte (fake) nesnelerle calisiyor, o yuzden cok hizlilar.

## Kullanilan teknolojiler

- Go 1.24 (web framework koymadim, standart kutuphanenin http routerini kullandim)
- Redis (`go-redis` kutuphanesi)
- `miniredis` (sadece testlerde)
- Docker / docker compose

## Neden boyle yaptim

Aslinda URL kisaltma cok basit bir is, isterseniz tek dosyada 50 satirda da yazabilirsiniz. Ama ben burada daha cok "kod nasil duzgun bolunur" onu gostermek istedim. O yuzden:

- Her katman tek bir seyden sorumlu
- Bagimliliklar arayuz uzerinden veriliyor (dependency injection)
- Hatalar `%w` ile sarmalanarak yukari tasiniyor, nerede patladigi kayboluyor degil
- Testler hizli calisiyor cunku disariya bagimli degil
- Kapanirken graceful shutdown var, yarim kalan istek aniden kesilmiyor

Kisacasi kucuk ama duzgun bir iskelet. Uzerine yeni ozellik eklemek kolay olsun diye ugrastim.

## Lisans

MIT. Detay icin `LICENSE` dosyasina bakabilirsiniz.

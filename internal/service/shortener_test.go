package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bsyilmaz/url-kisaltici/internal/repository"
)

// sahteRepo testlerde kullanilan bellek ici repository. Gercek Redis yerine
// gecer ve hata durumlarini taklit edebilmemizi saglar.
type sahteRepo struct {
	veriler   map[string]string
	existsErr error
	saveErr   error
	hepVar    bool // true ise her kod zaten varmis gibi davranir (cakisma testi icin)
}

func yeniSahteRepo() *sahteRepo {
	return &sahteRepo{veriler: make(map[string]string)}
}

func (s *sahteRepo) Save(_ context.Context, code, originalURL string) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.veriler[code] = originalURL
	return nil
}

func (s *sahteRepo) Find(_ context.Context, code string) (string, error) {
	v, ok := s.veriler[code]
	if !ok {
		return "", repository.ErrNotFound
	}
	return v, nil
}

func (s *sahteRepo) Exists(_ context.Context, code string) (bool, error) {
	if s.existsErr != nil {
		return false, s.existsErr
	}
	if s.hepVar {
		return true, nil
	}
	_, ok := s.veriler[code]
	return ok, nil
}

func TestShortenBasarili(t *testing.T) {
	repo := yeniSahteRepo()
	svc := NewShortener(repo)

	res, err := svc.Shorten(context.Background(), "https://ornek.com/uzun/adres")
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if res.Code == "" {
		t.Error("bir kod uretilmeliydi")
	}
	if len(res.Code) != codeLength {
		t.Errorf("kod uzunlugu %d olmaliydi, %d geldi", codeLength, len(res.Code))
	}
	if res.OriginalURL != "https://ornek.com/uzun/adres" {
		t.Errorf("orijinal adres yanlis: %s", res.OriginalURL)
	}
	if repo.veriler[res.Code] != res.OriginalURL {
		t.Error("kayit repository'e yazilmadi")
	}
}

func TestShortenGecersizURL(t *testing.T) {
	repo := yeniSahteRepo()
	svc := NewShortener(repo)

	// Bos, sema icermeyen veya http/https disindaki adresler reddedilmeli.
	testler := []string{"", "  ", "duz yazi", "ftp://dosya.com", "://eksik", "javascript:alert(1)"}
	for _, girdi := range testler {
		_, err := svc.Shorten(context.Background(), girdi)
		if !errors.Is(err, ErrGecersizURL) {
			t.Errorf("%q icin ErrGecersizURL bekleniyordu, gelen %v", girdi, err)
		}
	}
}

func TestShortenBoslukTemizler(t *testing.T) {
	repo := yeniSahteRepo()
	svc := NewShortener(repo)

	res, err := svc.Shorten(context.Background(), "   https://ornek.com   ")
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if res.OriginalURL != "https://ornek.com" {
		t.Errorf("bastaki ve sondaki bosluklar temizlenmeliydi, gelen %q", res.OriginalURL)
	}
}

func TestShortenCakismadaVazgecer(t *testing.T) {
	repo := yeniSahteRepo()
	repo.hepVar = true // her uretilen kod zaten varmis gibi davranacak
	svc := NewShortener(repo)

	_, err := svc.Shorten(context.Background(), "https://ornek.com")
	if !errors.Is(err, ErrUretilemedi) {
		t.Errorf("surekli cakisma olunca ErrUretilemedi bekleniyordu, gelen %v", err)
	}
}

func TestShortenRepositoryHatasi(t *testing.T) {
	repo := yeniSahteRepo()
	repo.saveErr = errors.New("redis coktu")
	svc := NewShortener(repo)

	_, err := svc.Shorten(context.Background(), "https://ornek.com")
	if err == nil {
		t.Error("repository hata verince Shorten da hata dondurmeliydi")
	}
}

func TestResolveBasarili(t *testing.T) {
	repo := yeniSahteRepo()
	repo.veriler["abc"] = "https://ornek.com"
	svc := NewShortener(repo)

	val, err := svc.Resolve(context.Background(), "abc")
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if val != "https://ornek.com" {
		t.Errorf("yanlis adres dondu: %s", val)
	}
}

func TestResolveBosKod(t *testing.T) {
	repo := yeniSahteRepo()
	svc := NewShortener(repo)

	_, err := svc.Resolve(context.Background(), "   ")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("bos kod icin ErrNotFound bekleniyordu, gelen %v", err)
	}
}

func TestResolveBulunamadi(t *testing.T) {
	repo := yeniSahteRepo()
	svc := NewShortener(repo)

	_, err := svc.Resolve(context.Background(), "olmayan")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("ErrNotFound bekleniyordu, gelen %v", err)
	}
}

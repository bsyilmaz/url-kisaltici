// Package service kisaltma is kurallarini barindirir. Bu katman http veya
// Redis gibi detaylari bilmez; sadece repository arayuzunu kullanir.
package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bsyilmaz/url-kisaltici/internal/model"
	"github.com/bsyilmaz/url-kisaltici/internal/repository"
	"github.com/bsyilmaz/url-kisaltici/pkg/shortcode"
)

const (
	// codeLength uretilen kisa kodlarin uzunlugu.
	codeLength = 7
	// maxRetry kod cakismasi durumunda en fazla kac kez yeniden denenecegi.
	maxRetry = 5
)

var (
	// ErrGecersizURL gonderilen adres bos ya da http/https degilse doner.
	ErrGecersizURL = errors.New("gecersiz url")
	// ErrUretilemedi tekrar tekrar denenmesine ragmen bos bir kod bulunamazsa doner.
	ErrUretilemedi = errors.New("kisa kod uretilemedi")
)

// Shortener kisaltma ve cozme islemlerini yurutur.
type Shortener struct {
	repo repository.Repository
}

// NewShortener verilen repository ile yeni bir Shortener olusturur.
func NewShortener(repo repository.Repository) *Shortener {
	return &Shortener{repo: repo}
}

// Shorten verilen adresi dogrular, benzersiz bir kod uretir ve kaydeder.
func (s *Shortener) Shorten(ctx context.Context, originalURL string) (model.URL, error) {
	originalURL = strings.TrimSpace(originalURL)
	if !gecerliMi(originalURL) {
		return model.URL{}, ErrGecersizURL
	}

	for i := 0; i < maxRetry; i++ {
		code, err := shortcode.Generate(codeLength)
		if err != nil {
			return model.URL{}, fmt.Errorf("kod uretilemedi: %w", err)
		}

		exists, err := s.repo.Exists(ctx, code)
		if err != nil {
			return model.URL{}, err
		}
		if exists {
			// Cok dusuk ihtimal ama ayni kod zaten varsa yenisini dene.
			continue
		}

		if err := s.repo.Save(ctx, code, originalURL); err != nil {
			return model.URL{}, err
		}
		return model.URL{Code: code, OriginalURL: originalURL}, nil
	}

	return model.URL{}, ErrUretilemedi
}

// Resolve kisa koda karsilik gelen asil adresi doner.
func (s *Shortener) Resolve(ctx context.Context, code string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", repository.ErrNotFound
	}
	return s.repo.Find(ctx, code)
}

// gecerliMi adresin bos olmadigini, http/https semasi tasidigini ve bir alan
// adi (host) icerdigini kontrol eder. Boylece "https://" gibi semasi olup
// host'u olmayan adresler reddedilir.
func gecerliMi(raw string) bool {
	if raw == "" {
		return false
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if u.Host == "" {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

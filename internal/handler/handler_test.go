package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bsyilmaz/url-kisaltici/internal/model"
	"github.com/bsyilmaz/url-kisaltici/internal/repository"
	"github.com/bsyilmaz/url-kisaltici/internal/service"
)

// sahteServis Shortener arayuzunu testler icin taklit eder. Hangi cevabin
// donecegini her testte fonksiyon olarak veriyoruz.
type sahteServis struct {
	shortenFn func(ctx context.Context, url string) (model.URL, error)
	resolveFn func(ctx context.Context, code string) (string, error)
}

func (s sahteServis) Shorten(ctx context.Context, url string) (model.URL, error) {
	return s.shortenFn(ctx, url)
}

func (s sahteServis) Resolve(ctx context.Context, code string) (string, error) {
	return s.resolveFn(ctx, code)
}

func TestShortenHandlerBasarili(t *testing.T) {
	svc := sahteServis{
		shortenFn: func(_ context.Context, url string) (model.URL, error) {
			return model.URL{Code: "aB3xK9z", OriginalURL: url}, nil
		},
	}
	h := New(svc, "http://kisa.lt")

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url":"https://ornek.com"}`))
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("201 bekleniyordu, gelen %d", rec.Code)
	}

	var resp shortenResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("cevap cozulemedi: %v", err)
	}
	if resp.ShortURL != "http://kisa.lt/aB3xK9z" {
		t.Errorf("yanlis kisa adres: %s", resp.ShortURL)
	}
	if resp.Code != "aB3xK9z" {
		t.Errorf("yanlis kod: %s", resp.Code)
	}
}

func TestShortenHandlerBozukGovde(t *testing.T) {
	svc := sahteServis{
		shortenFn: func(_ context.Context, _ string) (model.URL, error) {
			return model.URL{}, nil
		},
	}
	h := New(svc, "http://kisa.lt")

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader("{bu bozuk json"))
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("bozuk json icin 400 bekleniyordu, gelen %d", rec.Code)
	}
}

func TestShortenHandlerGecersizURL(t *testing.T) {
	svc := sahteServis{
		shortenFn: func(_ context.Context, _ string) (model.URL, error) {
			return model.URL{}, service.ErrGecersizURL
		},
	}
	h := New(svc, "http://kisa.lt")

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url":"gecersiz"}`))
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("gecersiz url icin 400 bekleniyordu, gelen %d", rec.Code)
	}
}

func TestRedirectBasarili(t *testing.T) {
	svc := sahteServis{
		resolveFn: func(_ context.Context, _ string) (string, error) {
			return "https://ornek.com", nil
		},
	}
	h := New(svc, "http://kisa.lt")

	req := httptest.NewRequest(http.MethodGet, "/aB3xK9z", nil)
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("302 bekleniyordu, gelen %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://ornek.com" {
		t.Errorf("yanlis yonlendirme adresi: %s", loc)
	}
}

func TestRedirectBulunamadi(t *testing.T) {
	svc := sahteServis{
		resolveFn: func(_ context.Context, _ string) (string, error) {
			return "", repository.ErrNotFound
		},
	}
	h := New(svc, "http://kisa.lt")

	req := httptest.NewRequest(http.MethodGet, "/olmayankod", nil)
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("404 bekleniyordu, gelen %d", rec.Code)
	}
}

func TestHealthHandler(t *testing.T) {
	h := New(sahteServis{}, "http://kisa.lt")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("200 bekleniyordu, gelen %d", rec.Code)
	}
}

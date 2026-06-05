// Package handler http katmanidir. Gelen istekleri karsilar, is katmanini
// cagirir ve sonucu json olarak dondurur. Is kurallari burada bulunmaz.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bsyilmaz/url-kisaltici/internal/model"
	"github.com/bsyilmaz/url-kisaltici/internal/repository"
	"github.com/bsyilmaz/url-kisaltici/internal/service"
)

// Shortener handler'in ihtiyac duydugu is katmani davranisini tanimlar.
// Arayuz olarak tutuldugu icin testlerde sahte bir nesne verilebiliyor.
type Shortener interface {
	Shorten(ctx context.Context, originalURL string) (model.URL, error)
	Resolve(ctx context.Context, code string) (string, error)
}

// Handler http isteklerini is katmanina baglar.
type Handler struct {
	svc     Shortener
	baseURL string
}

// New yeni bir Handler olusturur. baseURL kisa linklerin basina eklenir.
func New(svc Shortener, baseURL string) *Handler {
	return &Handler{svc: svc, baseURL: baseURL}
}

// Routes servisin http yollarini kaydeder ve hazir bir handler doner.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /shorten", h.shorten)
	mux.HandleFunc("GET /{code}", h.redirect)
	return mux
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code        string `json:"code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "istek govdesi okunamadi"})
		return
	}

	res, err := h.svc.Shorten(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, service.ErrGecersizURL) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "gecersiz url adresi"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "kisaltma yapilamadi"})
		return
	}

	writeJSON(w, http.StatusCreated, shortenResponse{
		Code:        res.Code,
		ShortURL:    h.baseURL + "/" + res.Code,
		OriginalURL: res.OriginalURL,
	})
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	original, err := h.svc.Resolve(r.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "adres bulunamadi"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "adres cozulemedi"})
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// writeJSON verilen veriyi json olarak yazar ve durum kodunu ayarlar.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

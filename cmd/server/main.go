// Komut url-kisaltici servisini ayaga kaldirir. Burasi uygulamanin giris
// noktasi; yapilandirma okunur, bagimliliklar kurulur ve sunucu baslatilir.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bsyilmaz/url-kisaltici/internal/handler"
	"github.com/bsyilmaz/url-kisaltici/internal/repository"
	"github.com/bsyilmaz/url-kisaltici/internal/service"
)

func main() {
	cfg := yapilandirmaYukle()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.redisAddr,
		Password: cfg.redisPassword,
	})

	// Acilista Redis'e ulasabiliyor muyuz diye bir kontrol edelim.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redise baglanilamadi: %v", err)
	}

	repo := repository.NewRedisRepository(rdb, cfg.ttl)
	svc := service.NewShortener(repo)
	h := handler.New(svc, cfg.baseURL)

	srv := &http.Server{
		Addr:         cfg.httpAddr,
		Handler:      h.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Sunucuyu ayri bir goroutine'de calistir ki kapatma sinyalini dinleyebilelim.
	go func() {
		log.Printf("servis %s adresinde dinliyor", cfg.httpAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("sunucu hatasi: %v", err)
		}
	}()

	// Ctrl+C veya kill geldiginde duzgun bir sekilde kapan (graceful shutdown).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("kapatma sinyali alindi, servis kapatiliyor")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("kapatma sirasinda hata olustu: %v", err)
	}
	_ = rdb.Close()
	log.Println("servis kapandi")
}

// yapilandirma servisin calismasi icin gereken ayarlari tutar.
type yapilandirma struct {
	httpAddr      string
	redisAddr     string
	redisPassword string
	baseURL       string
	ttl           time.Duration
}

// yapilandirmaYukle ayarlari ortam degiskenlerinden okur, yoksa varsayilani kullanir.
func yapilandirmaYukle() yapilandirma {
	return yapilandirma{
		httpAddr:      envOrDefault("HTTP_ADDR", ":8080"),
		redisAddr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		redisPassword: envOrDefault("REDIS_PASSWORD", ""),
		baseURL:       envOrDefault("BASE_URL", "http://localhost:8080"),
		ttl:           24 * time.Hour,
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

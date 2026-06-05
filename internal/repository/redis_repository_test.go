package repository

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// testRepo bellek ici sahte bir Redis (miniredis) ile repository hazirlar.
// Boylece testlerin calismasi icin gercek bir Redis kurulu olmasina gerek yok.
func testRepo(t *testing.T) (*RedisRepository, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis baslatilamadi: %v", err)
	}
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() {
		_ = client.Close()
		mr.Close()
	})
	return NewRedisRepository(client, time.Hour), mr
}

func TestSaveVeFind(t *testing.T) {
	repo, _ := testRepo(t)
	ctx := context.Background()

	if err := repo.Save(ctx, "abc123", "https://ornek.com"); err != nil {
		t.Fatalf("kaydetme hatasi: %v", err)
	}

	val, err := repo.Find(ctx, "abc123")
	if err != nil {
		t.Fatalf("okuma hatasi: %v", err)
	}
	if val != "https://ornek.com" {
		t.Errorf("beklenen https://ornek.com, gelen %s", val)
	}
}

func TestFindBulunamadi(t *testing.T) {
	repo, _ := testRepo(t)

	_, err := repo.Find(context.Background(), "olmayankod")
	if err != ErrNotFound {
		t.Errorf("ErrNotFound bekleniyordu, gelen %v", err)
	}
}

func TestExists(t *testing.T) {
	repo, _ := testRepo(t)
	ctx := context.Background()

	mevcut, err := repo.Exists(ctx, "kod1")
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if mevcut {
		t.Error("kod henuz kaydedilmedi, mevcut gorunmemeli")
	}

	if err := repo.Save(ctx, "kod1", "https://ornek.com"); err != nil {
		t.Fatalf("kaydetme hatasi: %v", err)
	}

	mevcut, err = repo.Exists(ctx, "kod1")
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if !mevcut {
		t.Error("kod kaydedildi, mevcut gorunmeli")
	}
}

func TestSaveTTLDolunca(t *testing.T) {
	repo, mr := testRepo(t)
	ctx := context.Background()

	if err := repo.Save(ctx, "gecici", "https://ornek.com"); err != nil {
		t.Fatalf("kaydetme hatasi: %v", err)
	}

	// miniredis ile zamani ileri sararak TTL'in dolmasini taklit ediyoruz.
	mr.FastForward(2 * time.Hour)

	_, err := repo.Find(ctx, "gecici")
	if err != ErrNotFound {
		t.Errorf("ttl dolunca ErrNotFound bekleniyordu, gelen %v", err)
	}
}

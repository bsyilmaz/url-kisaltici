package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// keyPrefix Redis tarafinda anahtarlarin basina eklenir, boylece
// ayni veritabaninda baska veriler ile karismaz.
const keyPrefix = "url:"

// RedisRepository, Repository arayuzunu Redis kullanarak uygular.
type RedisRepository struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisRepository yeni bir Redis tabanli repository olusturur.
// ttl, kayitlarin ne kadar sure sonra otomatik silinecegini belirler.
func NewRedisRepository(client *redis.Client, ttl time.Duration) *RedisRepository {
	return &RedisRepository{client: client, ttl: ttl}
}

func (r *RedisRepository) Save(ctx context.Context, code, originalURL string) error {
	if err := r.client.Set(ctx, keyPrefix+code, originalURL, r.ttl).Err(); err != nil {
		return fmt.Errorf("redise yazilamadi: %w", err)
	}
	return nil
}

func (r *RedisRepository) Find(ctx context.Context, code string) (string, error) {
	val, err := r.client.Get(ctx, keyPrefix+code).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("redisten okunamadi: %w", err)
	}
	return val, nil
}

func (r *RedisRepository) Exists(ctx context.Context, code string) (bool, error) {
	n, err := r.client.Exists(ctx, keyPrefix+code).Result()
	if err != nil {
		return false, fmt.Errorf("redis kontrol edilemedi: %w", err)
	}
	return n > 0, nil
}

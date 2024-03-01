package cache

import (
	bole "boletia/api/internal"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository struct {
	client         *redis.Client
	expirationTime time.Duration
	nameSpace      string
}

// NewCacheRepository initializes a Postgres-based implementation of bole.CacheRepository.
func NewCacheRepository(client *redis.Client, expirationTime time.Duration, nameSpace string) *CacheRepository {
	return &CacheRepository{
		client:         client,
		expirationTime: expirationTime,
		nameSpace:      nameSpace,
	}
}

// Get data from cache
func (r *CacheRepository) Get(ctx context.Context, hash string) ([]bole.Currency, error) {
	key := r.makeNamespace(hash)
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return []bole.Currency{}, nil
	}
	if err != nil {
		return []bole.Currency{}, err
	}
	var profile []bole.Currency
	if err = json.Unmarshal(val, &profile); err != nil {
		return []bole.Currency{}, err
	}
	return profile, err
}

// Set cache a data on redis
func (r *CacheRepository) Set(ctx context.Context, hash string, data []bole.Currency) error {
	key := r.makeNamespace(hash)
	profile, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := r.client.Set(ctx, key, profile, r.expirationTime).Err(); err != nil {
		return err
	}
	return nil
}

func (r *CacheRepository) makeNamespace(key string) string {
	return fmt.Sprintf("%s-%s", r.nameSpace, key)
}

package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		Client: client,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, bool, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return value, true, nil
}

func (r *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisCache) DelByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	pattern := prefix + "*"

	for {
		keys, nextCursor, err := r.Client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.Client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if nextCursor == 0 {
			break
		}

		cursor = nextCursor
	}

	return nil
}

package cache

import (
	"context"
	"time"

	redisStorage "github.com/gofiber/storage/redis"
)

type Service interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Ping(ctx context.Context) error
	GetAndDel(ctx context.Context, key string) (string, bool, error)
	SAdd(ctx context.Context, key string, members any, ttl time.Duration) error
	SRem(ctx context.Context, key string, members ...any) error
	SMembers(ctx context.Context, key string) ([]string, error)
	Close() error
	GetStorage() *redisStorage.Storage
}

package cache

import (
	"context" 
	"fmt"
	"rest-fiber/config"
	"runtime"
	"strconv"
	"strings"
	"time"

	redisStorage "github.com/gofiber/storage/redis"
	"github.com/redis/go-redis/v9"
)

type serviceImpl struct {
	client  *redis.Client
	env     config.Env
	storage *redisStorage.Storage
}
func NewService(env config.Env) (Service, error) {
	host := strings.TrimSpace(env.RedisHost)
	portStr := strings.TrimSpace(env.RedisPort)
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return nil, fmt.Errorf("invalid REDIS_PORT=%q", portStr)
	}
	password := env.RedisPassword
	address := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed (%s): %w", address, err)
	}

	storage := redisStorage.New(redisStorage.Config{
		Host:     host,
		Port:     port,
		Password: password,
		Database: 0,
		PoolSize: 10 * runtime.GOMAXPROCS(0),
	})

	return &serviceImpl{
		client:  client,
		env:     env,
		storage: storage,
	}, nil
}

func (s *serviceImpl) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *serviceImpl) Set(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *serviceImpl) Del(ctx context.Context, keys ...string) error {
	return s.client.Del(ctx, keys...).Err()
}

func (s *serviceImpl) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *serviceImpl) Close() error {
	return s.client.Close()
}

func (s *serviceImpl) GetAndDel(ctx context.Context, key string) (string, bool, error) {
	script := redis.NewScript(`
	local v = redis.call("GET", KEYS[1])
	if not v then return "" end
	redis.call("DEL", KEYS[1])
	return v
	`)
	res, err := script.Run(ctx, s.client, []string{key}).Result()
	if err != nil {
		return "", false, err
	}
	val, _ := res.(string)
	if val == "" {
		return "", false, nil
	}
	return val, true, nil
}

func (s *serviceImpl) SAdd(ctx context.Context, key string, members any, ttl time.Duration) error {
	pipe := s.client.Pipeline()
	pipe.SAdd(ctx, key, members)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *serviceImpl) SRem(ctx context.Context, key string, members ...any) error {
	return s.client.SRem(ctx, key, members...).Err()
}

func (s *serviceImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	return s.client.SMembers(ctx, key).Result()
}

func (s *serviceImpl) Incr(ctx context.Context, key string) (int64, error) {
	return s.client.Incr(ctx, key).Result()
}

func (s *serviceImpl) GetStorage() *redisStorage.Storage {
	return s.storage
}

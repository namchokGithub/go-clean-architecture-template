package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps go-redis for typed access.
type Client struct {
	rdb *redis.Client
}

// Config holds Redis connection settings.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// New creates a Redis client and pings to verify connectivity.
func New(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Client{rdb: rdb}, nil
}

func (c *Client) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, val, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

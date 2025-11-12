package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(addr, password string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &Client{Client: rdb}
}

func (c *Client) SetHold(ctx context.Context, key string, bookingID string, ttl time.Duration) (bool, error) {
	return c.SetNX(ctx, key, bookingID, ttl).Result()
}

func (c *Client) GetHold(ctx context.Context, key string) (string, error) {
	return c.Get(ctx, key).Result()
}

func (c *Client) DeleteHold(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.Client.HSet(ctx, key, values...).Err()
}

func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.Client.HGet(ctx, key, field).Result()
}

func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.Client.Exists(ctx, keys...).Result()
}




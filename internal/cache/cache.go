package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"social-network/internal/app"
	"time"
)

type Cache struct {
	ctx  context.Context
	conn *redis.Client
	dsn  string
}

// New ...
func New(ctx context.Context, dsn string) *Cache {
	return &Cache{
		ctx: ctx,
		dsn: dsn,
	}
}

// Connect ...
func (c *Cache) Connect(_ context.Context) app.Cache {
	opts, err := redis.ParseURL(c.dsn)
	if err != nil {
		panic(err)
	}

	c.conn = redis.NewClient(opts)
	return c
}

// Get ...
func (c *Cache) Get(key string) (result []app.PostDB, exists bool, err error) {
	var val string
	val, err = c.conn.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil, false, errors.New("key does not exist") // Cache miss
	} else if err != nil {
		return nil, false, errors.New("internal server error")
	}

	// Unmarshal JSON into result
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return nil, false, err
	}

	return result, true, nil
}

// Set ...
func (c *Cache) Set(key string, value []app.PostDB, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.conn.Set(c.ctx, key, jsonData, expiration).Err()
}

// Clear ...
func (c *Cache) Clear(key string) error {
	err := c.conn.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

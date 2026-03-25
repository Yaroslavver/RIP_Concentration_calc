package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"lab/internal/app/config"
)

type Client struct {
	client *redis.Client
	cfg    config.RedisConfig
}

func New(cfg config.RedisConfig) (*Client, error) {
	// Формируем адрес
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	
	// Создаем клиент с паролем
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,  // ← обязательно передаём пароль
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.ReadTimeout,
	})
	
	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	
	return &Client{client: rdb, cfg: cfg}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// Blacklist methods
const blacklistPrefix = "jwt_blacklist:"

func (c *Client) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := blacklistPrefix + token
	return c.client.Set(ctx, key, true, ttl).Err()
}

func (c *Client) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := blacklistPrefix + token
	_, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
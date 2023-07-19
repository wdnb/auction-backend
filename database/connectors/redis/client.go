package rdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var client *redis.Client

type Config struct {
	Host     string
	Port     string
	Password string
	DBName   int
}

func GetClient(cfg *Config) *redis.Client {
	if client == nil {
		client = initClient(cfg)
	}

	return client
}

func initClient(cfg *Config) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DBName,
	})
	_, err := connectWithRetry(context.Background(), client)
	if err != nil {
		panic(err)
	}

	return client
}

func connectWithRetry(ctx context.Context, client *redis.Client) (string, error) {
	maxRetries := 3
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		err := client.Ping(ctx).Err()
		if err == nil {
			return "OK", nil
		}

		zap.L().Error("Failed to connect to Redis",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Error(err),
		)

		if i == maxRetries-1 {
			return "", fmt.Errorf("failed to connect after %d retries", maxRetries)
		}

		time.Sleep(retryInterval)
	}
	return "", errors.New("unknown error")
}

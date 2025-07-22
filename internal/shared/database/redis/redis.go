package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"protravel-finance/internal/config"
	"protravel-finance/pkg/logger"
)

func NewRedis(ctx context.Context, cfg config.RedisConfig, log logger.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	log.Info("Connected to Redis Successfully")

	return client, nil
}

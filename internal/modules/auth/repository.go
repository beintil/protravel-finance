package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type authRepository struct{}

func NewAuthRepository() Repository {
	return &authRepository{}
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, redisClient *redis.Client, refreshToken string, jti string) error {
	newCacheKey := fmt.Sprintf("refresh:%s", jti)
	err := redisClient.Set(ctx, newCacheKey, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("SaveRefreshToken/Set: %w", err)
	}
	return nil
}

func (r *authRepository) GetRefreshToken(ctx context.Context, redisClient *redis.Client, jti string) (string, error) {
	cacheKey := fmt.Sprintf("refresh:%s", jti)
	refreshToken, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", RepositoryErrorRefreshTokenNotFound
		}
		return "", fmt.Errorf("GetRefreshToken/Get: %w", err)
	}
	return refreshToken, nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, redisClient *redis.Client, jti string) error {
	cacheKey := fmt.Sprintf("refresh:%s", jti)
	err := redisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		return fmt.Errorf("DeleteRefreshToken/Del: %w", err)
	}
	return nil
}

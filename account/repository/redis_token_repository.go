package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
	"github.com/go-redis/redis/v8"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("保存 refreshToken 失败，userID/TokenID-%s/%s：%v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Del(ctx, key).Err(); err != nil {
		log.Printf("删除 refreshToken 失败，userID/TokenID-%s/%s：%v\n", userID, tokenID, err)
		return apperrors.NewInternal()
	}
	return nil
}

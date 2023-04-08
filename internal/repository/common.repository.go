package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bagastri07/boilerplate-service/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func HSetWithExpiry(ctx context.Context, redisClient *redis.Client, bucketCacheKey string, field string, data any) error {
	cacheData, err := json.Marshal(data)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	err = redisClient.HSet(ctx, bucketCacheKey, field, cacheData).Err()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	err = redisClient.ExpireNX(ctx, bucketCacheKey, config.RedisCacheTTL()).Err()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func SetWithExpiry(ctx context.Context, redisClient *redis.Client, cacheKey string, data any) error {
	cacheData, err := json.Marshal(data)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	err = redisClient.Set(ctx, cacheKey, cacheData, config.RedisCacheTTL()).Err()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func Get(ctx context.Context, redisClient *redis.Client, cacheKey string) ([]byte, error) {
	cachedData, err := redisClient.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		logrus.WithField("cacheKey", cacheKey).Error(err.Error())
		return nil, err
	}
	return cachedData, nil
}

func DeleteByKeys(ctx context.Context, redisClient *redis.Client, cacheKeys []string) error {
	for _, cacheKey := range cacheKeys {
		err := redisClient.Del(ctx, cacheKey).Err()
		if err != nil && !errors.Is(err, redis.Nil) {
			logrus.WithField("cacheKey", cacheKey).Error(err.Error())
			return err
		}
	}
	return nil
}

func HGet(ctx context.Context, redisClient *redis.Client, bucketCacheKey string, field string) ([]byte, error) {
	cachedData, err := redisClient.HGet(ctx, bucketCacheKey, field).Bytes()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return nil, err
	}
	return cachedData, nil
}

package infrastructure

import (
	"github.com/bagastri07/authorization-service/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	RedisClient *redis.Client
)

func InitializeRedisCon() {
	redisURL, err := redis.ParseURL(config.RedisCacheHost())

	if err != nil {
		logrus.Panic(err)
	}

	redisOpts := &redis.Options{
		Network:      redisURL.Network,
		Addr:         redisURL.Addr,
		DB:           redisURL.DB,
		Username:     redisURL.Username,
		Password:     redisURL.Password,
		DialTimeout:  config.RedisDialTimeout(),
		WriteTimeout: config.RedisWriteTimeout(),
		ReadTimeout:  config.RedisReadTimeout(),
	}

	RedisClient = redis.NewClient(redisOpts)
}

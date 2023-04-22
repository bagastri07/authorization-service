package repository

import (
	"context"

	"github.com/bagastri07/authorization-service/internal/helper"
	"github.com/bagastri07/authorization-service/internal/model"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type tokenRepository struct {
	redisClient *goredis.Client
}

func NewTokenRepository(redisClient *goredis.Client) model.TokenRepository {
	return &tokenRepository{
		redisClient: redisClient,
	}
}

func (r *tokenRepository) Create(ctx context.Context, userID uuid.UUID, token *model.Token) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   helper.DumpIncomingContext(ctx),
		"token": helper.Dump(token),
	})

	accessTokenCacheKey := model.AccessTokenCacheKey(userID.String(), token.AccessToken)
	err := SetWithExpiry(ctx, r.redisClient, accessTokenCacheKey, token.AccessToken)
	if err != nil {
		logger.WithField("cacheKey", accessTokenCacheKey).Error(err)
		return err
	}

	refreshTokenCacheKey := model.RefreshTokenCacheKey(userID.String(), token.RefreshToken)
	err = SetWithExpiry(ctx, r.redisClient, refreshTokenCacheKey, token.RefreshToken)
	if err != nil {
		logger.WithField("cacheKey", refreshTokenCacheKey).Error(err)
		return err
	}

	return nil
}

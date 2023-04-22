package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

type TokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, token *Token) error
}

func RefreshTokenCacheKey(userID string, token string) string {
	return fmt.Sprintf("refresh-token:%s:%s", userID, token)
}

func AccessTokenCacheKey(userID string, token string) string {
	return fmt.Sprintf("access-token:%s:%s", userID, token)
}

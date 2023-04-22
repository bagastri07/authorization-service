package helper

import (
	"time"

	"github.com/bagastri07/authorization-service/internal/config"
	"github.com/bagastri07/authorization-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func GenerateJwtToken(ID string) (*model.Token, error) {

	accessTokenExpiration := time.Now().Add(config.AccessTokenDuration())
	claims := &model.Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenStr, err := accessToken.SignedString(config.JwtAccessTokenSecret())
	if err != nil {
		logrus.WithField("ID", ID).Error(err)
		return nil, err
	}

	refreshTokenExpiration := time.Now().Add(config.RefreshTokenDuration())
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(refreshTokenExpiration)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenStr, err := refreshToken.SignedString(config.JwtRefreshTokenSecret())
	if err != nil {
		logrus.WithField("ID", ID).Error(err)
		return nil, err
	}

	return &model.Token{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil

}

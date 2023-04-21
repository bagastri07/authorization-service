package helper

import (
	"time"

	"github.com/bagastri07/authorization-service/internal/config"
	"github.com/bagastri07/authorization-service/internal/constant"
	"github.com/bagastri07/authorization-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func GenerateJwtToken(ID string) (*model.TokenResp, error) {

	accessTokenExpiration := time.Now().Add(constant.JwtAccessTokenExpirationTime)
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

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenStr, err := refreshToken.SignedString(config.JwtRefreshTokenSecret())
	if err != nil {
		logrus.WithField("ID", ID).Error(err)
		return nil, err
	}

	return &model.TokenResp{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil

}

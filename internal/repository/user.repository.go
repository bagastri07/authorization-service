package repository

import (
	"context"

	"github.com/bagastri07/authorization-service/internal/helper"
	"github.com/bagastri07/authorization-service/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	DB          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(DB *gorm.DB, redis *redis.Client) model.UserRepository {
	return &userRepository{
		DB:          DB,
		redisClient: redis,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*uuid.UUID, error) {
	tx := helper.GetTxFromContext(ctx, r.DB)

	err := tx.WithContext(ctx).Save(user).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  helper.DumpIncomingContext(ctx),
			"user": helper.Dump(user),
		}).Error(err)
		return nil, err
	}

	return &user.ID, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)

	err := r.DB.WithContext(ctx).
		Where("email = ?", email).
		First(user).
		Error

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":   helper.DumpIncomingContext(ctx),
			"email": email,
		}).Error(err)
		return nil, err
	}

	return user, nil
}

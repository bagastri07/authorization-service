package repository

import (
	"context"

	"github.com/bagastri07/authorization-service/internal/helper"
	"github.com/bagastri07/authorization-service/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	DB *gorm.DB
}

func NewUserRepoRepository(DB *gorm.DB) model.UserRoleRepository {
	return &userRoleRepository{
		DB: DB,
	}
}

func (r *userRoleRepository) Create(ctx context.Context, userRole *model.UserRole) error {
	tx := helper.GetTxFromContext(ctx, r.DB)

	err := tx.WithContext(ctx).
		Create(userRole).
		Error

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      helper.DumpIncomingContext(ctx),
			"userRole": helper.Dump(userRole),
		}).Error(err)
		return err
	}

	return nil
}

package usecase

import (
	"context"
	"errors"

	"github.com/bagastri07/authorization-service/internal/constant"
	cErr "github.com/bagastri07/authorization-service/internal/constant/customerror"
	"github.com/bagastri07/authorization-service/internal/helper"
	"github.com/bagastri07/authorization-service/internal/infrastructure"
	"github.com/bagastri07/authorization-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepo     model.UserRepository
	userRoleRepo model.UserRoleRepository
	tokenRepo    model.TokenRepository
}

func NewUserUsecase() model.UserUsecase {
	return &userUsecase{}
}

func (uc *userUsecase) Register(ctx context.Context, user *model.User) (*model.Token, error) {
	var err error
	tx := infrastructure.PostgreSQL.Begin()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logrus.Fatal(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	ctx = helper.NewTxContext(ctx, tx)

	userID, err := uc.createUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := uc.createAndCacheToken(ctx, *userID)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *userUsecase) createAndCacheToken(ctx context.Context, userID uuid.UUID) (*model.Token, error) {
	token, err := helper.GenerateJwtToken(userID.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    helper.DumpIncomingContext(ctx),
			"userID": helper.Dump(userID),
		}).Error(err)
		return nil, err
	}

	err = uc.tokenRepo.Create(ctx, userID, token)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":   helper.DumpIncomingContext(ctx),
			"token": helper.Dump(token),
		}).Error(err)
		return nil, err
	}

	return token, nil
}

func (uc *userUsecase) createUser(ctx context.Context, user *model.User) (*uuid.UUID, error) {
	err := uc.checkEmailAlreadyRegistered(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	user.Password, err = helper.HashPassword(user.Password)
	if err != nil {
		logrus.WithField("ctx", helper.DumpIncomingContext(ctx)).Error(err)
		return nil, err
	}

	userID, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  helper.DumpIncomingContext(ctx),
			"user": helper.Dump(user),
		}).Error(err)
		return nil, err
	}

	err = uc.userRoleRepo.Create(ctx, &model.UserRole{
		RoleID: uuid.UUID(constant.CustomerRoleID),
		UserID: *userID,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    helper.DumpIncomingContext(ctx),
			"userID": helper.Dump(userID),
		}).Error(err)
		return nil, err
	}

	return userID, nil
}

func (uc *userUsecase) checkEmailAlreadyRegistered(ctx context.Context, email string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithFields(logrus.Fields{
			"ctx":   helper.DumpIncomingContext(ctx),
			"email": email,
		}).Error(err)

		return err
	}

	if user != nil {
		return cErr.ErrorEmailAlreadyExist
	}

	return nil
}

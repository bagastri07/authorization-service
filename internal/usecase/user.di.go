package usecase

import "github.com/bagastri07/authorization-service/internal/model"

func (uc *userUsecase) InjectUserRepository(userRepo model.UserRepository) {
	uc.userRepo = userRepo
}

func (uc *userUsecase) InjectUserRoleRepository(userRoleRepo model.UserRoleRepository) {
	uc.userRoleRepo = userRoleRepo
}

func (uc *userUsecase) InjectTokenRepository(tokenRepo model.TokenRepository) {
	uc.tokenRepo = tokenRepo
}

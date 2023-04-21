package grpc

import "github.com/bagastri07/authorization-service/internal/model"

func (s *Server) InjectUserUsecase(userUC model.UserUsecase) {
	s.userUC = userUC
}

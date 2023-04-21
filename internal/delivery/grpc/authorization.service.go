package grpc

import (
	"github.com/bagastri07/authorization-service/internal/model"
	pb "github.com/bagastri07/authorization-service/pb/authorization"
)

type Server struct {
	userUC model.UserUsecase
	pb.UnimplementedProductServiceServer
}

func NewGRPCServer() *Server {
	return &Server{}
}

package grpc

import (
	"github.com/bagastri07/boilerplate-service/internal/model"
	pb "github.com/bagastri07/boilerplate-service/pb/boilerplate"
)

type Server struct {
	productUC model.ProductUsecase
	pb.UnimplementedProductServiceServer
}

func NewGRPCServer() *Server {
	return new(Server)
}

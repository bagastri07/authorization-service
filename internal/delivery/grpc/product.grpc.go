package grpc

import (
	"context"
	"time"

	cerr "github.com/bagastri07/boilerplate-service/internal/constant/customerror"
	"github.com/bagastri07/boilerplate-service/internal/model"
	pb "github.com/bagastri07/boilerplate-service/pb/boilerplate"
	"github.com/bagastri07/boilerplate-service/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	err := t.productUC.Create(ctx, &model.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       float64(req.GetPrice()),
	})

	if err != nil {
		logrus.WithContext(ctx).WithField("req", utils.Dump(req)).Error(err)
		return nil, err
	}

	return &pb.CreateProductResponse{
		Message: "created",
	}, nil
}

func (t *Server) FindByID(ctx context.Context, req *pb.FindByIDRequest) (*pb.Product, error) {
	product, err := t.productUC.FindByID(ctx, int(req.GetId()))

	switch err {
	case nil:
	case cerr.ErrorProductNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &pb.Product{
		Id:          int64(product.ID),
		Name:        product.Name,
		Description: product.Description,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339Nano),
	}, nil
}

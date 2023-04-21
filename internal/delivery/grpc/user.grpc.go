package grpc

import (
	"context"

	"github.com/bagastri07/authorization-service/internal/helper"
	"github.com/bagastri07/authorization-service/internal/model"
	pb "github.com/bagastri07/authorization-service/pb/authorization"
	"github.com/sirupsen/logrus"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user := model.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		FullName: req.GetFullName(),
	}

	tokenResp, err := s.userUC.Register(ctx, &user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": helper.DumpIncomingContext(ctx),
			"req": helper.Dump(req),
		}).Error(err)
		return nil, err
	}

	return &pb.AuthResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}, nil
}

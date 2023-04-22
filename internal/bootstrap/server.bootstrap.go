package bootstrap

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/bagastri07/authorization-service/internal/config"
	"github.com/bagastri07/authorization-service/internal/constant"
	grpcDelivery "github.com/bagastri07/authorization-service/internal/delivery/grpc"
	"github.com/bagastri07/authorization-service/internal/infrastructure"
	"github.com/bagastri07/authorization-service/internal/repository"
	"github.com/bagastri07/authorization-service/internal/usecase"
	pb "github.com/bagastri07/authorization-service/pb/authorization"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServer() {
	infrastructure.InitializePostgresConn()
	infrastructure.InitializeRedisCon()

	pgDB, err := infrastructure.PostgreSQL.DB()
	continueOrFatal(err)

	// // init repositories
	userRepository := repository.NewUserRepository(infrastructure.PostgreSQL, infrastructure.RedisClient)
	userRoleRepository := repository.NewUserRepoRepository(infrastructure.PostgreSQL)

	// // init usecases
	userUsecase := usecase.NewUserUsecase()
	userUsecase.InjectUserRepository(userRepository)
	userUsecase.InjectUserRoleRepository(userRoleRepository)

	// init grpc
	grpcDelivery := grpcDelivery.NewGRPCServer()
	grpcDelivery.InjectUserUsecase(userUsecase)

	authorizationGRPCServer := grpc.NewServer()

	pb.RegisterProductServiceServer(authorizationGRPCServer, grpcDelivery)
	if config.Env() == constant.EnvDevelopment {
		reflection.Register(authorizationGRPCServer)
	}

	lis, err := net.Listen("tcp", ":"+config.GRPCPort())
	continueOrFatal(err)

	go func() {
		err = authorizationGRPCServer.Serve(lis)
		continueOrFatal(err)
	}()

	startingMessage()
	setupPrometheus()

	wait := gracefulShutdown(context.Background(), config.GracefulShutdownTimeOut(), map[string]operation{
		"postgressql connection": func(ctx context.Context) error {
			return pgDB.Close()
		},
	})
	<-wait
}

func startingMessage() {
	logrus.Info(fmt.Sprintf("%s@%s is starting", config.ServiceName(), config.ServiceVersion()))
	logrus.Info(fmt.Sprintf("grpc server started on :%s", config.GRPCPort()))
}

func setupPrometheus() {
	http.Handle("/metrics", promhttp.Handler())

	svc := &http.Server{
		ReadTimeout:  config.MetricsReadTimeout(),
		WriteTimeout: config.MetricsWriteTimeout(),
		Addr:         fmt.Sprintf(":%s", config.MetricsPort()),
	}

	go func() {
		_ = svc.ListenAndServe()
	}()
	logrus.Info(fmt.Sprintf("metrics server started on :%s", config.MetricsPort()))
}

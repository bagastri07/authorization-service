package bootstrap

import (
	"context"
	"fmt"
	"net"
	"net/http"

	grpcTransport "github.com/bagastri07/boilerplate-service/internal/delivery/grpc"
	pb "github.com/bagastri07/boilerplate-service/pb/boilerplate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/bagastri07/boilerplate-service/internal/config"
	"github.com/bagastri07/boilerplate-service/internal/constant"
	"github.com/bagastri07/boilerplate-service/internal/infrastructure"
	"github.com/bagastri07/boilerplate-service/internal/repository"
	"github.com/bagastri07/boilerplate-service/internal/usecase"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	infrastructure.InitializePostgresConn()
	infrastructure.InitializeRedisCon()

	pgDB, err := infrastructure.PostgreSQL.DB()
	continueOrFatal(err)

	// init repositories
	productRepository := repository.NewProductRepository(infrastructure.PostgreSQL, infrastructure.RedisClient)

	// init usecases
	productUsecase := usecase.NewProductUsecase(productRepository)

	// init grpc
	grpcDelivery := grpcTransport.NewGRPCServer()
	err = grpcDelivery.InjectProductUsecase(productUsecase)
	continueOrFatal(err)

	productGRPCServer := grpc.NewServer()

	pb.RegisterProductServiceServer(productGRPCServer, grpcDelivery)
	if config.Env() == constant.EnvDevelopment {
		reflection.Register(productGRPCServer)
	}

	lis, err := net.Listen("tcp", ":"+config.GRPCPort())
	continueOrFatal(err)

	go func() {
		err = productGRPCServer.Serve(lis)
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

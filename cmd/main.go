package main

import (
	"flag"
	"github.com/mariaiu/yandex-eda-app/internal/config"
	"github.com/mariaiu/yandex-eda-app/internal/handler"
	"github.com/mariaiu/yandex-eda-app/internal/repository"
	"github.com/mariaiu/yandex-eda-app/internal/service"
	"github.com/mariaiu/yandex-eda-app/internal/validator"
	pb "github.com/mariaiu/yandex-eda-app/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.SetUpConfig(); if err != nil {
		logger.Fatalf("error initializing config: %s", err.Error())
	}

	validator := validating.NewValidator()
	if err = validator.RegisterRules(cfg.App.MaxWorkers); err != nil {
		logger.Fatalf("error when registering validation rules: %s", err.Error())
	}

	if err = validator.ValidateStruct(cfg); err != nil {
		logger.Fatalf("error validating config: %s", err.Error())
	}

	flag.Float64Var(&cfg.App.Latitude, "latitude",  cfg.App.Latitude, "set latitude")
	flag.Float64Var(&cfg.App.Longitude, "longitude",  cfg.App.Longitude, "set longitude")
	flag.Float64Var(&cfg.App.MinRating, "rating",  cfg.App.MinRating, "set minimal rating")
	flag.Parse()

	db, err := repository.NewPostgresDB(cfg); if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)

	svc := service.NewService(repo, logger, &cfg.App)

	grpcHandler := handler.NewGRPCHandler(svc, logger, &cfg.App, validator)
	httpHandler := handler.NewHTTPHandler(svc, logger, &cfg.App, &cfg.Auth, validator)

	go func() {
		httpServer := &http.Server{Addr: "0.0.0.0:8080", Handler: httpHandler.ConfigureRouter()}

		logger.Println("start HTTP server")

		logger.Fatal(httpServer.ListenAndServe())
	}()

	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:8081"); if err != nil {
			logger.Fatalln(err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterYandexEdaServer(grpcServer, grpcHandler)

		logger.Println("start GRPC server")

		logger.Fatal(grpcServer.Serve(listener))
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Println("closing")
}


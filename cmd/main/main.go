package main

import (
	"context"
	"fmt"
	"order_service/internal/config"
	"order_service/internal/repository"
	"order_service/internal/service"
	"order_service/internal/transport/grpc"
	"order_service/pkg/db/cache"
	"order_service/pkg/db/postgres"
	"order_service/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "order_service"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	cfg := config.New()
	if cfg == nil {
		panic("failed to load config")
	}

	redis := cache.New(cfg.RedisConfig)
	fmt.Println(redis.Ping(ctx))

	db, err := postgres.New(cfg.Config)
	if err != nil {
		panic(err)
	}
	repo := repository.NewOrderRepository(db)
	service := service.NewOrderService(repo)

	grpcServer, err := grpc.New(ctx, cfg.GRPCServerPort, cfg.RestServerPort, service)
	if err != nil {
		mainLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера
	go func() {
		if err = grpcServer.Start(ctx); err != nil {
			mainLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	err = grpcServer.Stop(ctx)
	if err != nil {
		mainLogger.Error(ctx, err.Error())
	}
	mainLogger.Info(ctx, "Server stopped")
	fmt.Println("Server stopped")
}

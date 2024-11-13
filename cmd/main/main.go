package main

import (
	"context"
	"fmt"
	"order_service/internal/models"
	"order_service/internal/repository"
	"order_service/internal/service"
	"order_service/internal/transport/grpc"
	"order_service/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	serviceName = "order_service"
	grpcPort    = 50051
	restPort    = 8080
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	mu := sync.RWMutex{}
	db := make(map[string]models.Order)
	repo := repository.NewOrderRepository(db, &mu)
	service := service.NewOrderService(repo)

	grpcServer, err := grpc.New(ctx, grpcPort, restPort, service)
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

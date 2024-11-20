package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	client "order_service/pkg/api/3_1"
	"order_service/pkg/logger"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *http.Server
	listener   net.Listener
}

func New(ctx context.Context, grpcPort, restPort int, service Service) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	client.RegisterOrderServiceServer(grpcServer, NewOrderService(service))

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	gwmux := runtime.NewServeMux()
	if err := client.RegisterOrderServiceHandler(context.Background(), gwmux, conn); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: gwmux,
	}

	return &Server{grpcServer, gwServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		lg := logger.GetLoggerFromCtx(ctx)
		if lg != nil {
			lg.Info(ctx, "starting grpc server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
		}

		return s.grpcServer.Serve(s.listener)
	})

	eg.Go(func() error {
		lg := logger.GetLoggerFromCtx(ctx)
		if lg != nil {
			lg.Info(ctx, "starting rest server", zap.String("port", s.restServer.Addr))
		}

		return s.restServer.ListenAndServe()
	})

	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()

	l := logger.GetLoggerFromCtx(ctx)
	if l != nil {
		l.Info(ctx, "grpc server stopped")
	}

	return s.restServer.Shutdown(ctx)
}

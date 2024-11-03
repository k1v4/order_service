package grpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"order_service/pkg/api/order"
	"order_service/pkg/logger"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(ctx context.Context, port int, service Service) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	order.RegisterOrderServiceServer(grpcServer, NewOrderService(service))

	return &Server{grpcServer, lis}, nil
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

	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	l := logger.GetLoggerFromCtx(ctx)
	if l != nil {
		l.Info(ctx, "grpc server stopped")
	}

	return nil
}

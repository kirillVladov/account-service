package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/kirillVladov/account-service/internal/gen/grpc"
)

type Handlers struct {
	Account pb.AccountServiceServer
}

type Server struct {
	server *grpc.Server
	port   int
}

func NewServer(port int, handlers Handlers) *Server {
	grpcServer := grpc.NewServer()

	pb.RegisterAccountServiceServer(grpcServer, handlers.Account)

	return &Server{
		server: grpcServer,
		port:   port,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/kirillVladov/account-service/internal/debug_server"
	"github.com/kirillVladov/account-service/internal/di"
	"github.com/kirillVladov/account-service/internal/transport/grpc"
	"github.com/kirillVladov/account-service/pkg/logger"
)

type shutdownFn = func(ctx context.Context) error

type App struct {
	shutdownFn []shutdownFn

	di     *di.DI
	logger *zap.Logger
}

const (
	shutdownTimeout = 30 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	container := di.New()
	l := container.Logger()

	a := &App{
		di:     container,
		logger: l,
	}

	ctx = logger.WithLogger(ctx, l)

	l.Info("service started")

	go a.startDebugServer()
	go a.startGRPCServer(ctx)

	<-ctx.Done()

	a.logger.Info("shutting down...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	for _, fn := range a.shutdownFn {
		fn(shutdownCtx)
	}

	a.logger.Info("stopped")
}

func (a *App) startDebugServer() {
	a.logger.Info("debug server started")

	server := debug_server.New(fmt.Sprintf(":%d", a.di.Config().DebugPort), a.logger)
	_ = server.Start()

	a.shutdownFn = append(a.shutdownFn, server.Shutdown)
}

func (a *App) startGRPCServer(ctx context.Context) {
	a.logger.Info("grpc server started")

	cfg := a.di.Config()

	handlers := grpc.Handlers{
		Account: a.di.AccountHandler(),
	}

	server := grpc.NewServer(cfg.GRPCPort, handlers)

	go func() {
		if err := server.Start(); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	a.shutdownFn = append(a.shutdownFn, server.Shutdown)
}

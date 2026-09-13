package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kirillVladov/account-service/internal/debug_server"
	"github.com/kirillVladov/account-service/internal/di"
	"github.com/kirillVladov/account-service/pkg/logger"
)

type shutdownFn = func(ctx context.Context) error

type App struct {
	shutdownFn []shutdownFn

	di *di.DI
}

// todo: review
func main() {
	ctx := context.Background()

	app := App{
		di: di.New(),
	}

	ctx = logger.WithLogger(ctx, app.di.Logger())

	go app.startDebugServer()
	go app.startAccountConfirmationQueue(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for _, fn := range app.shutdownFn {
		if err := fn(ctx); err != nil {
			log.Print("forced shutdown")
			os.Exit(1)
		}
	}

	log.Print("service stopped")
}

func (a *App) startDebugServer() {
	server := debug_server.New(fmt.Sprintf(":%d", a.di.Config().DebugPort), a.di.Logger())
	_ = server.Start()

	a.shutdownFn = append(a.shutdownFn, server.Shutdown)
}

func (a *App) startAccountConfirmationQueue(ctx context.Context) {
	q := a.di.AccountConfirmationQueue()

	q.Start(ctx)

	a.shutdownFn = append(a.shutdownFn, q.Stop)
}

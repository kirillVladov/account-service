package debug_server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	echo  *echo.Echo
	ready func() bool
	log   *zap.Logger
}

// New создаёт debug Server на указанном addr (например ":9090").
func New(addr string, logger *zap.Logger) *Server {
	s := &Server{
		log: zap.NewNop(),
	}

	if logger != nil {
		s.log = logger
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				s.log.Error("debug request",
					zap.String("method", v.Method),
					zap.String("uri", v.URI),
					zap.Int("status", v.Status),
					zap.Error(v.Error),
				)
				return nil
			}
			s.log.Debug("debug request",
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/liveness", s.handleLiveness)
	e.GET("/readiness", s.handleReadiness)

	// todo: refactor
	e.Server.Addr = addr
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 5 * time.Second
	e.Server.IdleTimeout = 30 * time.Second

	s.echo = e
	return s
}

func (s *Server) Start() error {
	s.log.Info("debug server starting", zap.String("addr", s.echo.Server.Addr))

	go func() {
		if err := s.echo.Start(s.echo.Server.Addr); err != nil && err != http.ErrServerClosed {
			s.log.Error("debug server error", zap.Error(err))
		}
	}()

	return nil
}

// Shutdown корректно останавливает сервер.
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("debug server shutting down")
	return s.echo.Shutdown(ctx)
}

func (s *Server) handleLiveness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (s *Server) handleReadiness(c echo.Context) error {
	if s.ready != nil && !s.ready() {
		return c.String(http.StatusServiceUnavailable, "not ready")
	}
	return c.String(http.StatusOK, "ok")
}

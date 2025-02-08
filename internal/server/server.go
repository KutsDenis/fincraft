package server

import (
	"context"
	"errors"
	"github.com/KutsDenis/logzap"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

type Server struct {
	httpServer *http.Server
}

// NewServer создает HTTP-сервер
func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
	}
}

// Start запускает HTTP-сервер
func (s *Server) Start() {
	go func() {
		logzap.Info("starting server", zap.String("port", s.httpServer.Addr))

		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logzap.Fatal("failed to start server", zap.Error(err))
		}
	}()
}

// GracefulShutdown останавливает HTTP-сервер
func (s *Server) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logzap.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logzap.Error("failed to shutdown server", zap.Error(err))
	}
	logzap.Info("server shutdown")
}

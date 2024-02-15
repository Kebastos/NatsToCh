package server

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
	logger *log.Log
	cfg    *config.HTTPConfig
}

func NewHTTPServer(cfg *config.HTTPConfig, logger *log.Log) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		logger: logger,
		server: &http.Server{
			Addr:         cfg.ListenAddr,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *HTTPServer) Serve() {
	http.Handle("/metrics", promhttp.Handler())

	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Fatalf("failed to start http server. %s", err)
	}
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

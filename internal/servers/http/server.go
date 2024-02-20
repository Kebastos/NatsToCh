package http

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Server struct {
	server *http.Server
	logger *log.Log
	cfg    *config.HTTPConfig
}

func NewServer(cfg *config.HTTPConfig, logger *log.Log) *Server {
	return &Server{
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

func (s *Server) Serve() {
	http.Handle("/metrics", promhttp.Handler())

	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Fatalf("failed to start http server. %s", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

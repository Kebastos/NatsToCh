package main

import (
	"context"
	"flag"
	"github.com/Kebastos/NatsToCh/internal/clickhouse"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/Kebastos/NatsToCh/internal/nats"
	"github.com/Kebastos/NatsToCh/internal/servers/http"
	"github.com/Kebastos/NatsToCh/internal/servers/nats2ch"
	"os/signal"
	"syscall"
)

var Version = "0.0.0"
var configFile = flag.String("config", "config_example.yaml", "Proxy configuration filename")

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := log.MustConfig()

	cfg, err := config.NewConfig(*configFile)
	if err != nil {
		logger.Fatalf("failed to read config. %s", err)
	}
	if cfg.Debug {
		logger.SetDebug(cfg.Debug)
		logger.Debugf("debug mode enabled")
	}

	m, err := metrics.NewMetrics(cfg)
	if err != nil {
		logger.Fatalf("failed to create metrics: %s", err)
	}
	m.MustRegister()

	httpServer := http.NewServer(&cfg.HTTPConfig, logger)
	go httpServer.Serve()

	natsClient := nats.NewClient(&cfg.NATSConfig, logger, m)
	err = natsClient.Connect()
	if err != nil {
		logger.Fatalf("failed to connect to NATS server. %s", err)
	}

	chClient := clickhouse.NewClickhouseClient(&cfg.CHConfig, logger, m)
	if err = chClient.Connect(); err != nil {
		logger.Fatalf("failed to connect to ClickHouse. %s", err)
	}

	n := nats2ch.NewServer(cfg, natsClient, chClient, logger, m)
	err = n.Start(ctx)
	if err != nil {
		logger.Fatalf("failed to start NATS to ClickHouse server. %s", err)
	}

	logger.Infof("application started. version - %s", Version)

	<-ctx.Done()
	logger.Infof("shutting down...")
	natsClient.Shutdown()
	err = chClient.Close()
	if err != nil {
		logger.Fatalf("failed to close ClickHouse. %s", err)
	}
	err = httpServer.Shutdown(ctx)
	if err != nil {
		logger.Fatalf("failed to shutdown http server. %s", err)
	}
}

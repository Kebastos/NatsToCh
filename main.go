package main

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/Kebastos/NatsToCh/internal/server"
	"github.com/Kebastos/NatsToCh/internal/workers"
	"os/signal"
	"syscall"
)

var Version = "0.0.0"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.MustConfig()

	logger.Infof("loading config...")
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalf("failed to read config. %s", err)
	}
	if cfg.Debug {
		logger.SetDebug(cfg.Debug)
		logger.Debugf("debug mode run")
	}

	metrics.MustRegister()

	httpServer := server.NewHTTPServer(&cfg.HTTPConfig, logger)
	go httpServer.Serve()

	natsClient := clients.NewNatsClient(&cfg.NATSConfig, logger)
	err = natsClient.Connect()
	if err != nil {
		logger.Fatalf("failed to connect to NATS server. %s", err)
	}

	chClient := clients.NewClickhouseClient(&cfg.CHConfig, logger)
	if err = chClient.Connect(); err != nil {
		logger.Fatalf("failed to connect to ClickHouse. %s", err)
	}

	natsWorker := workers.NewNatsWorker(cfg, natsClient, chClient)
	if err = natsWorker.Start(ctx); err != nil {
		logger.Fatalf("failed to start nats worker. %s", err)
	}

	go func() {
		<-ctx.Done()
		logger.Infof("shutting down...")
		natsClient.Shutdown()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			logger.Fatalf("failed to shutdown http server. %s", err)
		}
	}()
	logger.Infof("application started. version - %s", Version)
	select {}
}

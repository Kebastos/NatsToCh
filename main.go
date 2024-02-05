package main

import (
	"context"
	"github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/config"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/Kebastos/NatsToCh/internal/workers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx); err != nil {
		log.Fatalf("failed to run server: %s", err)
	}
}

func runServer(ctx context.Context) error {
	log.Infof("loading config...")
	cfg := config.MustConfig()

	if cfg.Server.Debug {
		log.SetDebug(cfg.Server.Debug)
		log.Debugf("debug mode run")
	}

	metrics.MustRegister(cfg)
	log.Infof("metrics registered")

	server := cfg.Server
	if len(server.HTTP.ListenAddr) == 0 {
		panic("wrong config section - `listen_addr` is not configured")
	}

	go serveHTTP(server.HTTP)

	log.Infof("http server is starting at address: %s", cfg.Server.HTTP.ListenAddr)

	natsClient := clients.NewNatsClient(&cfg.NATSConfig)
	err := natsClient.Connect()
	if err != nil {
		log.Fatalf("failed to connect to NATS server. %s", err)
	}

	chClient := clients.NewClickhouseClient(&cfg.CHConfig)
	if err = chClient.Connect(); err != nil {
		log.Fatalf("failed to connect to ClickHouse. %s", err)
	}

	natsWorker := workers.NewNatsWorker(cfg, natsClient, chClient)
	if err = natsWorker.Start(ctx); err != nil {
		log.Fatalf("failed to start nats worker. %s", err)
	}

	log.Infof("application start")

	<-ctx.Done()
	log.Infof("shutting down server gracefully")

	select {}
}

func serveHTTP(cfg config.HTTP) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(cfg.ListenAddr, nil)
	if err != nil {
		log.Fatalf("failed to start http server. %s", err)
	}
}

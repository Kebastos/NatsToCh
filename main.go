package main

import (
	"github.com/Kebastos/NatsToCh/config"
	clients2 "github.com/Kebastos/NatsToCh/internal/clients"
	"github.com/Kebastos/NatsToCh/internal/log"
	"github.com/Kebastos/NatsToCh/internal/metrics"
	"github.com/Kebastos/NatsToCh/internal/workers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
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

	natsClient := clients2.NewNatsClient(&cfg.NATSConfig)
	err := natsClient.Connect()
	if err != nil {
		log.Fatalf("failed to connect to NATS server. %s", err)
	}

	chClient := clients2.NewClickhouseClient(&cfg.CHConfig)
	if err = chClient.Connect(); err != nil {
		log.Fatalf("failed to connect to ClickHouse. %s", err)
	}

	natsWorker := workers.NewNatsWorker(cfg, natsClient, chClient)
	if err = natsWorker.Start(); err != nil {
		log.Fatalf("failed to start nats worker. %s", err)
	}

	log.Infof("application start")

	select {}
}

func serveHTTP(cfg config.HTTP) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(cfg.ListenAddr, nil)
	if err != nil {
		log.Fatalf("failed to start http server. %s", err)
	}
}

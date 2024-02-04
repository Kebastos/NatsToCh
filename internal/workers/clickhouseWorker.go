package workers

import (
	"github.com/Kebastos/NatsToCh/config"
	"github.com/Kebastos/NatsToCh/internal/clients"
)

type ClickhouseWorker struct {
	cfg *config.Server
	ch  *clients.ClickhouseClient
}

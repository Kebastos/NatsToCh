package workers

import (
	"github.com/Kebastos/NatsToCh/clients"
	"github.com/Kebastos/NatsToCh/config"
)

type ClickhouseWorker struct {
	cfg *config.Server
	ch  *clients.ClickhouseClient
}

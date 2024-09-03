package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Debug      bool       `yaml:"debug"`
	HTTPConfig HTTPConfig `yaml:"http" env-required:"true"`
	NATSConfig NATSConfig `yaml:"nats" env-required:"true"`
	CHConfig   CHConfig   `yaml:"clickhouse" env-required:"true"`
	Subjects   []Subject  `yaml:"subjects" env-required:"true"`
}

type HTTPConfig struct {
	ListenAddr   string        `yaml:"listen_addr" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"5s"`
}

type CHConfig struct {
	Host            string        `yaml:"host" env-required:"true"`
	Port            int           `yaml:"port" env-required:"true"`
	Database        string        `yaml:"db" env-required:"true"`
	User            string        `yaml:"user" env-required:"true"`
	Password        string        `yaml:"password"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"10s"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"10"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"10"`
}

type NATSConfig struct {
	ClientName     string `yaml:"client_name" env-required:"true"`
	User           string `yaml:"user" env-required:"true"`
	Password       string `yaml:"password" env-required:"true"`
	Server         string `yaml:"server" env-required:"true"`
	MaxReconnect   int    `yaml:"max_reconnect" env-default:"5"`
	ReconnectWait  int    `yaml:"reconnect_wait" env-default:"5"`
	ConnectTimeout int    `yaml:"connect_timeout" env-default:"5"`
	MaxWait        int    `yaml:"max_wait" env-default:"10"`
}

type Subject struct {
	Name         string       `yaml:"name" env-required:"true"`
	Queue        string       `yaml:"queue"`
	TableName    string       `yaml:"table_name" env-required:"true"`
	UseBuffer    bool         `yaml:"use_buffer" env-default:"true"`
	BufferConfig BufferConfig `yaml:"buffer_config"`
}

type BufferConfig struct {
	MaxSize int           `yaml:"max_size" env-required:"true"`
	MaxWait time.Duration `yaml:"max_wait" env-required:"true"`
}

func NewConfig(configFile string) (*Config, error) {
	flag.Parse()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf(fmt.Sprintf("config file doesn't exist: %s", configFile))
	}

	var config Config

	if err := cleanenv.ReadConfig(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

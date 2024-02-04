package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Server     Server
	NATSConfig NATSConfig `yaml:"nats"`
	CHConfig   CHConfig   `yaml:"clickhouse"`
	Subjects   []Subject  `yaml:"subjects"`
}

type Server struct {
	Env       string `yaml:"env" env:"ENV" envDefault:"local"`
	Namespace string `yaml:"namespace" env:"NAMESPACE" envDefault:"natstoch"`
	Debug     bool   `yaml:"debug" env:"DEBUG" envDefault:"false"`
	HTTP      HTTP   `yaml:"http"`
}

type HTTP struct {
	ListenAddr string `yaml:"listen_addr,omitempty"`
}

type CHConfig struct {
	Host        string        `yaml:"host" env:"CH_HOST" envDefault:"localhost"`
	Port        uint16        `yaml:"port" env:"CH_PORT" envDefault:"9000"`
	Database    string        `yaml:"db" env:"CH_DATABASE" envDefault:"default"`
	User        string        `yaml:"user" env:"CH_USER" envDefault:"default"`
	Password    string        `yaml:"password" env:"CH_PASSWORD" envDefault:"default"`
	Timeout     time.Duration `yaml:"timeout" env:"CH_TIMEOUT" envDefault:"2s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"CH_IDLE_TIMEOUT" envDefault:"30s"`
}

type NATSConfig struct {
	ClientName             string `yaml:"client_name" env:"NATS_CLIENT_NAME" envDefault:"natsToCh"`
	User                   string `yaml:"user" env:"NATS_USER" envDefault:"localUser"`
	Password               string `yaml:"password" env:"NATS_PASSWORD" envDefault:"localPassword"`
	Server                 string `yaml:"server" env:"NATS_SERVER" envDefault:"localhost"`
	MaxReconnect           int    `yaml:"max_reconnect" env:"NATS_MAX_RECONNECT" envDefault:"-1"`
	ReconnectWait          int    `yaml:"reconnect_wait" env:"NATS_RECONNECT_WAIT" envDefault:"10"`
	ConnectTimeout         int    `yaml:"connect_timeout" env:"NATS_CONNECT_TIMEOUT" envDefault:"10"`
	MaxWait                int    `yaml:"max_wait" env:"NATS_MAX_WAIT" envDefault:"10"`
	PublishAsyncMaxPending int    `yaml:"publish_async_max_pending" env:"NATS_PUBLISH_ASYNC_MAX_PENDING" envDefault:"10"`
}

type Subject struct {
	Name         string       `yaml:"name" `
	Queue        string       `yaml:"queue"`
	UseBuffer    bool         `yaml:"use_buffer"`
	TableName    string       `yaml:"table_name"`
	BufferConfig BufferConfig `yaml:"buffer_config,omitempty"`
}

type BufferConfig struct {
	MaxSize     int `yaml:"max_size,omitempty"`
	MaxWait     int `yaml:"max_wait,omitempty"`
	MaxByteSize int `yaml:"max_byte_size,omitempty"`
}

func MustConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file doesn't exist: %s", configPath))
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic(fmt.Sprintf("can't read config file: %s", err))
	}

	return &config
}

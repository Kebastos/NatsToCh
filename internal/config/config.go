package config

import (
	"flag"
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
	Env       string `yaml:"env"`
	Namespace string `yaml:"namespace"`
	Debug     bool   `yaml:"debug"`
	HTTP      HTTP   `yaml:"http"`
}

type HTTP struct {
	ListenAddr string `yaml:"listen_addr,omitempty"`
}

type CHConfig struct {
	Host            string        `yaml:"host"`
	Port            uint16        `yaml:"port"`
	Database        string        `yaml:"db"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
}

type NATSConfig struct {
	ClientName             string `yaml:"client_name"`
	User                   string `yaml:"user"`
	Password               string `yaml:"password"`
	Server                 string `yaml:"server"`
	MaxReconnect           int    `yaml:"max_reconnect"`
	ReconnectWait          int    `yaml:"reconnect_wait"`
	ConnectTimeout         int    `yaml:"connect_timeout"`
	MaxWait                int    `yaml:"max_wait"`
	PublishAsyncMaxPending int    `yaml:"publish_async_max_pending"`
}

type Subject struct {
	Name              string            `yaml:"name" `
	Queue             string            `yaml:"queue"`
	TableName         string            `yaml:"table_name"`
	Async             bool              `yaml:"async"`
	AsyncInsertConfig AsyncInsertConfig `yaml:"async_insert_config,omitempty"`
	UseBuffer         bool              `yaml:"use_buffer"`
	BufferConfig      BufferConfig      `yaml:"buffer_config,omitempty"`
}

type AsyncInsertConfig struct {
	Wait bool `yaml:"wait,omitempty"`
}

type BufferConfig struct {
	MaxSize     int           `yaml:"max_size,omitempty"`
	MaxWait     time.Duration `yaml:"max_wait,omitempty"`
	MaxByteSize int           `yaml:"max_byte_size,omitempty"`
}

var configFile = flag.String("config", "", "Proxy configuration filename")

func MustConfig() *Config {
	flag.Parse()

	if *configFile == "" {
		*configFile = "config/local.yaml"
	}

	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file doesn't exist: %s", *configFile))
	}

	var config Config

	if err := cleanenv.ReadConfig(*configFile, &config); err != nil {
		panic(fmt.Sprintf("can't read config file: %s", err))
	}

	return &config
}
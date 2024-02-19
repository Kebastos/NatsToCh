package config_test

import (
	"github.com/Kebastos/NatsToCh/internal/config"
	"testing"
)

var (
	configFile      = "../../config/dev.yaml"
	emptyConfigFile = ""
)

func TestConfigNewWithUserFile(t *testing.T) {
	_, err := config.NewConfig(configFile)
	if err != nil {
		t.Errorf("failed to read config. %s", err)
	}
}

func TestConfigNewWithEmptyFile(t *testing.T) {
	_, err := config.NewConfig(emptyConfigFile)
	if err == nil {
		t.Errorf("should have failed to create config. %s", err)
	}
}

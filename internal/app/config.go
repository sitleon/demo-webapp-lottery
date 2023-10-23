package app

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type HttpServerConfig struct {
	Port int `yaml:"port"`
}

type LotteryConfig struct {
	Interval int64 `yaml:"interval"`
}

type AppConfig struct {
	Http struct {
		Server HttpServerConfig `yaml:"server"`
	} `yaml:"http"`
	Db struct {
		ConnURI string `yaml:"connection_uri"`
	} `yaml:"db"`
	Lottery LotteryConfig `yaml:"lottery"`
}

func LoadCfg(path string) *AppConfig {
	f, err := os.Open(path)
	if err != nil {
		logrus.Errorf("failed to load configuration file [%s]: %s", path, err)
		return nil
	}

	b, err := io.ReadAll(f)
	if err != nil {
		logrus.Errorf("failed to load configuration file [%s]: %s", path, err)
		return nil
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		logrus.Errorf("invalid configurations [%s]: %s", path, err)
		return nil
	}

	return &cfg
}

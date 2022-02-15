package config

import (
	"github.com/Aserose/ModusOperandi/pkg/logger"
)

type ClientConfig struct {
	ymlFilename string
	log         logger.Logger
}

type CfgClient struct {
	Name string `yaml:"name"`
	Main struct {
		WindowHeight float32 `yaml:"height"`
		WindowWidth  float32 `yaml:"width"`
	} `yaml:"main_window"`
	Shortcut struct {
		Name         string  `yaml:"name"`
		WindowHeight float32 `yaml:"height"`
		WindowWidth  float32 `yaml:"width"`
	} `yaml:"shortcut_window"`
	Notification struct {
		Name         string  `yaml:"name"`
		WindowHeight float32 `yaml:"height"`
		WindowWidth  float32 `yaml:"width"`
	} `yaml:"notification_window"`
	Creation struct {
		Name         string  `yaml:"name"`
		WindowHeight float32 `yaml:"height"`
		WindowWidth  float32 `yaml:"width"`
	} `yaml:"creation_window"`
}

func NewClientCfg(clientFilename string, log logger.Logger) *ClientConfig {
	return &ClientConfig{
		ymlFilename: clientFilename,
		log:         log,
	}
}

func (c ClientConfig) InitCfgClient() CfgClient {
	var cfg CfgClient

	unmarshalYaml(c.ymlFilename, c.log, &cfg)

	return cfg
}

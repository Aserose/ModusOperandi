package config

import (
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Infrastructure interface {
	InitInfrastructure() CfgBolt
}

type Client interface {
	InitCfgClient() CfgClient
}

type Config struct {
	Infrastructure
	Client
}

func NewConfig(ymlInfrastuctureFilename string, clientFilename string, log logger.Logger) *Config {
	return &Config{
		Infrastructure: NewInfrastructure(ymlInfrastuctureFilename, log),
		Client:         NewClientCfg(clientFilename, log),
	}
}

func unmarshalYaml(filename string, log logger.Logger, outs ...interface{}) {
	ymlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("%s: %s", log.CallInfoStr(), err.Error())
	}

	for _, out := range outs {
		err = yaml.Unmarshal(ymlFile, out)
		if err != nil {
			log.Panicf("%s: %s", log.CallInfoStr(), err.Error())
		}
	}
}

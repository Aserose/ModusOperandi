package config

import (
	"github.com/Aserose/ModusOperandi/pkg/logger"
)

type InfrastructureConfig struct {
	ymlFilename string
	log         logger.Logger
}

type (
	CfgBolt struct {
		RootName              string `yaml:"rootName"`
		BucketInstructionName string `yaml:"bucketInstructionName"`
	}
)

func NewInfrastructure(infrastructureFilename string, log logger.Logger) *InfrastructureConfig {
	return &InfrastructureConfig{
		ymlFilename: infrastructureFilename,
		log:         log,
	}
}

func (i InfrastructureConfig) InitInfrastructure() CfgBolt {
	var cfgBolt CfgBolt

	unmarshalYaml(i.ymlFilename, i.log, &cfgBolt)

	return cfgBolt
}

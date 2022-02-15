package app

import (
	"github.com/Aserose/ModusOperandi/pkg/client"
	"github.com/Aserose/ModusOperandi/pkg/client/handler"
	"github.com/Aserose/ModusOperandi/pkg/config"
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"github.com/Aserose/ModusOperandi/pkg/repository"
	"github.com/Aserose/ModusOperandi/pkg/repository/boltDB"
	"github.com/Aserose/ModusOperandi/pkg/repository/boltDB/data"
	"github.com/Aserose/ModusOperandi/pkg/service"
)

const (
	YmlInfrastructureFilename = `configs/infrastructure.yml`
	YmlClientFilename         = `configs/client.yml`
)

func Start() {
	log := logger.NewLogger()

	cfg := config.NewConfig(YmlInfrastructureFilename, YmlClientFilename, log)
	cfgBolt := cfg.InitInfrastructure()
	cfgClient := cfg.Client.InitCfgClient()

	bolt := data.NewBoltData(boltDB.ConnectBoltDB(log, cfgBolt), cfgBolt, log)
	db := repository.NewDB(bolt)

	services := service.NewService(db, log)

	handlers := handler.NewHandler(services, log)

	clt := client.NewClient(handlers, cfgClient)

	clt.Start()
}

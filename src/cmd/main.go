package main

import (
	"os"

	"github.com/alpardfm/e-payment/src/business/domain"
	"github.com/alpardfm/e-payment/src/business/usecase"
	rest "github.com/alpardfm/e-payment/src/handler"
	"github.com/alpardfm/e-payment/src/utils/config"
	"github.com/alpardfm/go-toolkit/configbuilder"
	"github.com/alpardfm/go-toolkit/configreader"
	"github.com/alpardfm/go-toolkit/files"
	"github.com/alpardfm/go-toolkit/log"
	"github.com/alpardfm/go-toolkit/parser"
	"github.com/alpardfm/go-toolkit/sql"
)

const (
	configFile = "./etc/cfg/config.json"
	appName    = "E-Payment API"
)

func main() {
	if !files.IsExist(configFile) {
		configbuilder.Init(configbuilder.Options{
			Env:        os.Getenv("EC_APP_ENVIRONMENT"),
			Key:        os.Getenv("EC_APP_KEY"),
			Secret:     os.Getenv("EC_APP_SECRET"),
			Region:     os.Getenv("EC_APP_REGION"),
			ConfigFile: configFile,
			Namespace:  appName,
		}).BuildConfig()
	}

	cfg := config.Init()
	configreader := configreader.Init(configreader.Options{
		ConfigFile: configFile,
	})
	configreader.ReadConfig(&cfg)

	log := log.Init(cfg.Log)

	parser := parser.InitParser(log, cfg.Parser)

	JSONParser := parser.JSONParser()

	db := sql.Init(cfg.SQL, log)

	d := domain.Init(log, db, JSONParser, cfg)

	uc := usecase.Init(log, d, JSONParser, cfg)

	r := rest.Init(cfg, configreader, log, parser.JSONParser(), uc)
	r.Run()
}

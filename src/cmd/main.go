package main

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase"
	"go-clean/src/handler/rest"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/cloud_storage"
	"go-clean/src/lib/configreader"
	"go-clean/src/lib/midtrans"
	"go-clean/src/lib/sql"
	"go-clean/src/utils/config"

	_ "go-clean/docs/swagger"
)

// @contact.name   Rakhmad Giffari Nurfadhilah
// @contact.url    https://fadhilmail.tech/
// @contact.email  rakhmadgiffari14@gmail.com

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization

const (
	configFile       string = "./etc/cfg/config.json"
	cloudStorageFile string = "./etc/cfg/officebuddy-388208-2a9032588b19.json"
)

func main() {
	cfg := config.Init()
	configReader := configreader.Init(configreader.Options{
		ConfigFile: configFile,
	})
	configReader.ReadConfig(&cfg)

	configReader = configreader.Init(configreader.Options{
		ConfigFile: cloudStorageFile,
	})
	configReader.ReadConfig(&cfg.GoogleStorage)

	cs := cloud_storage.Init(cloudStorageFile, cfg.GoogleStorage)

	auth := auth.Init()

	db := sql.Init(cfg.SQL)

	midtrans := midtrans.Init(cfg.Midtrans)

	d := domain.Init(db, midtrans, cs)

	uc := usecase.Init(auth, d, cs)

	r := rest.Init(cfg.Gin, configReader, uc, auth)

	r.Run()
}

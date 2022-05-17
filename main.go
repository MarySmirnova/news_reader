package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/MarySmirnova/news_reader/internal"
	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/caarlos0/env/v6"
	"github.com/chatex-com/process-manager"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var cfg config.Application

func init() {
	godotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	rssConf, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(rssConf, &cfg.RSS); err != nil {
		panic(err)
	}

	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(lvl)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})

	process.SetLogger(internal.NewProcessLogger())
}

func main() {
	app, err := internal.NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	app.Run()
}

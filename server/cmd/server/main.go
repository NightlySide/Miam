package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/api"
	"io.github.nightlyside/miam/pkg/config"
)

var ConfigPathFlag = flag.String("config", "", "config file path")

func main() {
	flag.Parse()

	// load config
	cfg, err := config.LoadConfig(*ConfigPathFlag)
	if err != nil {
		log.Fatal(err)
	}

	// create router
	router := api.CreateRouter(cfg)

	// start server
	listenOn := fmt.Sprintf("%s:%d", cfg.Api.Host, cfg.Api.Port)
	log.Infof("Starting server of %s", listenOn)
	http.ListenAndServe(listenOn, router)
}

package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/api"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
)

var ConfigPathFlag = flag.String("config", "", "config file path")

func main() {
	flag.Parse()

	// load config
	cfg, err := config.LoadConfig(*ConfigPathFlag)
	if err != nil {
		log.Fatal(err)
	}

	// create handler
	conn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	handler := &api.Handler{
		Config: cfg,
		DB:     conn,
	}

	// create router and server
	router := api.CreateRouter(handler)
	listenOn := fmt.Sprintf("%s:%d", cfg.Api.Host, cfg.Api.Port)
	srv := &http.Server{
		Addr:    listenOn,
		Handler: router,
	}

	// start server
	log.Infof("Starting server on %s", listenOn)
	log.Fatal(srv.ListenAndServe())
}

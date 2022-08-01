package main

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig(filepath.Join("..", "config", "config.toml"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// start the migration
	if err = conn.AutoMigrate(&db.Recipe{}); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"math"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
	"io.github.nightlyside/miam/pkg/utils"
)

var Searchflag = flag.String("search", "", "search food in list")

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig(filepath.Join("..", "config", "config.toml"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	names, err := conn.GetFoodNames("fr")
	if err != nil {
		log.Fatal(err)
	}

	matches := utils.Autocomplete(*Searchflag, names)
	for _, match := range matches[:int(math.Min(10., float64(len(matches))))] {
		log.Info(match)
	}
}

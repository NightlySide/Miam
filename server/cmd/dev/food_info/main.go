package main

import (
	"flag"
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

	match := utils.Autocomplete(*Searchflag, names)[0]

	food := conn.FoodFromName(match)
	log.Info("\n" + food.Info())
}

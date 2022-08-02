package main

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig(filepath.Join("..", "config", "config.toml"))
	if err != nil {
		logrus.Fatal(err)
	}

	conn, err := db.ConnectDB(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	var recipe db.Recipe
	err = conn.Model(&recipe).Preload("Ingredients").Preload("Steps").Preload("Diet").First(&recipe, "id = 2").Error
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("\n%s", recipe.Info())
}

package db

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io.github.nightlyside/miam/pkg/ciqual"
	"io.github.nightlyside/miam/pkg/config"
)

type Database struct {
	gorm.DB
	Cache cache.Cache
}

func ConnectDB(cfg *config.Config) (*Database, error) {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Paris", cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Database, cfg.Database.Port)
	log.Debugf("DSN: %s", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger:               logger.Default.LogMode(logger.Silent),
		FullSaveAssociations: true,
	})
	if err != nil {
		return nil, err
	}

	return &Database{
		DB:    *db,
		Cache: *cache.New(10*time.Minute, 20*time.Minute),
	}, nil
}

func (db *Database) CiqualMigrate() error {
	return db.AutoMigrate(ciqual.Food{}, ciqual.Component{}, ciqual.Composition{}, ciqual.FoodGroup{}, ciqual.Source{})
}

func (db *Database) RecipeMigrate() error {
	return db.AutoMigrate(Recipe{}, Diet{}, Ingredient{}, Step{})
}

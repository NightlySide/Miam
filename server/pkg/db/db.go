package db

import (
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io.github.nightlyside/miam/pkg/ciqual"
	"io.github.nightlyside/miam/pkg/config"
)

type Database struct {
	gorm.DB
	Redis redis.Client
}

func ConnectDB(cfg *config.Config) (*Database, error) {
	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Paris", cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Database, cfg.Database.Port)
	log.Debugf("DSN: %s", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger:               logger.Default.LogMode(logger.Silent),
		FullSaveAssociations: true,
	})
	if err != nil {
		return nil, err
	}

	// Connect to redis
	redc := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	return &Database{
		DB:    *db,
		Redis: *redc,
	}, nil
}

func (db *Database) CiqualMigrate() error {
	return db.AutoMigrate(ciqual.Food{}, ciqual.Component{}, ciqual.Composition{}, ciqual.FoodGroup{}, ciqual.Source{})
}

func (db *Database) RecipeMigrate() error {
	return db.AutoMigrate(Recipe{}, Diet{}, Ingredient{}, Step{})
}

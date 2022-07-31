package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io.github.nightlyside/miam/pkg/ciqual"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/miam?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(ciqual.Food{}, ciqual.Component{}, ciqual.Composition{}, ciqual.FoodGroup{}, ciqual.Source{})
}

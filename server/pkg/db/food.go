package db

import (
	"fmt"

	"gorm.io/gorm"
	"io.github.nightlyside/miam/pkg/ciqual"
)

func GetFoodNames(db *gorm.DB, lang string) ([]string, error) {
	var selector string
	if lang == "fr" {
		selector = "name_fr"
	} else if lang == "eng" {
		selector = "name_eng"
	} else {
		return nil, fmt.Errorf("lang has to be 'fr' or 'eng'")
	}

	var names []string
	err := db.Model(&ciqual.Food{}).Select(selector).Scan(&names).Error
	if err != nil {
		return nil, err
	}

	return names, nil
}

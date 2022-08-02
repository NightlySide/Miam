package db

import (
	"fmt"

	"github.com/patrickmn/go-cache"
	"io.github.nightlyside/miam/pkg/api"
	"io.github.nightlyside/miam/pkg/ciqual"
)

func (db *Database) GetFoodNames(lang string) ([]string, error) {
	const FUNC_KEY = "get_food_names"

	// select lang
	var selector string
	if lang == "fr" {
		selector = "name_fr"
	} else if lang == "eng" {
		selector = "name_eng"
	} else {
		return nil, fmt.Errorf("lang has to be 'fr' or 'eng'")
	}

	// check if not already in cache
	food_names, found := db.Cache.Get(FUNC_KEY)
	if found {
		db.Cache.Set(FUNC_KEY, food_names, cache.DefaultExpiration)
		return food_names.([]string), nil
	}

	// not in cache, fetch data
	var names []string
	err := db.Model(&ciqual.Food{}).Select(selector).Scan(&names).Error
	if err != nil {
		return nil, err
	}
	db.Cache.Set(FUNC_KEY, names, cache.DefaultExpiration)

	return names, nil
}

func (db *Database) FoodFromName(name string) (*api.Food, error) {
	var FUNC_KEY = "food_from_name_" + name

	// check if not already in cache
	item, found := db.Cache.Get(FUNC_KEY)
	if found {
		db.Cache.Set(FUNC_KEY, item, cache.DefaultExpiration)
		return item.(*api.Food), nil
	}

	// not found, let's fetch the data
	var food ciqual.Food
	err := db.First(&food, "name_fr = ?", name).Error
	if err != nil {
		return nil, err
	}

	var group ciqual.FoodGroup
	err = db.First(&group, "code = ? AND sub_group_code = ? AND sub_sub_group_code = ?", food.GroupCode, food.SubGroupCode, food.SubSubGroupCode).Error
	if err != nil && err.Error() == "record not found" {
		err = db.First(&group, "code = ? AND sub_group_code = ?", food.GroupCode, food.SubGroupCode).Error
		if err.Error() == "record not found" {
			err = db.First(&group, "code = ?", food.GroupCode).Error
			if err != nil {
				return nil, err
			}
		}
	}

	var composition []ciqual.Composition
	err = db.Model(&ciqual.Composition{}).Where("food_code = ?", food.Code).Scan(&composition).Error
	if err != nil {
		return nil, err
	}

	components := map[int]ciqual.Component{}
	for _, compo := range composition {
		var component ciqual.Component
		db.First(&component, "code = ?", compo.ComponentCode)
		components[compo.ComponentCode] = component
	}

	api_food := &api.Food{
		Food:        food,
		Group:       group,
		Composition: composition,
		Components:  components,
	}
	db.Cache.Set(FUNC_KEY, api_food, cache.DefaultExpiration)

	return api_food, nil
}

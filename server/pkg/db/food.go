package db

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v9"
	"io.github.nightlyside/miam/pkg/ciqual"
	"io.github.nightlyside/miam/pkg/models"
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
	food_names_json, err := db.Redis.Get(FUNC_KEY).Bytes()
	if err == nil {
		var food_names []string
		err = json.Unmarshal(food_names_json, &food_names)
		if err != nil {
			return nil, err
		}
		return food_names, nil
	}

	// not in cache, fetch data
	var names []string
	err = db.Model(&ciqual.Food{}).Select(selector).Scan(&names).Error
	if err != nil {
		return nil, err
	}
	names_json, err := json.Marshal(names)
	if err != nil {
		return nil, err
	}
	db.Redis.Set(FUNC_KEY, names_json, redis.KeepTTL)

	return names, nil
}

func (db *Database) FoodFromName(name string) (*models.Food, error) {
	var FUNC_KEY = "food_from_name_" + name

	// check if not already in cache
	item_json, err := db.Redis.Get(FUNC_KEY).Bytes()
	if err == nil {
		var item models.Food
		err = json.Unmarshal(item_json, &item)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}

	// not found, let's fetch the data
	var food ciqual.Food
	err = db.First(&food, "name_fr = ?", name).Error
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

	api_food := &models.Food{
		Food:        food,
		Group:       group,
		Composition: composition,
		Components:  components,
	}

	api_food_json, err := json.Marshal(api_food)
	if err != nil {
		return nil, err
	}
	db.Redis.Set(FUNC_KEY, api_food_json, redis.KeepTTL)

	return api_food, nil
}

func (db *Database) FoodFromCode(code string) (*models.Food, error) {
	var FUNC_KEY = "food_from_code_" + code

	// check if not already in cache
	item_json, err := db.Redis.Get(FUNC_KEY).Bytes()
	if err == nil {
		var item models.Food
		err = json.Unmarshal(item_json, &item)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}

	// not found, let's fetch the data
	var food ciqual.Food
	err = db.First(&food, "code = ?", code).Error
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
	err = db.Model(&ciqual.Composition{}).Where("food_code = ? AND content <> '-'", food.Code).Scan(&composition).Error
	if err != nil {
		return nil, err
	}

	components := map[int]ciqual.Component{}
	for _, compo := range composition {
		var component ciqual.Component
		db.First(&component, "code = ?", compo.ComponentCode)
		components[compo.ComponentCode] = component
	}

	api_food := &models.Food{
		Food:        food,
		Group:       group,
		Composition: composition,
		Components:  components,
	}

	api_food_json, err := json.Marshal(api_food)
	if err != nil {
		return nil, err
	}
	db.Redis.Set(FUNC_KEY, api_food_json, redis.KeepTTL)

	return api_food, nil
}

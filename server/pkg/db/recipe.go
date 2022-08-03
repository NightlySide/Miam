package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v9"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"io.github.nightlyside/miam/pkg/models"
)

type Recipe struct {
	gorm.Model
	Name        string
	Description string
	Yield       string
	Ingredients []Ingredient
	Steps       []Step
	Difficulty  string
	CookTime    string
	PrepTime    string
	Diet        Diet
	Cost        string
	ImageURL    string
}

type Ingredient struct {
	gorm.Model
	RecipeID uint
	FoodCode int
	Quantity float64
	Unit     string
}

type Step struct {
	gorm.Model
	RecipeID        uint
	Text            string
	IngredientsCode pq.Int32Array `gorm:"type:integer[]"`
}

type Diet struct {
	gorm.Model
	RecipeID     uint
	IsVegetarian bool
	IsVegan      bool
	HasGluten    bool
}

func (r *Recipe) Info() string {
	res := ""
	res += fmt.Sprintf("Name: %s\n", r.Name)
	res += fmt.Sprintf("PrepTime: %s\n", r.PrepTime)
	res += fmt.Sprintf("CookTime: %s\n", r.CookTime)
	res += fmt.Sprintf("Yield: %s\n", r.Yield)
	res += fmt.Sprintf("Description: %s\n", r.Description)

	res += "Ingredients:\n"
	for _, ing := range r.Ingredients {
		res += fmt.Sprintf("\t- %.1f %s %d\n", ing.Quantity, ing.Unit, ing.FoodCode)
	}

	res += "Steps:\n"
	for _, step := range r.Steps {
		res += fmt.Sprintf("\t- %s\n", step.Text)
	}

	return res
}

func (db *Database) RecipeFromId(id string) (*models.Recipe, error) {
	var FUNC_KEY = "recipe_from_code_" + id

	// TODO: uncomment cache part
	// check if not already in cache
	// item_json, err := db.Redis.Get(FUNC_KEY).Bytes()
	// if err == nil {
	// 	var item models.Recipe
	// 	err = json.Unmarshal(item_json, &item)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &item, nil
	// }

	// not found, let's fetch the data
	var recipe Recipe
	err := db.Model(&recipe).Preload("Ingredients").Preload("Steps").Preload("Diet").First(&recipe, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// get ingredients
	ingredients := []models.RecipeIngredient{}
	for _, ing := range recipe.Ingredients {
		code := strconv.Itoa(ing.FoodCode)
		food, err := db.FoodFromCode(code)
		if err != nil {
			return nil, err
		}

		ingredients = append(ingredients, models.RecipeIngredient{
			Quantity: ing.Quantity,
			Unit:     ing.Unit,
			Food:     *food,
		})
	}

	// get steps
	steps := []models.RecipeSteps{}
	for _, step := range recipe.Steps {
		ingcodes := []int{}
		for _, code := range step.IngredientsCode {
			ingcodes = append(ingcodes, int(code))
		}

		steps = append(steps, models.RecipeSteps{
			Text:            step.Text,
			IngredientCodes: ingcodes,
		})
	}

	// compute composition
	composition := map[string]models.RecipeComposition{}
	for _, ing := range ingredients {
		for _, ing_compo := range ing.Food.Compositions {
			raw_qty := strings.ReplaceAll(ing_compo.Content, ",", ".")
			if strings.Contains(raw_qty, "<") || raw_qty == "traces" {
				raw_qty = "0"
			}

			qty, err := strconv.ParseFloat(raw_qty, 64)
			if err != nil {
				return nil, err
			}

			// if not registered
			if _, ok := composition[ing_compo.NameFr]; !ok {
				composition[ing_compo.NameFr] = models.RecipeComposition{
					NameFr:   ing_compo.NameFr,
					NameEng:  ing_compo.NameEng,
					Quantity: qty,
					Unit:     ing_compo.NameFr, // TODO: find good unit
				}
			} else {
				compo := composition[ing_compo.NameFr]
				compo.Quantity += qty
				composition[ing_compo.NameFr] = compo
			}
		}
	}

	api_recipe := &models.Recipe{
		Name:        recipe.Name,
		Description: recipe.Description,
		Yield:       recipe.Yield,
		CookTime:    recipe.CookTime,
		PrepTime:    recipe.PrepTime,
		Difficulty:  recipe.Difficulty,
		Cost:        recipe.Cost,
		ImageURL:    recipe.ImageURL,
		Diet: models.RecipeDiet{
			IsVegetarian: recipe.Diet.IsVegetarian,
			IsVegan:      recipe.Diet.IsVegan,
			HasGluten:    recipe.Diet.HasGluten,
		},
		Ingredients: ingredients,
		Steps:       steps,
		Composition: composition,
	}

	api_recipe_json, err := json.Marshal(api_recipe)
	if err != nil {
		return nil, err
	}
	db.Redis.Set(FUNC_KEY, api_recipe_json, redis.KeepTTL)

	return api_recipe, nil
}

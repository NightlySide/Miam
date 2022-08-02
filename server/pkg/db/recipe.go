package db

import (
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
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

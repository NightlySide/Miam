package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
	"io.github.nightlyside/miam/pkg/scraper"
	"io.github.nightlyside/miam/pkg/utils"
)

var ConfigPath = flag.String("config", "", "path to the config file")

func main() {
	// ----- Get the config path
	flag.Parse()
	if *ConfigPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	logrus.SetLevel(logrus.DebugLevel)

	// ----- Get recipe URL

	input := textinput.New("Enter a recipe URL: ")
	input.InitialValue = "https://www.marmiton.org/recettes/recette_pate-a-crepes_12372.aspx"
	input.Placeholder = "https://www.marmiton.org/recettes/recette_pate-a-crepes_12372.aspx"

	url, err := input.RunPrompt()
	if err != nil {
		logrus.Fatal(err)
	}

	// ----- Get infos from URL
	logrus.Debug("Get infos from url")
	scrapedRecipe, err := scraper.Scrape(url)
	if err != nil {
		logrus.Fatal(err)
	}

	// parse ingredients
	ingredients := []*scraper.Ingredient{}
	for _, rawIng := range scrapedRecipe.RecipeIngredient {
		ing, err := scraper.ParseIngredient(rawIng)
		if err != nil {
			logrus.Fatal(err)
		}
		ingredients = append(ingredients, ing)
	}

	// ------ check ingredient with db
	logrus.Debug("Load config")
	cfg, err := config.LoadConfig(*ConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Debug("Connect to db")
	conn, err := db.ConnectDB(cfg)
	if err != nil {
		logrus.Fatalf("connectdb :%w", err)
	}
	logrus.Debug("Get food names list")
	food_names, err := conn.GetFoodNames("fr")
	if err != nil {
		logrus.Fatal(err)
	}

	// for each found ingredient
	logrus.Debug("Process food items")
	foods := []db.Ingredient{}
	for _, ing := range ingredients {
		matches := utils.Autocomplete(ing.Ingredient, food_names)
		if len(matches) == 0 {
			matches = utils.Autocomplete(strings.Split(ing.Ingredient, " ")[0], food_names)
			if len(matches) == 0 {
				logrus.Errorf("no match found for \"%s\"", ing.Ingredient)
				continue
			}
		}

		// matches = matches[:int(math.Min(10., float64(len(matches))))]
		choices := []*selection.Choice{}
		for _, match := range matches {
			choices = append(choices, selection.NewChoice(match))
		}

		sp := selection.New(fmt.Sprintf("Pick an ingredient for: \"%s\"", ing.Ingredient), choices)
		sp.PageSize = 5

		// get the user choice
		choice, err := sp.RunPrompt()
		if err != nil {
			logrus.Fatal(err)
		}

		// get the ingredient from this choice
		food, err := conn.FoodFromName(choice.Value.(string))
		if err != nil {
			logrus.Fatal(err)
		}
		foods = append(foods, db.Ingredient{
			Quantity: ing.Quantity,
			Unit:     ing.Unit,
			FoodCode: food.Code,
		})
	}

	// extract steps
	steps := []db.Step{}
	for _, instr := range scrapedRecipe.RecipeInstructions {
		steps = append(steps, db.Step{
			Text: instr.Text,
			// TODO: missing ingredients
		})
	}

	// ----- build recipe
	recipe := db.Recipe{
		Name:        scrapedRecipe.Name,
		Yield:       scrapedRecipe.RecipeYield,
		Ingredients: foods,
		Steps:       steps,
		Description: scrapedRecipe.Description,
		CookTime:    scrapedRecipe.CookTime,
		PrepTime:    scrapedRecipe.PrepTime,
		ImageURL:    scrapedRecipe.Image[0],
	}

	// ----- Save recipe to db
	err = conn.Create(&recipe).Error
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Recipe: \"%s\" added successfully", recipe.Name)
}

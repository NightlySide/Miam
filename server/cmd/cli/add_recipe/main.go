package main

import (
	"flag"
	"fmt"
	"os"

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

	// ----- Get recipe URL

	input := textinput.New("Enter a recipe URL: ")
	input.InitialValue = "https://www.marmiton.org/recettes/recette_pate-a-crepes_12372.aspx"
	input.Placeholder = "https://www.marmiton.org/recettes/recette_pate-a-crepes_12372.aspx"

	url, err := input.RunPrompt()
	if err != nil {
		logrus.Fatal(err)
	}

	// ----- Get infos from URL
	recipe, err := scraper.Scrape(url)
	if err != nil {
		logrus.Fatal(err)
	}

	// parse ingredients
	ingredients := []*scraper.Ingredient{}
	for _, rawIng := range recipe.RecipeIngredient {
		ing, err := scraper.ParseIngredient(rawIng)
		if err != nil {
			logrus.Fatal(err)
		}
		ingredients = append(ingredients, ing)
	}

	// ------ check ingredient with db
	cfg, err := config.LoadConfig(*ConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}
	conn, err := db.ConnectDB(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	food_names, err := conn.GetFoodNames("fr")
	if err != nil {
		logrus.Fatal(err)
	}
	for _, ing := range ingredients {
		matches := utils.Autocomplete(ing.Ingredient, food_names)
		if len(matches) == 0 {
			logrus.Errorf("no match found for \"%s\"", ing.Ingredient)
			continue
		}

		// matches = matches[:int(math.Min(10., float64(len(matches))))]
		choices := []*selection.Choice{}
		for _, match := range matches {
			choices = append(choices, selection.NewChoice(match))
		}

		sp := selection.New(fmt.Sprintf("Pick an ingredient for: \"%s\"", ing.Ingredient), choices)
		sp.PageSize = 5

		choice, err := sp.RunPrompt()
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info(choice)
	}

}

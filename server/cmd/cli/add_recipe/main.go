package main

import (
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/sirupsen/logrus"
	"io.github.nightlyside/miam/pkg/scraper"
)

func main() {
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
	logrus.Info(recipe.RecipeIngredient)
}

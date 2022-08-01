package api

type Recipe struct {
	Name        string
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
	Food     Food
	Quantity string
}

type Step struct {
	Text        string
	Ingredients []Food
}

type Diet struct {
	IsVegetarian bool
	IsVegan      bool
	HasGluten    bool
}

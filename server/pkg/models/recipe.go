package models

type Recipe struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Yield       string                       `json:"yield"`
	Ingredients []RecipeIngredient           `json:"ingredients"`
	Steps       []RecipeSteps                `json:"steps"`
	Difficulty  string                       `json:"difficulty"`
	Cost        string                       `json:"cost"`
	CookTime    string                       `json:"cook_time"`
	PrepTime    string                       `json:"prep_time"`
	Diet        RecipeDiet                   `json:"diet"`
	ImageURL    string                       `json:"image_url"`
	Composition map[string]RecipeComposition `json:"composition"`
}

type RecipeIngredient struct {
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Food     Food    `json:"food"`
}

type RecipeSteps struct {
	Text            string `json:"text"`
	IngredientCodes []int  `json:"ingredient_codes"`
}

type RecipeDiet struct {
	IsVegetarian bool `json:"is_vegetarian"`
	IsVegan      bool `json:"is_vegan"`
	HasGluten    bool `json:"has_gluten"`
}

type RecipeComposition struct {
	NameFr   string  `json:"name_fr"`
	NameEng  string  `json:"name_eng"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

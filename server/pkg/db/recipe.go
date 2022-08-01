package db

type Recipe struct {
	RecipeId    int
	Name        string
	Yield       string
	Ingredients []Ingredient `gorm:"foreignKey:IngredientId;references:RecipeId"`
	Steps       []Step       `gorm:"foreignKey:StepId;references:RecipeId"`
	Difficulty  string
	CookTime    string
	PrepTime    string
	Diet        Diet `gorm:"foreignKey:DietId;references:RecipeId"`
	Cost        string
	ImageURL    string
}

type Ingredient struct {
	IngredientId int
	Code         int
	Quantity     string
}

type Step struct {
	StepId      int
	Text        string
	Ingredients []Ingredient `gorm:"foreignKey:IngredientId;references:StepId"`
}

type Diet struct {
	DietId       int
	IsVegetarian bool
	IsVegan      bool
	HasGluten    bool
}

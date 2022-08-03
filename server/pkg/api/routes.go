package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) SetRoutes(e *echo.Echo) {
	api := e.Group("/api")
	{
		api.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "hello") })
		api.GET("/import/:url", h.ImportRecipe)
	}

	food := api.Group("/food")
	{
		food.GET("/:code", h.GetFoodFromCode)
		food.GET("/autocomplete/:val", h.AutocompleteFood)
	}

	recipe := api.Group("/recipe")
	{
		recipe.GET("/:id", h.GetRecipeFromId)
	}
}

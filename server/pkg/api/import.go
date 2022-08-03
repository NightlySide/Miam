package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"io.github.nightlyside/miam/pkg/scraper"
)

func (h *Handler) ImportRecipe(c echo.Context) error {
	url := c.Param("url")
	var FUNC_KEY = "import_recipe_" + url

	// check if the site is in the cache
	recipe_json, err := h.DB.Redis.Get(FUNC_KEY).Bytes()
	if err == nil {
		// in the cache
		var recipe scraper.ScrapedRecipe
		err = json.Unmarshal(recipe_json, &recipe)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status: false,
				Error:  err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Status: true,
			Data:   recipe,
		})
	}

	recipe, err := scraper.Scrape(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status: false,
			Error:  err.Error(),
		})
	}

	// set the recipe in the cache
	recipe_json, err = json.Marshal(recipe)
	if err != nil {
		return c.JSON(http.StatusOK, Response{
			Status: true,
			Data:   recipe,
		})
	}
	h.DB.Redis.Set(FUNC_KEY, recipe_json, redis.KeepTTL)

	return c.JSON(http.StatusOK, Response{
		Status: true,
		Data:   recipe,
	})
}

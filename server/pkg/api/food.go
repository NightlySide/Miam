package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"io.github.nightlyside/miam/pkg/utils"
)

func (h *Handler) GetFoodFromCode(c echo.Context) error {
	code := c.Param("code")
	food, err := h.DB.FoodFromCode(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status: false,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, Response{
		Status: true,
		Data:   food,
	})
}

func (h *Handler) AutocompleteFood(c echo.Context) error {
	val := c.Param("val")
	var names []string

	// get names from redis
	names_json, err := h.DB.Redis.Get("food_names").Bytes()
	if err == nil {
		// found something
		err := json.Unmarshal(names_json, &names)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status: false,
				Error:  err.Error(),
			})
		}
	}

	// not in cache, fetching data
	// get the food names
	names, err = h.DB.GetFoodNames("fr")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status: false,
			Error:  err.Error(),
		})
	}
	// save to redis
	names_json, err = json.Marshal(names)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status: false,
			Error:  err.Error(),
		})
	}
	h.DB.Redis.Set("food_names", names_json, redis.KeepTTL)

	// return matches
	matches := utils.Autocomplete(val, names)
	return c.JSON(http.StatusOK, Response{
		Status: true,
		Data:   matches,
	})
}

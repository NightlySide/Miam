package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetRecipeFromId(c echo.Context) error {
	id := c.Param("id")
	food, err := h.DB.RecipeFromId(id)
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

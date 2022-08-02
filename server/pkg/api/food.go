package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

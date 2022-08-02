package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) SetRoutes(e *echo.Echo) {
	api := e.Group("/api")
	{
		api.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "hello") })
	}

	food := api.Group("/food")
	{
		food.GET("/:code", h.GetFoodFromCode)
	}
}

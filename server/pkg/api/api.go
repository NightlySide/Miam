package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io.github.nightlyside/miam/pkg/config"
	"io.github.nightlyside/miam/pkg/db"
)

var version = "0.2"

type Handler struct {
	Config *config.Config
	DB     *db.Database
}

func CreateRouter(handler *Handler) *echo.Echo {
	e := echo.New()

	// Logger middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] - [${method}] ${uri} (${status}) ${error}\n",
	}))

	// Recover middleware
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowCredentials: true,
	}))

	// Default route
	e.GET("/ready", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	e.GET("/meta", func(c echo.Context) error { return c.JSON(http.StatusOK, echo.Map{"version": version}) })

	// Set routes
	handler.SetRoutes(e)

	return e
}

type Response struct {
	Status bool        `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

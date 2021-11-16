package route

import (
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	// TODO: config for e.Debug()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	c := config.Init()
	h := handlers.NewHandler(*c)
	// Routes
	e.GET("/", h.GetHello)

	return e
}

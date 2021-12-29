package route

import (
	"log"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	// TODO: config for e.Debug()
	// Middleware
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Fatal(err)
	}
	h := handlers.NewHandler(*c, d)
	// ROUTES

	e.GET("/", h.GetHello)

	// Formals
	e.GET("/formals", h.GetFormals)
	e.POST("/formals/:id/tickets", h.AddTicket)
	e.DELETE("/formals/:id/tickets", h.CancelTickets)

	// Tickets
	e.GET("/tickets", h.GetTickets)
	e.POST("/tickets", h.BuyTicket)
	e.DELETE("/tickets/:id", h.CancelTicket)
	e.PUT("/tickets/:id", h.EditTicket)

	return e
}

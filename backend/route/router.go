package route

import (
	"log"

	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	// TODO: config for e.Debug()
	// Middleware
	if c.Debug {
		e.Debug = true
		// CORS for testing localhost on different ports
		e.Use(middleware.CORS())
	} else {
		e.Use(middleware.Recover())
	}

	// ROUTES
	a := auth.Init(c)
	h := handlers.NewHandler(*c, d, a)
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

	// Auth
	e.GET("/auth/redirect", h.AuthRedirect)
	e.GET("/auth/callback", h.AuthCallback)

	return e
}

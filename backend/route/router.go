package route

import (
	"log"

	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/handlers"
	"github.com/kcsu/store/middleware"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(em.Logger())

	// TODO: config for e.Debug()
	// Middleware
	if c.Debug {
		e.Debug = true
		// CORS for testing localhost on different ports
		e.Use(em.CORS())
	} else {
		e.Use(em.Recover())
	}

	// ROUTES
	a := auth.Init(c)
	h := handlers.NewHandler(*c, d, a)
	e.GET("/", h.GetHello)

	// Auth
	e.GET("/auth/redirect", h.AuthRedirect)
	e.GET("/auth/callback", h.AuthCallback)

	requireAuth := middleware.JWTAuth(c)

	formals := e.Group("/formals", requireAuth)
	// Formals
	formals.GET("", h.GetFormals)
	formals.POST("/:id/tickets", h.AddTicket)
	formals.DELETE("/:id/tickets", h.CancelTickets)

	tickets := e.Group("/tickets", requireAuth)
	// Tickets
	tickets.GET("", h.GetTickets)
	tickets.POST("", h.BuyTicket)
	tickets.DELETE("/:id", h.CancelTicket)
	tickets.PUT("/:id", h.EditTicket)

	return e
}

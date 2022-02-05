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

// Initialise the router
func Init() *echo.Echo {
	// Load config from the environment
	c := config.Init()

	// Initialise database connection
	d, err := db.Init(c)
	if err != nil {
		log.Panic(err)
	}

	// Create the router
	e := echo.New()

	// Configure global middleware
	e.Use(em.Logger())
	if c.Debug {
		e.Debug = true
		// CORS for testing localhost on different ports
		e.Use(em.CORSWithConfig(em.CORSConfig{
			AllowCredentials: true,
		}))
	} else {
		e.HideBanner = true
		e.Use(em.Recover())
	}

	api := e.Group("/api")

	// ROUTE HANDLERS
	// Create the handler object which stores useful data and methods
	store := auth.InitSessionStore(c)
	a := auth.Init(c, store)
	h := handlers.NewHandler(*c, d, a)

	// Authentication middleware
	requireAuth := middleware.JWTAuth(c)

	// TODO: remove
	api.GET("/", h.GetHello)

	// Auth Routes
	api.GET("/auth/user", h.GetUser, requireAuth)
	api.POST("/auth/logout", h.Logout, requireAuth)

	formals := api.Group("/formals", requireAuth)
	// Formal routes
	formals.GET("", h.GetFormals)
	formals.POST("/:id/tickets", h.AddTicket)
	formals.DELETE("/:id/tickets", h.CancelTickets)

	tickets := api.Group("/tickets", requireAuth)
	// Ticket routes
	tickets.GET("", h.GetTickets)
	tickets.POST("", h.BuyTicket)
	tickets.DELETE("/:id", h.CancelTicket)
	tickets.PUT("/:id", h.EditTicket)

	return e
}

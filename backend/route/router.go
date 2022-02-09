package route

import (
	"log"

	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/handlers"
	"github.com/kcsu/store/middleware"
	"github.com/labstack/echo-contrib/session"
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
	e.Validator = middleware.NewValidator()

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

	// ROUTE HANDLERS
	// Create the handler object which stores useful data and methods
	store := auth.InitSessionStore(c)
	e.Use(session.Middleware(store))
	a := auth.Init(c, store)
	h := handlers.NewHandler(*c, d, a)

	api := e.Group("/api")

	// Authentication middleware
	requireAuth := middleware.JWTAuth(c)

	// TODO: change to health
	api.GET("/", h.GetHello)

	// Auth Routes
	auth := api.Group("/oauth")
	auth.GET("/redirect", h.AuthRedirect)
	auth.GET("/callback", h.AuthCallback)
	auth.GET("/user", h.GetUser, requireAuth)
	auth.POST("/logout", h.Logout, requireAuth)

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

package route

import (
	"log"

	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/handlers"
	"github.com/kcsu/store/handlers/admin"
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

	// Authentication middleware
	requireAuth := middleware.JWTAuth(c)
	api := e.Group("/api")
	ApiRoutes(api, h, requireAuth)

	// TODO: Protect!
	adminApi := api.Group("/admin", requireAuth)
	ah := admin.NewHandler(h, d)
	AdminRoutes(adminApi, ah)

	return e
}

func ApiRoutes(api *echo.Group, h *handlers.Handler, requireAuth echo.MiddlewareFunc) {

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
}

func AdminRoutes(a *echo.Group, ah *admin.AdminHandler) {
	// FIXME: CRUD or R/W/D?
	rbac := middleware.NewRBAC(middleware.RbacConfig{
		Auth:  ah.Auth,
		Users: ah.Users,
	})

	formals := a.Group("/formals")
	formals.GET("", ah.GetFormals, rbac.M("formals", "read"))
	formals.POST("", ah.CreateFormal, rbac.M("formals", "write"))
	formals.GET("/:id", ah.GetFormal, rbac.M("formals", "read"))
	formals.PUT("/:id", ah.UpdateFormal, rbac.M("formals", "write"))
	formals.PUT("/:id/groups", ah.UpdateFormalGroups, rbac.M("formals", "write"))
	formals.DELETE("/:id", ah.DeleteFormal, rbac.M("formals", "delete"))

	tickets := a.Group("/tickets")
	tickets.POST("/manual", ah.CreateManualTicket, rbac.M("tickets", "write"))
	tickets.DELETE("/manual/:id", ah.CancelManualTicket, rbac.M("tickets", "delete"))
	tickets.PUT("/manual/:id", ah.EditManualTicket, rbac.M("tickets", "write"))
	tickets.DELETE("/:id", ah.CancelTicket, rbac.M("tickets", "delete"))
	tickets.PUT("/:id", ah.EditTicket, rbac.M("tickets", "write"))

	groups := a.Group("/groups")
	groups.GET("", ah.GetGroups, rbac.M("groups", "read"))
	groups.POST("", ah.CreateGroup, rbac.M("groups", "write"))
	groups.GET("/:id", ah.GetGroup, rbac.M("groups", "read"))
	groups.PUT("/:id", ah.UpdateGroup, rbac.M("groups", "write"))
	groups.DELETE("/:id", ah.DeleteGroup, rbac.M("groups", "delete"))
	groups.POST("/:id/users", ah.AddGroupUser, rbac.M("groups", "write"))
	groups.DELETE("/:id/users", ah.RemoveGroupUser, rbac.M("groups", "write"))
	groups.POST("/:id/users/lookup", ah.LookupGroupUsers, rbac.M("groups", "write"))

	roles := a.Group("/roles")
	roles.GET("", ah.GetRoles, rbac.M("roles", "read"))
	roles.POST("", ah.CreateRole, rbac.M("roles", "write"))
	roles.PUT("/:id", ah.UpdateRole, rbac.M("roles", "write"))
	roles.GET("/users", ah.GetUserRoles, rbac.M("roles", "read"))
	roles.POST("/users", ah.AddUserRole, rbac.M("roles", "write"))
	roles.DELETE("/users", ah.RemoveUserRole, rbac.M("roles", "write"))
	roles.DELETE("/:id", ah.DeleteRole, rbac.M("roles", "delete"))

	permissions := a.Group("/permissions")
	permissions.POST("", ah.CreatePermission, rbac.M("permissions", "write"))
	permissions.DELETE("/:id", ah.DeletePermission, rbac.M("permissions", "delete"))

	bills := a.Group("/bills")
	bills.GET("", ah.GetBills, rbac.M("billing", "read"))
	bills.GET("/:id", ah.GetBill, rbac.M("billing", "read"))
	bills.GET("/:id/stats", ah.GetBillStats, rbac.M("billing", "read"))
	bills.GET("/:id/stats/formals.csv", ah.GetBillFormalStatsCSV, rbac.M("billing", "read"))
	bills.GET("/:id/stats/users.csv", ah.GetBillUserStatsCSV, rbac.M("billing", "read"))
	bills.PUT("/:id", ah.UpdateBill, rbac.M("billing", "write"))
	bills.POST("/:id/formals", ah.AddBillFormals, rbac.M("billing", "write"))
	bills.DELETE("/:id/formals/:formalId", ah.RemoveBillFormal, rbac.M("billing", "write"))
}

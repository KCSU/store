package admin

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/handlers"
)

// Helper struct containing useful data and methods for admin handlers to use
type AdminHandler struct {
	Config  config.Config
	Formals db.FormalStore
	Tickets db.TicketStore
	Users   db.UserStore
	Auth    auth.Auth
}

// Initialise the handler helper
func NewHandler(h *handlers.Handler) *AdminHandler {
	return &AdminHandler{
		h.Config,
		h.Formals,
		h.Tickets,
		h.Users,
		h.Auth,
	}
}

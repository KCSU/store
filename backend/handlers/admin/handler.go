package admin

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/handlers"
	"github.com/kcsu/store/lookup"
	"gorm.io/gorm"
)

// Helper struct containing useful data and methods for admin handlers to use
type AdminHandler struct {
	Config  config.Config
	Formals db.FormalStore
	Tickets db.TicketStore
	Users   db.UserStore
	Groups  db.GroupStore
	Auth    auth.Auth
	Lookup  lookup.Lookup
}

// Initialise the handler helper
func NewHandler(h *handlers.Handler, d *gorm.DB) *AdminHandler {
	groups := db.NewGroupStore(d)
	lookup := lookup.New(h.Config.LookupApiUrl, groups)
	return &AdminHandler{
		h.Config,
		h.Formals,
		h.Tickets,
		h.Users,
		groups,
		h.Auth,
		lookup,
	}
}

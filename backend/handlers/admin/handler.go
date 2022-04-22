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
	Config        config.Config
	Formals       db.FormalStore
	Tickets       db.TicketStore
	ManualTickets db.ManualTicketStore
	Users         db.UserStore
	Groups        db.GroupStore
	Roles         db.RoleStore
	Bills         db.BillStore
	Auth          auth.Auth
	Lookup        lookup.Lookup
}

// Initialise the handler helper
func NewHandler(h *handlers.Handler, d *gorm.DB) *AdminHandler {
	groups := db.NewGroupStore(d)
	roles := db.NewRoleStore(d)
	bills := db.NewBillStore(d)
	lookup := lookup.New(h.Config.LookupApiUrl, groups)
	manualTickets := db.NewManualTicketStore(d)
	return &AdminHandler{
		h.Config,
		h.Formals,
		h.Tickets,
		manualTickets,
		h.Users,
		groups,
		roles,
		bills,
		h.Auth,
		lookup,
	}
}

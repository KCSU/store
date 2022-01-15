package handlers

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"gorm.io/gorm"
)

// Helper struct containing useful data and methods for handlers to use
type Handler struct {
	config  config.Config
	formals db.FormalStore
	tickets db.TicketStore
	users   db.UserStore
	auth    *auth.Auth
}

// Initialise the handler helper
func NewHandler(c config.Config, d *gorm.DB, a *auth.Auth) *Handler {
	return &Handler{
		config:  c,
		formals: db.NewFormalStore(d),
		tickets: db.NewTicketStore(d),
		users:   db.NewUserStore(d),
		auth:    a,
	}
}

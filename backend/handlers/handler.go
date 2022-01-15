package handlers

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"gorm.io/gorm"
)

// Helper struct containing useful data and methods for handlers to use
type Handler struct {
	Config  config.Config
	Formals db.FormalStore
	Tickets db.TicketStore
	Users   db.UserStore
	Auth    *auth.Auth
}

// Initialise the handler helper
func NewHandler(c config.Config, d *gorm.DB, a *auth.Auth) *Handler {
	return &Handler{
		Config:  c,
		Formals: db.NewFormalStore(d),
		Tickets: db.NewTicketStore(d),
		Users:   db.NewUserStore(d),
		Auth:    a,
	}
}

package handlers

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"gorm.io/gorm"
)

type Handler struct {
	config  config.Config
	formals *db.FormalStore
	tickets *db.TicketStore
	auth    *auth.Auth
}

func NewHandler(c config.Config, d *gorm.DB, a *auth.Auth) *Handler {
	return &Handler{
		config:  c,
		formals: db.NewFormalStore(d),
		tickets: db.NewTicketStore(d),
		auth:    a,
	}
}

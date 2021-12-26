package handlers

import (
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"gorm.io/gorm"
)

type Handler struct {
	config  config.Config
	db      *gorm.DB
	formals *db.FormalStore
}

func NewHandler(c config.Config, d *gorm.DB) *Handler {
	return &Handler{
		config:  c,
		db:      d,
		formals: db.NewFormalStore(d),
	}
}

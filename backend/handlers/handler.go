package handlers

import (
	"github.com/kcsu/store/config"
	"gorm.io/gorm"
)

type Handler struct {
	config config.Config
	db     *gorm.DB
}

func NewHandler(c config.Config, d *gorm.DB) *Handler {
	return &Handler{
		config: c,
		db:     d,
	}
}

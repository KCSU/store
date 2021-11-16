package handlers

import "github.com/kcsu/store/config"

type Handler struct {
	config config.Config
}

func NewHandler(c config.Config) *Handler {
	return &Handler{
		config: c,
	}
}

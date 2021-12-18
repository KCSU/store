package handlers

import (
	"net/http"

	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetTickets(c echo.Context) error {
	var tickets []model.Ticket
	h.db.Find(&tickets)
	return c.JSON(http.StatusOK, &tickets)
}

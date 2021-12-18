package handlers

import (
	"net/http"

	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetFormals(c echo.Context) error {
	var formals []model.Formal
	// TODO: handle errors
	h.db.Find(&formals)
	formalData := make([]dto.FormalDto, len(formals))
	for i, f := range formals {
		formalData[i].Formal = f
		formalData[i].TicketsRemaining = f.Tickets - uint(
			h.db.Model(f).Not("is_queue OR is_guest").Association("TicketSales").Count(),
		)
		formalData[i].GuestTicketsRemaining = f.GuestTickets - uint(
			h.db.Model(f).Not("is_queue").Where("is_guest").Association("TicketSales").Count(),
		)
	}
	return c.JSON(http.StatusOK, &formalData)
}

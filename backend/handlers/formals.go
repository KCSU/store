package handlers

import (
	"net/http"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetFormals(c echo.Context) error {
	// TODO: handle errors
	// TODO: filter by time
	formals, err := h.formals.Get()
	if err != nil {
		return echo.ErrInternalServerError
	}
	formalData := make([]dto.FormalDto, len(formals))
	for i, f := range formals {
		formalData[i].Formal = f
		formalData[i].TicketsRemaining = h.formals.TicketsRemaining(&f, false)
		formalData[i].GuestTicketsRemaining = h.formals.TicketsRemaining(&f, true)
	}
	return c.JSON(http.StatusOK, &formalData)
}

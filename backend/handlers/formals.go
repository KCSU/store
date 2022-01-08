package handlers

import (
	"net/http"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

// Handler to fetch a list of upcoming formals
func (h *Handler) GetFormals(c echo.Context) error {
	// TODO: handle errors
	formals, err := h.formals.Get()
	if err != nil {
		return err
	}
	// Create JSON response
	formalData := make([]dto.FormalDto, len(formals))
	for i, f := range formals {
		formalData[i].Formal = f
		formalData[i].TicketsRemaining = h.formals.TicketsRemaining(&f, false)
		formalData[i].GuestTicketsRemaining = h.formals.TicketsRemaining(&f, true)
	}
	return c.JSON(http.StatusOK, &formalData)
}

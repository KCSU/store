package handlers

import (
	"net/http"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

// Handler to fetch a list of upcoming formals
func (h *Handler) GetFormals(c echo.Context) error {
	// TODO: handle errors
	userId := h.Auth.GetUserId(c)
	formals, err := h.Formals.GetWithUserData(userId)
	if err != nil {
		return err
	}
	// Create JSON response
	// TODO: DRY!!!!!
	formalData := make([]dto.FormalDto, len(formals))
	for i, f := range formals {
		formalData[i].Formal = f
		// FIXME: This is horribly inefficient!!
		formalData[i].TicketsRemaining = h.Formals.TicketsRemaining(&f, false)
		formalData[i].GuestTicketsRemaining = h.Formals.TicketsRemaining(&f, true)
		formalData[i].MyTickets = f.TicketSales
		groups := make([]dto.GroupDto, len(f.Groups))
		for j, g := range f.Groups {
			groups[j] = dto.GroupDto{
				ID:   g.ID,
				Name: g.Name,
			}
		}
		formalData[i].Groups = groups
	}
	return c.JSON(http.StatusOK, &formalData)
}

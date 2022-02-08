package handlers

import (
	"net/http"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

// Handler to fetch a list of upcoming formals
func (h *Handler) GetFormals(c echo.Context) error {
	// TODO: handle errors
	formals, err := h.Formals.GetWithGroups()
	if err != nil {
		return err
	}
	// Create JSON response
	formalData := make([]dto.FormalDto, len(formals))
	for i, f := range formals {
		formalData[i].Formal = f
		formalData[i].TicketsRemaining = h.Formals.TicketsRemaining(&f, false)
		formalData[i].GuestTicketsRemaining = h.Formals.TicketsRemaining(&f, true)
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

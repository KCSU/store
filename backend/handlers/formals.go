package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

// Handler to fetch a guest list for a formal
func (h *Handler) GetFormalGuestList(c echo.Context) error {
	id := c.Param("id")
	formalId, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	formal, err := h.Formals.Find(formalId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if !formal.HasGuestList || !formal.IsVisible {
		return echo.ErrForbidden
	}
	guests, err := h.Formals.FindGuestList(formalId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &guests)
}

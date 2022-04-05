package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Cancel a given ticket
//
// If the ticket is not a guest ticket, also cancel
// all the guest tickets for the same formal and user
func (ah *AdminHandler) CancelTicket(c echo.Context) error {
	id := c.Param("id")
	ticketID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}
	ticket, err := ah.Tickets.Find(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}

	if !ticket.IsGuest {
		// Delete all guest tickets for the same formal and user
		if err := ah.Tickets.DeleteByFormal(ticket.FormalID, ticket.UserID); err != nil {
			return err
		}
	} else {
		// Delete the ticket
		if err := ah.Tickets.Delete(ticketID); err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusOK)
}

package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
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

// Update a specified ticket
func (ah *AdminHandler) EditTicket(c echo.Context) error {
	id := c.Param("id")
	ticketID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}

	t := new(dto.TicketRequestDto)
	if err := c.Bind(t); err != nil {
		return err
	}
	if err := c.Validate(t); err != nil {
		return err
	}
	if err := ah.Tickets.Update(ticketID, t); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	return c.NoContent(http.StatusOK)
}

// Create a manual ticket
func (ah *AdminHandler) CreateManualTicket(c echo.Context) error {
	t := new(dto.ManualTicketDto)
	if err := c.Bind(t); err != nil {
		return err
	}
	if err := c.Validate(t); err != nil {
		return err
	}
	ticket := model.ManualTicket{
		MealOption:    t.MealOption,
		FormalID:      t.FormalID,
		Type:          t.Type,
		Name:          t.Name,
		Justification: t.Justification,
		Email:         t.Email,
	}
	// TODO: check if the formal exists
	if err := ah.ManualTickets.Create(&ticket); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

// Delete a manual ticket
func (ah *AdminHandler) DeleteManualTicket(c echo.Context) error {
	id := c.Param("id")
	ticketID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}
	if _, err := ah.ManualTickets.Find(ticketID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.ManualTickets.Delete(ticketID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

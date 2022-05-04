package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
	ticketID, err := uuid.Parse(id)
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
		if err := ah.Access.Log(c,
			fmt.Sprintf(
				"cancelled tickets for user %q, formal %q",
				ticket.User.Name, ticket.Formal.Name,
			),
			map[string]string{
				"formalId":  ticket.FormalID.String(),
				"userEmail": ticket.User.Email,
				"ticketId":  ticket.ID.String(),
			},
		); err != nil {
			return err
		}
	} else {
		// Delete the ticket
		if err := ah.Tickets.Delete(ticketID); err != nil {
			return err
		}
		if err := ah.Access.Log(c,
			fmt.Sprintf(
				"cancelled a guest ticket for user %q, formal %q",
				ticket.User.Name, ticket.Formal.Name,
			),
			map[string]string{
				"formalId":  ticket.FormalID.String(),
				"userEmail": ticket.User.Email,
				"ticketId":  ticket.ID.String(),
			},
		); err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusOK)
}

// Update a specified ticket
func (ah *AdminHandler) EditTicket(c echo.Context) error {
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
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
	// TODO: fetch named metadata!!
	if err := ah.Access.Log(c, "updated a ticket", map[string]string{
		"ticketId": ticketID.String(),
	}); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// Create a manual ticket
func (ah *AdminHandler) CreateManualTicket(c echo.Context) error {
	t := new(dto.CreateManualTicketDto)
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
	f, err := ah.Formals.Find(ticket.FormalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound // FIXME: this seems bad
		}
		return err
	}
	if err = ah.ManualTickets.Create(&ticket); err != nil {
		return err
	}
	if err := ah.Access.Log(c,
		fmt.Sprintf(
			"created a manual ticket for user %q, formal %q",
			ticket.Name, f.Name,
		),
		map[string]string{
			"formalId":      ticket.FormalID.String(),
			"userEmail":     ticket.Email,
			"justification": ticket.Justification,
			"ticketId":      ticket.ID.String(),
		},
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

// Delete a manual ticket
func (ah *AdminHandler) CancelManualTicket(c echo.Context) error {
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	ticket, err := ah.ManualTickets.Find(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.ManualTickets.Delete(ticketID); err != nil {
		return err
	}
	if err := ah.Access.Log(c,
		fmt.Sprintf(
			"cancelled a manual ticket for %q",
			ticket.Name, // XXX: Formal name?
		),
		map[string]string{
			"ticketId":  ticketID.String(),
			"formalId":  ticket.FormalID.String(),
			"userEmail": ticket.Email,
		},
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (ah *AdminHandler) EditManualTicket(c echo.Context) error {
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	t := new(dto.EditManualTicketDto)
	if err := c.Bind(t); err != nil {
		return err
	}
	if err := c.Validate(t); err != nil {
		return err
	}
	ticket, err := ah.ManualTickets.Find(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	ticket.MealOption = t.MealOption
	ticket.Type = t.Type
	ticket.Name = t.Name
	ticket.Justification = t.Justification
	ticket.Email = t.Email
	if err := ah.ManualTickets.Update(&ticket); err != nil {
		return err
	}
	if err := ah.Access.Log(c,
		fmt.Sprintf(
			"updated a manual ticket for %q",
			ticket.Name, // XXX: Formal name?
		),
		map[string]string{
			"ticketId":  ticketID.String(),
			"formalId":  ticket.FormalID.String(),
			"userEmail": ticket.Email,
		},
	); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

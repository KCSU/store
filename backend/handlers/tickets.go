package handlers

import (
	"net/http"

	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetTickets(c echo.Context) error {
	// Load tickets from the database
	var tickets []model.Ticket
	h.db.Preload("Formal").Find(&tickets)
	// Group tickets by formal and guest status
	dtos := map[uint]dto.TicketDto{}
	for _, t := range tickets {
		if myDto, hasDto := dtos[t.Formal.ID]; hasDto {
			if t.IsGuest {
				myDto.GuestTickets = append(myDto.GuestTickets, t)
			} else {
				myDto.Ticket = t
			}
			dtos[t.Formal.ID] = myDto
		} else {
			myDto := dto.TicketDto{
				Formal: *t.Formal,
			}
			if t.IsGuest {
				myDto.GuestTickets = []model.Ticket{t}
			} else {
				myDto.GuestTickets = []model.Ticket{}
				myDto.Ticket = t
			}
			dtos[t.Formal.ID] = myDto
		}

	}
	// Convert to slice and return
	v := make([]dto.TicketDto, 0, len(dtos))
	for _, val := range dtos {
		v = append(v, val)
	}
	return c.JSON(http.StatusOK, v)
}

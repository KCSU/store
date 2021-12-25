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
	dtoList := make([]dto.TicketDto, 0, len(dtos))
	for _, val := range dtos {
		dtoList = append(dtoList, val)
	}
	return c.JSON(http.StatusOK, dtoList)
}

func (h *Handler) BuyTicket(c echo.Context) error {
	// Bind request data
	t := new(dto.BuyTicketDto)
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO:
	// if err := c.Validate(t); err != nil {
	// 	return err
	// }

	// Check that ticket does not already exist
	// TODO: users
	var ticketExists int64
	h.db.Model(&model.Ticket{}).Not("is_guest").Where("formal_id = ?", t.FormalId).Count(&ticketExists)
	if ticketExists > 0 {
		return echo.NewHTTPError(http.StatusConflict, "Ticket already exists.")
	}

	// Instantiate DB objects to be inserted
	tickets := make([]model.Ticket, len(t.GuestTickets)+1)
	tickets[0] = model.Ticket{
		FormalID:   int(t.FormalId),
		IsGuest:    false,
		IsQueue:    true,
		MealOption: t.Ticket.MealOption,
	}
	for i, gt := range t.GuestTickets {
		tickets[i+1] = model.Ticket{
			FormalID:   int(t.FormalId),
			IsGuest:    true,
			IsQueue:    true,
			MealOption: gt.MealOption,
		}
	}
	// Insert into DB
	if err := h.db.Create(tickets).Error; err != nil {
		return err
	}
	// TODO: h.db.Clauses(clause.OnConflict{DoNothing: true})
	// Should this return some data?
	return c.NoContent(http.StatusOK)
}

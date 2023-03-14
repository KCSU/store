package handlers

import (
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// TODO: implement 404 checks everywhere they are needed

// Get a list of the user's tickets, grouped by formal
func (h *Handler) GetTickets(c echo.Context) error {
	userId := h.Auth.GetUserId(c)
	// Load tickets from the database
	tickets, err := h.Tickets.Get(userId)
	if err != nil {
		return err
	}
	// Group tickets by formal and guest status
	dtos := map[uuid.UUID]dto.TicketDto{}
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
	// Sort by datetime
	sort.Slice(dtoList, func(i, j int) bool {
		return dtoList[i].Formal.DateTime.Before(dtoList[j].Formal.DateTime)
	})
	return c.JSON(http.StatusOK, dtoList)
}

// Buy tickets for a formal, potentially with guest tickets
func (h *Handler) BuyTicket(c echo.Context) error {
	// Get the logged-in user
	userId := h.Auth.GetUserId(c)
	user, err := h.Users.Find(userId)
	if err != nil {
		return echo.ErrUnauthorized
	}
	// Bind request data
	t := new(dto.BuyTicketDto)
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(t); err != nil {
		return err
	}

	formal, err := h.Formals.FindWithGroups(t.FormalId)
	if err != nil {
		// Check whether error is formal existence?
		return err
	}

	// Check that the user belongs to the requisite group
	canBuy, err := h.canBuyTickets(&user, &formal)
	if err != nil {
		return err
	}
	if !canBuy {
		return echo.ErrForbidden
	}

	// Check that tickets are still on sale
	if time.Now().After(formal.SaleEnd) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Sales have closed.")
	}

	// Check that formal permits this many tickets
	if len(t.GuestTickets) > int(formal.GuestLimit) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Too many guest tickets requested.")
	}

	// Check that ticket does not already exist
	// TODO: users
	ticketExists, err := h.Tickets.ExistsByFormal(t.FormalId, userId)
	if err != nil {
		return err
	}
	if ticketExists {
		return echo.NewHTTPError(http.StatusConflict, "Ticket already exists.")
	}

	// Instantiate DB objects to be inserted
	tickets := make([]model.Ticket, len(t.GuestTickets)+1)
	tickets[0] = model.Ticket{
		FormalID:   t.FormalId,
		IsGuest:    false,
		IsQueue:    true,
		MealOption: t.Ticket.MealOption,
		UserID:     userId,
	}
	for i, gt := range t.GuestTickets {
		tickets[i+1] = model.Ticket{
			FormalID:   t.FormalId,
			IsGuest:    true,
			IsQueue:    true,
			MealOption: gt.MealOption,
			UserID:     userId,
		}
	}
	// Insert into DB
	if err := h.Tickets.BatchCreate(tickets); err != nil {
		return err
	}
	// TODO: h.db.Clauses(clause.OnConflict{DoNothing: true})
	// Should this return some data?
	return c.NoContent(http.StatusCreated)
}

// Cancel the user's tickets for a given formal
func (h *Handler) CancelTickets(c echo.Context) error {
	userId := h.Auth.GetUserId(c)

	// Get the formal ID from query
	id := c.Param("id")
	formalID, err := uuid.Parse(id)
	if err != nil {
		// TODO: NewHTTPError?
		return echo.ErrNotFound
	}

	formal, err := h.Formals.Find(formalID)
	if err != nil {
		// Check whether error is formal existence?
		return err
	}

	// Check that tickets are still on sale
	if time.Now().After(formal.SaleEnd) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Sales have closed.")
	}

	// [Soft-]delete ticket from the database
	if err := h.Tickets.DeleteByFormal(formalID, userId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	return c.NoContent(http.StatusOK)
}

// Cancel a specified ticket
func (h *Handler) CancelTicket(c echo.Context) error {
	userId := h.Auth.GetUserId(c)
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		// Should this be a different error?
		return echo.ErrNotFound
	}
	// Fetch ticket from the database
	ticket, err := h.Tickets.FindWithFormal(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}

	// Cannot cancel after sales have ended
	if time.Now().After(ticket.Formal.SaleEnd) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Sales have closed.")
	}

	// Check the ticket really belongs to the logged-in user
	if ticket.UserID != userId {
		return echo.ErrForbidden
	}
	// The ticket must be a guest ticket
	if !ticket.IsGuest {
		return echo.NewHTTPError(http.StatusForbidden, "Non-guest tickets must be cancelled as a group")
	}
	// Delete the ticket from the database
	if err := h.Tickets.Delete(ticketID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// Update a specified ticket
func (h *Handler) EditTicket(c echo.Context) error {
	userId := h.Auth.GetUserId(c)
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	// Fetch ticket from the database
	ticket, err := h.Tickets.FindWithFormal(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}

	// Cannot edit after sales have ended
	if time.Now().After(ticket.Formal.SaleEnd) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Sales have closed.")
	}

	// Check the ticket really belongs to the logged-in user
	if ticket.UserID != userId {
		return echo.ErrForbidden
	}
	// Update the ticket based on request data
	t := new(dto.TicketRequestDto)
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(t); err != nil {
		return err
	}
	if err := h.Tickets.Update(ticketID, t); err != nil {
		return err
	}
	// TODO: return the new model
	return c.NoContent(http.StatusOK)
}

// Add a guest ticket to the user's tickets for a specified formal
//
// TODO: reduce DB calls
func (h *Handler) AddTicket(c echo.Context) error {
	userId := h.Auth.GetUserId(c)
	// Parse request data
	id := c.Param("id")
	formalID, err := uuid.Parse(id)
	if err != nil {
		// TODO: NewHTTPError?
		return echo.ErrNotFound
	}
	t := new(dto.TicketRequestDto)
	if err := c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(t); err != nil {
		return err
	}

	// Make sure the user already has a ticket to this formal
	exists, err := h.Tickets.ExistsByFormal(formalID, userId)
	if err != nil {
		return err
	}
	if !exists {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "You must purchase a normal ticket first.")
	}

	// Make sure the user doesn't have too many guest tickets already
	count, err := h.Tickets.CountGuestByFormal(formalID, userId)
	if err != nil {
		return err
	}
	formal, err := h.Formals.Find(formalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}

	// Check that tickets are still on sale
	if time.Now().After(formal.SaleEnd) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Sales have closed.")
	}

	if int64(formal.GuestLimit)-count <= 0 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Too many guest tickets requested.")
	}

	// Create the ticket in the DB
	ticket := model.Ticket{
		FormalID:   formalID,
		IsGuest:    true,
		IsQueue:    true,
		MealOption: t.MealOption,
		UserID:     userId,
	}
	if err := h.Tickets.Create(&ticket); err != nil {
		return err
	}
	// TODO: return the new model
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) ScanTicket(c echo.Context) error {
	// Parse request data
	id := c.Param("id")
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	ticket, err := h.Tickets.Find(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if ticket.IsQueue {
		return echo.NewHTTPError(
			http.StatusUnprocessableEntity,
			"Cannot scan a queue ticket",
		)
	}
	// Update the ticket in the DB
	if err := h.Tickets.Scan(ticketID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.ScannedTicketDto{
		ID:         ticket.ID,
		IsGuest:    ticket.IsGuest,
		IsScanned:  ticket.IsScanned,
		MealOption: ticket.MealOption,
		FormalID:   ticket.FormalID,
		FormalName: ticket.Formal.Name,
		FormalDate: ticket.Formal.DateTime,
		UserName:   ticket.User.Name,
	})
}

// Check if the specified user can buy tickets to the specified formal
func (h *Handler) canBuyTickets(user *model.User, formal *model.Formal) (bool, error) {
	userGroups, err := h.Users.Groups(user)
	if err != nil {
		return false, err
	}
	s := map[uuid.UUID]bool{}
	for _, g := range userGroups {
		s[g.ID] = true
	}
	for _, g := range formal.Groups {
		if _, ok := s[g.ID]; ok {
			return true, nil
		}
	}
	return false, nil
}

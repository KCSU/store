package dto

import "github.com/kcsu/store/model"

type TicketDto struct {
	// TODO: maybe use pointers??
	Formal       model.Formal   `json:"formal"`
	Ticket       model.Ticket   `json:"ticket"`
	GuestTickets []model.Ticket `json:"guestTickets"`
}

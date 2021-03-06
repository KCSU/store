package dto

import "github.com/kcsu/store/model"

type FormalDto struct {
	model.Formal
	TicketsRemaining      uint           `json:"ticketsRemaining"`
	GuestTicketsRemaining uint           `json:"guestTicketsRemaining"`
	Groups                []GroupDto     `json:"groups"`
	MyTickets             []model.Ticket `json:"myTickets,omitempty"`
}

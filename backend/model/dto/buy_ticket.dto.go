package dto

import "github.com/google/uuid"

type BuyTicketDto struct {
	FormalId     uuid.UUID          `json:"formalId"`
	Ticket       TicketRequestDto   `json:"ticket"`
	GuestTickets []TicketRequestDto `json:"guestTickets"`
}

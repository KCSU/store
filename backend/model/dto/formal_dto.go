package dto

import "github.com/kcsu/store/model"

type FormalDto struct {
	model.Formal
	TicketsRemaining      uint `json:"ticketsRemaining"`
	GuestTicketsRemaining uint `json:"guestTicketsRemaining"`
}

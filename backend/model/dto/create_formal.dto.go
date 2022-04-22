package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kcsu/store/model"
)

type CreateFormalDto struct {
	Name         string      `json:"name" validate:"required,min=5,max=100"`
	Menu         string      `json:"menu"` // TODO: validate menu?
	Price        float32     `json:"price" validate:"min=0"`
	GuestPrice   float32     `json:"guestPrice" validate:"min=0"`
	GuestLimit   uint        `json:"guestLimit"`
	Tickets      uint        `json:"tickets"`
	GuestTickets uint        `json:"guestTickets"`
	SaleStart    time.Time   `json:"saleStart" validate:"required"`
	SaleEnd      time.Time   `json:"saleEnd" validate:"required"`  // TODO: gt SaleStart?
	DateTime     time.Time   `json:"dateTime" validate:"required"` // TODO: gt SaleEnd?
	Groups       []uuid.UUID `json:"groups"`
}

func (f *CreateFormalDto) Formal() model.Formal {
	return model.Formal{
		Name:         f.Name,
		Menu:         f.Menu,
		Price:        f.Price,
		GuestPrice:   f.GuestPrice,
		GuestLimit:   f.GuestLimit,
		Tickets:      f.Tickets,
		GuestTickets: f.GuestTickets,
		SaleStart:    f.SaleStart,
		SaleEnd:      f.SaleEnd,
		DateTime:     f.DateTime,
	}
}

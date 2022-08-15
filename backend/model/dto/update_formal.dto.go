package dto

import (
	"time"

	"github.com/kcsu/store/model"
)

type UpdateFormalDto struct {
	Name                   string    `json:"name" validate:"required,min=5,max=100"`
	Menu                   string    `json:"menu"` // TODO: validate menu?
	Price                  float32   `json:"price" validate:"min=0"`
	GuestPrice             float32   `json:"guestPrice" validate:"min=0"`
	GuestLimit             uint      `json:"guestLimit"`
	FirstSaleTickets       uint      `json:"firstSaleTickets"`
	FirstSaleGuestTickets  uint      `json:"firstSaleGuestTickets"`
	FirstSaleStart         time.Time `json:"firstSaleStart" validate:"required"`
	SecondSaleTickets      uint      `json:"secondSaleTickets"`
	SecondSaleGuestTickets uint      `json:"secondSaleGuestTickets"`
	SecondSaleStart        time.Time `json:"secondSaleStart" validate:"required"`
	SaleEnd                time.Time `json:"saleEnd" validate:"required"`  // TODO: gt SaleStart?
	DateTime               time.Time `json:"dateTime" validate:"required"` // TODO: gt SaleEnd?
	HasGuestList           bool      `json:"hasGuestList"`
	IsVisible              bool      `json:"isVisible"`
}

func (f *UpdateFormalDto) Formal() model.Formal {
	return model.Formal{
		Name:                   f.Name,
		Menu:                   f.Menu,
		Price:                  f.Price,
		GuestPrice:             f.GuestPrice,
		GuestLimit:             f.GuestLimit,
		FirstSaleTickets:       f.FirstSaleTickets,
		FirstSaleGuestTickets:  f.FirstSaleGuestTickets,
		FirstSaleStart:         f.FirstSaleStart,
		SecondSaleTickets:      f.SecondSaleTickets,
		SecondSaleGuestTickets: f.SecondSaleGuestTickets,
		SecondSaleStart:        f.SecondSaleStart,
		SaleEnd:                f.SaleEnd,
		DateTime:               f.DateTime,
		HasGuestList:           f.HasGuestList,
		IsVisible:              f.IsVisible,
	}
}

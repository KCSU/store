package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kcsu/store/handlers"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const expectedJSON = `[
	{
		"id": 1,
		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z",
		"deletedAt": null,
		"name": "Test 1",
		"menu": "A menu",
		"price": 21.3,
		"guestPrice": 11.6,
		"guestLimit": 0,
		"tickets": 0,
		"guestTickets": 0,
		"saleStart": "0001-01-01T00:00:00Z",
		"saleEnd": "0001-01-01T00:00:00Z",
		"dateTime": "0001-01-01T00:00:00Z",
		"ticketsRemaining": 24,
		"guestTicketsRemaining": 56
	},
	{
		"id": 6,
		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z",
		"deletedAt": null,
		"name": "Test 2",
		"menu": "Another menu",
		"price": 15.6,
		"guestPrice": 27.2,
		"guestLimit": 0,
		"tickets": 0,
		"guestTickets": 0,
		"saleStart": "0001-01-01T00:00:00Z",
		"saleEnd": "0001-01-01T00:00:00Z",
		"dateTime": "0001-01-01T00:00:00Z",
		"ticketsRemaining": 64,
		"guestTicketsRemaining": 31
	}
]`

func TestGetFormals(t *testing.T) {
	// Init handler
	h := new(Handler)
	f := new(mocks.FormalStore)
	h.Formals = f
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock database
	formals := []model.Formal{
		{
			Model:      model.Model{ID: 1},
			Name:       "Test 1",
			Menu:       "A menu",
			Price:      21.3,
			GuestPrice: 11.6,
		},
		{
			Model:      model.Model{ID: 6},
			Name:       "Test 2",
			Menu:       "Another menu",
			Price:      15.6,
			GuestPrice: 27.2,
		},
	}
	f.On("Get").Return(formals, nil)
	f.On("TicketsRemaining", &formals[0], true).Return(uint(56))
	f.On("TicketsRemaining", &formals[0], false).Return(uint(24))
	f.On("TicketsRemaining", &formals[1], true).Return(uint(31))
	f.On("TicketsRemaining", &formals[1], false).Return(uint(64))

	// Run test
	err := h.GetFormals(c)
	assert.NoError(t, err)
	f.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expectedJSON, rec.Body.String())
}

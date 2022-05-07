package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	. "github.com/kcsu/store/handlers"
	am "github.com/kcsu/store/mocks/auth"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetFormals(t *testing.T) {
	const expectedJSON = `[
		{
			"id": "215292b8-4911-4d93-81dd-ebafb1aa6489",
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
			"guestTicketsRemaining": 56,
			"groups": [],
			"myTickets": [{
				"id": "8bc2da87-88ea-4fcf-a69f-22a360b2606a",
				"createdAt": "0001-01-01T00:00:00Z",
				"updatedAt": "0001-01-01T00:00:00Z",
				"deletedAt": null,
				"isGuest": false,
				"isQueue": false,
				"formalId": "215292b8-4911-4d93-81dd-ebafb1aa6489",
				"userId": "290bc3be-12f6-48a0-b624-32b2eb5e05c9",
				"option": "Vegan"
			}]
		},
		{
			"id": "202ca011-4cf8-4bf4-b318-c644be23ba85",
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
			"guestTicketsRemaining": 31,
			"groups": [
				{
					"id": "ce22bfb1-b932-4073-a6a9-01e4787deecb",
					"name": "Group A"
				},
				{
					"id": "bcdb4d1c-7deb-49c0-aaaa-8adb053ecfc2",
					"name": "Group B"
				}
			]
		}
	]`
	uid := uuid.MustParse("290bc3be-12f6-48a0-b624-32b2eb5e05c9")
	// Init handler
	h := new(Handler)
	f := mocks.NewFormalStore(t)
	a := am.NewAuth(t)
	h.Formals = f
	h.Auth = a
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock database
	formals := []model.Formal{
		{
			Model:      model.Model{ID: uuid.MustParse("215292b8-4911-4d93-81dd-ebafb1aa6489")},
			Name:       "Test 1",
			Menu:       "A menu",
			Price:      21.3,
			GuestPrice: 11.6,
			TicketSales: []model.Ticket{{
				Model:      model.Model{ID: uuid.MustParse("8bc2da87-88ea-4fcf-a69f-22a360b2606a")},
				FormalID:   uuid.MustParse("215292b8-4911-4d93-81dd-ebafb1aa6489"),
				UserID:     uid,
				MealOption: "Vegan",
			}},
		},
		{
			Model:      model.Model{ID: uuid.MustParse("202ca011-4cf8-4bf4-b318-c644be23ba85")},
			Name:       "Test 2",
			Menu:       "Another menu",
			Price:      15.6,
			GuestPrice: 27.2,
			Groups: []model.Group{
				{
					Model: model.Model{ID: uuid.MustParse("ce22bfb1-b932-4073-a6a9-01e4787deecb")},
					Name:  "Group A",
				},
				{
					Model: model.Model{ID: uuid.MustParse("bcdb4d1c-7deb-49c0-aaaa-8adb053ecfc2")},
					Name:  "Group B",
				},
			},
		},
	}
	// FIXME: refactor to make it easier to add cases?
	a.On("GetUserId", c).Return(uid)
	f.On("GetWithUserData", uid).Return(formals, nil)
	f.On("TicketsRemaining", &formals[0], true).Return(uint(56))
	f.On("TicketsRemaining", &formals[0], false).Return(uint(24))
	f.On("TicketsRemaining", &formals[1], true).Return(uint(31))
	f.On("TicketsRemaining", &formals[1], false).Return(uint(64))

	// Run test
	err := h.GetFormals(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expectedJSON, rec.Body.String())
}

func TestGetFormalGuestList(t *testing.T) {
	expectedJSON := `[
		{
			"name": "Test 1",
			"email": "test1@cam.ac.uk",
			"guests": 3
		},
		{
			"name": "Test 2",
			"email": "test2@cam.ac.uk",
			"guests": 1
		}
	]`
	id := uuid.New()
	// Init handler
	h := new(Handler)
	f := mocks.NewFormalStore(t)
	h.Formals = f
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprint("/formals/", id.String(), "/guests"),
		nil,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())
	// Mock database
	guests := []model.FormalGuest{
		{
			Name:   "Test 1",
			Email:  "test1@cam.ac.uk",
			Guests: 3,
		},
		{
			Name:   "Test 2",
			Email:  "test2@cam.ac.uk",
			Guests: 1,
		},
	}
	f.On("FindGuestList", id).Return(guests, nil)
	// Run test
	err := h.GetFormalGuestList(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expectedJSON, rec.Body.String())
}

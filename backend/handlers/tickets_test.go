package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/kcsu/store/handlers"
	"github.com/kcsu/store/middleware"
	am "github.com/kcsu/store/mocks/auth"
	sm "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TicketSuite struct {
	suite.Suite
	h        *Handler
	auth     *am.Auth
	users    *sm.UserStore
	tickets  *sm.TicketStore
	formals  *sm.FormalStore
	e        *echo.Echo
	mockUser *model.User
}

var userId = uuid.MustParse("9cb93e84-43ac-456f-972b-71ffce3e6782")

func (t *TicketSuite) SetupTest() {
	t.auth = am.NewAuth(t.T())
	t.users = sm.NewUserStore(t.T())
	t.tickets = sm.NewTicketStore(t.T())
	t.formals = sm.NewFormalStore(t.T())
	t.h = &Handler{
		Auth:    t.auth,
		Users:   t.users,
		Tickets: t.tickets,
		Formals: t.formals,
	}
	// Config?
	// Some mocks we always want to use
	t.auth.On("GetUserId", mock.Anything).Maybe().Return(userId)
	t.mockUser = &model.User{
		Model: model.Model{
			ID: userId,
		},
		Name:           "Logan Roy",
		Email:          "lry555@cam.ac.uk",
		ProviderUserId: "abcdefg",
	}
	t.e = echo.New()
	// Should this be empty for unit tests? A mock nil-validator?
	t.e.Validator = middleware.NewValidator()
	t.users.On("Find", userId).Maybe().Return(*t.mockUser, nil)
}

func (t *TicketSuite) TestGetTickets() {
	formals := []model.Formal{{
		Name:  "Test Formal",
		Model: model.Model{ID: uuid.MustParse("d439376d-abd1-44c1-831c-e8d4565cac5a")},
		DateTime: time.Date(
			2022, time.January, 17,
			19, 30, 0, 0, time.UTC,
		),
	}, {
		Name:  "Another Formal",
		Model: model.Model{ID: uuid.MustParse("7d35e4c8-0603-4723-8b9c-756814d1c545")},
		DateTime: time.Date(
			2022, time.January, 19,
			19, 30, 0, 0, time.UTC,
		),
	}}
	tickets := []model.Ticket{
		{
			Model:      model.Model{ID: uuid.MustParse("50e20619-6383-48eb-aaaf-ad3c15ec919f")},
			UserID:     userId,
			IsGuest:    true,
			IsQueue:    true,
			FormalID:   formals[0].ID,
			Formal:     &formals[0],
			MealOption: "Pescetarian",
		},
		{
			Model:      model.Model{ID: uuid.MustParse("62fc9222-3629-4914-ba40-d4f8a46e0ddd")},
			UserID:     userId,
			IsGuest:    false,
			FormalID:   formals[1].ID,
			Formal:     &formals[1],
			MealOption: "Vegetarian",
		},
		{
			Model:      model.Model{ID: uuid.MustParse("d46c4abb-5e92-414a-aeef-8f7cc10261fb")},
			UserID:     userId,
			IsGuest:    false,
			FormalID:   formals[0].ID,
			Formal:     &formals[0],
			MealOption: "Normal",
		},
	}
	wantsJson := `[
		{
			"formal":{
				"id": "d439376d-abd1-44c1-831c-e8d4565cac5a",
				"createdAt":"0001-01-01T00:00:00Z",
				"updatedAt":"0001-01-01T00:00:00Z",
				"deletedAt":null,
				"name":"Test Formal",
				"menu":"",
				"price":0,
				"guestPrice":0,
				"guestLimit":0,
				"tickets":0,
				"guestTickets":0,
				"hasGuestList":false,
				"saleStart":"0001-01-01T00:00:00Z",
				"saleEnd":"0001-01-01T00:00:00Z",
				"dateTime":"2022-01-17T19:30:00Z"
			},
			"ticket":{
				"id": "d46c4abb-5e92-414a-aeef-8f7cc10261fb",
				"createdAt":"0001-01-01T00:00:00Z",
				"updatedAt":"0001-01-01T00:00:00Z",
				"deletedAt":null,
				"isGuest":false,
				"isQueue":false,
				"option":"Normal",
				"formalId": "d439376d-abd1-44c1-831c-e8d4565cac5a",
				"userId": "9cb93e84-43ac-456f-972b-71ffce3e6782"
			},
			"guestTickets":[
				{
					"id": "50e20619-6383-48eb-aaaf-ad3c15ec919f",
					"createdAt":"0001-01-01T00:00:00Z",
					"updatedAt":"0001-01-01T00:00:00Z",
					"deletedAt":null,
					"isGuest":true,
					"isQueue":true,
					"option":"Pescetarian",
					"formalId": "d439376d-abd1-44c1-831c-e8d4565cac5a",
					"userId": "9cb93e84-43ac-456f-972b-71ffce3e6782"
				}
			]
		},
		{
			"formal":{
				"id": "7d35e4c8-0603-4723-8b9c-756814d1c545",
				"createdAt":"0001-01-01T00:00:00Z",
				"updatedAt":"0001-01-01T00:00:00Z",
				"deletedAt":null,
				"name":"Another Formal",
				"menu":"",
				"price":0,
				"guestPrice":0,
				"guestLimit":0,
				"tickets":0,
				"guestTickets":0,
				"hasGuestList":false,
				"saleStart":"0001-01-01T00:00:00Z",
				"saleEnd":"0001-01-01T00:00:00Z",
				"dateTime":"2022-01-19T19:30:00Z"
			},
			"ticket":{
				"id": "62fc9222-3629-4914-ba40-d4f8a46e0ddd",
				"createdAt":"0001-01-01T00:00:00Z",
				"updatedAt":"0001-01-01T00:00:00Z",
				"deletedAt":null,
				"isGuest":false,
				"isQueue":false,
				"option":"Vegetarian",
				"formalId": "7d35e4c8-0603-4723-8b9c-756814d1c545",
				"userId": "9cb93e84-43ac-456f-972b-71ffce3e6782"
			},
			"guestTickets":[

			]
		}
	]`
	// HT
	req := httptest.NewRequest(http.MethodGet, "/tickets", nil)
	rec := httptest.NewRecorder()
	c := t.e.NewContext(req, rec)
	t.tickets.On("Get", userId).Return(tickets, nil)
	err := t.h.GetTickets(c)
	t.NoError(err)
	t.Equal(http.StatusOK, rec.Code)
	t.JSONEq(wantsJson, rec.Body.String())
}

func (t *TicketSuite) TestBuyTicket() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		body   string
		formal model.Formal
		exists bool
		wants  *wants
	}
	userGroups := []model.Group{
		{Model: model.Model{ID: uuid.MustParse("d5014dfb-1b6d-4b20-8da5-f753991e68bc")}},
		{Model: model.Model{ID: uuid.MustParse("4369f4be-9ad8-4a03-acbf-ebe5f96a0048")}},
	}
	tests := []test{
		{
			"Insufficient Groups",
			`{
				"formalId": "c9da6c4e-965d-4ac9-bb24-2f035d1be0a6",
				"ticket": {"option": "Normal"},
				"guestTickets": []
			}`,
			model.Formal{
				Model: model.Model{ID: uuid.MustParse("c9da6c4e-965d-4ac9-bb24-2f035d1be0a6")},
				Name:  "Wrong Group Formal",
				Groups: []model.Group{
					{Model: model.Model{ID: uuid.New()}},
					{Model: model.Model{ID: uuid.New()}},
				},
				SaleEnd: time.Now().AddDate(0, 0, 1),
			},
			false,
			&wants{http.StatusForbidden, "Forbidden"},
		},
		{
			"Guest Limit",
			`{
				"formalId": "4b47d6cd-3a0a-40dc-b066-d8711e46d8cb",
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]
			}`,
			model.Formal{
				Model:      model.Model{ID: uuid.MustParse("4b47d6cd-3a0a-40dc-b066-d8711e46d8cb")},
				Name:       "Wrong Number Formal",
				GuestLimit: 1,
				Groups: []model.Group{
					userGroups[1],
					{Model: model.Model{ID: uuid.New()}},
				},
				SaleEnd: time.Now().AddDate(0, 0, 1),
			},
			false,
			&wants{http.StatusUnprocessableEntity, "Too many guest tickets requested."},
		},
		{
			"Duplicate Ticket",
			`{
				"formalId": "1052f03b-0d86-41af-a00b-2a632a21ae01",
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]
			}`,
			model.Formal{
				Model:      model.Model{ID: uuid.MustParse("1052f03b-0d86-41af-a00b-2a632a21ae01")},
				Name:       "Existing Formal",
				GuestLimit: 2,
				Groups: []model.Group{
					userGroups[1],
					userGroups[0],
				},
				SaleEnd: time.Now().AddDate(0, 0, 1),
			},
			true,
			&wants{http.StatusConflict, "Ticket already exists."},
		},
		{
			"Sales Closed",
			`{
				"formalId": "cbc4d1e8-0036-4677-90eb-57d069fd1217",
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]
			}`,
			model.Formal{
				Model:      model.Model{ID: uuid.MustParse("cbc4d1e8-0036-4677-90eb-57d069fd1217")},
				Name:       "Wrong Group Formal",
				GuestLimit: 3,
				Groups: []model.Group{
					userGroups[0],
					{Model: model.Model{ID: uuid.New()}},
				},
				SaleEnd: time.Now().AddDate(0, 0, -3),
			},
			false,
			&wants{http.StatusUnprocessableEntity, "Sales have closed."},
		},
		{
			"Should Create",
			`{
				"formalId": "7ae71984-26b0-49bc-883e-8a77b90d00c1",
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]
			}`,
			model.Formal{
				Model:      model.Model{ID: uuid.MustParse("7ae71984-26b0-49bc-883e-8a77b90d00c1")},
				Name:       "Wrong Group Formal",
				GuestLimit: 3,
				Groups: []model.Group{
					userGroups[0],
					{Model: model.Model{ID: uuid.New()}},
				},
				SaleEnd: time.Now().AddDate(0, 0, 1),
			},
			false,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func() {
			// HTTP
			req := httptest.NewRequest(
				http.MethodPost, "/tickets", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := t.e.NewContext(req, rec)
			// Mock
			t.formals.On("Find", test.formal.ID).Return(test.formal, nil).Once()
			t.users.On("Groups", t.mockUser).Return(userGroups, nil)
			t.tickets.On(
				"ExistsByFormal", test.formal.ID, userId,
			).Maybe().Return(test.exists, nil)
			if test.wants == nil {
				t.tickets.On(
					"BatchCreate", mock.AnythingOfType("[]model.Ticket"),
				).Return(nil).Once()
			}
			// Test
			err := t.h.BuyTicket(c)
			if test.wants != nil {
				var he *echo.HTTPError
				if t.ErrorAs(err, &he) {
					t.Equal(test.wants.code, he.Code)
					t.Equal(test.wants.message, he.Message)
				}
			} else {
				t.NoError(err)
				t.Equal(http.StatusCreated, rec.Code)
			}
		})
	}
}

func (t *TicketSuite) TestCancelTickets() {
	// HTTP
	req := new(http.Request)
	rec := httptest.NewRecorder()
	c := t.e.NewContext(req, rec)
	c.SetParamNames("id")
	formalId := uuid.New()
	c.SetParamValues(formalId.String())
	t.formals.On("Find", formalId).Return(model.Formal{
		SaleEnd: time.Now().AddDate(0, 0, 7),
	}, nil)
	// Mock
	t.tickets.On("DeleteByFormal", formalId, userId).Return(nil)
	err := t.h.CancelTickets(c)
	t.NoError(err)
	t.Equal(http.StatusOK, rec.Code)
}

func (t *TicketSuite) TestCancelTicket() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		ticket model.Ticket
		wants  *wants
	}
	tests := []test{
		{"Incorrect User", model.Ticket{
			UserID:   uuid.New(),
			FormalID: uuid.New(),
			IsGuest:  true,
			IsQueue:  false,
		}, &wants{http.StatusForbidden, "Forbidden"}},
		{"Forbid Non-Guest", model.Ticket{
			UserID:   userId,
			FormalID: uuid.New(),
			IsGuest:  false,
			IsQueue:  true,
		}, &wants{http.StatusForbidden, "Non-guest tickets must be cancelled as a group"}},
		{"Sales Closed", model.Ticket{
			UserID:   userId,
			FormalID: uuid.New(),
			IsGuest:  true,
			IsQueue:  true,
			Formal:   &model.Formal{SaleEnd: time.Now().AddDate(0, 0, -5)},
		}, &wants{http.StatusUnprocessableEntity, "Sales have closed."}},
		{"Should Cancel", model.Ticket{
			UserID:   userId,
			FormalID: uuid.New(),
			IsGuest:  true,
			IsQueue:  true,
		}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func() {
			// HTTP
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := t.e.NewContext(req, rec)
			c.SetParamNames("id")
			ticketId := uuid.New()
			c.SetParamValues(ticketId.String())
			// Mock
			ticket := test.ticket
			if ticket.Formal == nil {
				ticket.Formal = &model.Formal{
					SaleEnd: time.Now().AddDate(0, 0, 1),
				}
			}
			t.tickets.On("FindWithFormal", ticketId).Return(ticket, nil).Once()
			if test.wants == nil {
				t.tickets.On("Delete", ticketId).Return(nil).Once()
			}
			// Test
			err := t.h.CancelTicket(c)
			if test.wants != nil {
				var he *echo.HTTPError
				if t.ErrorAs(err, &he) {
					t.Equal(test.wants.code, he.Code)
					t.Equal(test.wants.message, he.Message)
				}
			} else {
				t.NoError(err)
				t.Equal(http.StatusOK, rec.Code)
			}
		})
	}
}

func (t *TicketSuite) TestEditTicket() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		ticket model.Ticket
		option string
		wants  *wants
	}
	tests := []test{
		{"Incorrect User", model.Ticket{
			UserID:   uuid.New(),
			FormalID: uuid.New(),
			IsGuest:  true,
			IsQueue:  false,
		}, "Vegetarian", &wants{http.StatusForbidden, "Forbidden"}},
		{"Sale Closed", model.Ticket{
			UserID:   userId,
			FormalID: uuid.New(),
			IsGuest:  true,
			IsQueue:  true,
			Formal: &model.Formal{
				SaleEnd: time.Now().AddDate(0, 0, -7),
			},
		}, "Normal", &wants{http.StatusUnprocessableEntity, "Sales have closed."}},
		{"Should Update", model.Ticket{
			UserID:   userId,
			FormalID: uuid.New(),
			IsGuest:  true,
			IsQueue:  true,
		}, "Pescetarian", nil},
	}
	// HTTP
	for _, test := range tests {
		t.Run(test.name, func() {
			ticketId := uuid.New()
			route := fmt.Sprintf("/tickets/%s", ticketId.String())
			body, err := json.Marshal(map[string]string{
				"option": test.option,
			})
			t.Require().NoError(err)
			req := httptest.NewRequest(
				http.MethodPut, route, bytes.NewReader(body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := t.e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(ticketId.String())
			// Mock
			test.ticket.ID = ticketId
			if test.ticket.Formal == nil {
				test.ticket.Formal = &model.Formal{
					SaleEnd: time.Now().AddDate(0, 0, 12),
				}
			}
			t.tickets.On("FindWithFormal", ticketId).Return(test.ticket, nil).Once()
			if test.wants == nil {
				t.tickets.On(
					"Update", ticketId,
					&dto.TicketRequestDto{MealOption: test.option},
				).Return(nil).Once()
			}
			// Test
			err = t.h.EditTicket(c)
			if test.wants != nil {
				var he *echo.HTTPError
				if t.ErrorAs(err, &he) {
					t.Equal(test.wants.code, he.Code)
					t.Equal(test.wants.message, he.Message)
				}
			} else {
				t.NoError(err)
				t.Equal(http.StatusOK, rec.Code)
			}
		})
	}
}

func (t *TicketSuite) TestAddTicket() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		formal model.Formal
		option string
		wants  *wants
	}
	tests := []test{{
		"Should Add",
		model.Formal{
			GuestLimit: 4,
			SaleEnd:    time.Now().AddDate(0, 0, 1),
		},
		"Vegetarian",
		nil,
	}, {
		"Sale Closed",
		model.Formal{
			GuestLimit: 4,
			SaleEnd:    time.Now().AddDate(0, 0, -1),
		},
		"Normal",
		&wants{http.StatusUnprocessableEntity, "Sales have closed."},
	}, {
		"Guest Limit",
		model.Formal{
			GuestLimit: 2,
			SaleEnd:    time.Now().AddDate(0, 0, 1),
		},
		"Pescetarian",
		&wants{http.StatusUnprocessableEntity, "Too many guest tickets requested."},
	}}
	for _, test := range tests {
		t.Run(test.name, func() {
			// HTTP
			formalId := uuid.New()
			path := fmt.Sprintf("/formals/%s", formalId.String())
			body, err := json.Marshal(map[string]string{
				"option": test.option,
			})
			t.Require().NoError(err)
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := t.e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(formalId.String())
			// Mock
			t.tickets.On("ExistsByFormal", formalId, userId).Return(true, nil).Once()
			// The user already has 2 guest tickets
			t.tickets.On("CountGuestByFormal", formalId, userId).Return(int64(2), nil).Once()
			t.formals.On("Find", formalId).Return(test.formal, nil).Once()
			if test.wants == nil {
				t.tickets.On("Create", &model.Ticket{
					FormalID:   formalId,
					IsGuest:    true,
					IsQueue:    true,
					MealOption: test.option,
					UserID:     userId,
				}).Return(nil).Once()
			}
			// Test
			err = t.h.AddTicket(c)
			if test.wants != nil {
				var he *echo.HTTPError
				if t.ErrorAs(err, &he) {
					t.Equal(test.wants.code, he.Code)
					t.Equal(test.wants.message, he.Message)
				}
			} else {
				t.NoError(err)
				t.Equal(http.StatusCreated, rec.Code)
			}
		})
	}
}

func TestTicketSuite(t *testing.T) {
	suite.Run(t, new(TicketSuite))
}

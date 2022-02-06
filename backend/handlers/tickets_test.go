package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/kcsu/store/handlers"
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
	mockUser *model.User
}

func (t *TicketSuite) AssertExpectations() {
	t.auth.AssertExpectations(t.T())
	t.users.AssertExpectations(t.T())
	t.tickets.AssertExpectations(t.T())
	t.formals.AssertExpectations(t.T())
}

const userId = 5

func (t *TicketSuite) SetupTest() {
	t.auth = new(am.Auth)
	t.users = new(sm.UserStore)
	t.tickets = new(sm.TicketStore)
	t.formals = new(sm.FormalStore)
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
			ID: uint(userId),
		},
		Name:           "Logan Roy",
		Email:          "lry555@cam.ac.uk",
		ProviderUserId: "abcdefg",
	}
	t.users.On("Find", userId).Maybe().Return(*t.mockUser, nil)
}

func (t *TicketSuite) TestGetTickets() {
	formals := []model.Formal{{
		Name:  "Test Formal",
		Model: model.Model{ID: 52},
		DateTime: time.Date(
			2022, time.January, 17,
			19, 30, 0, 0, time.UTC,
		),
	}, {
		Name:  "Another Formal",
		Model: model.Model{ID: 56},
		DateTime: time.Date(
			2022, time.January, 19,
			19, 30, 0, 0, time.UTC,
		),
	}}
	tickets := []model.Ticket{
		{
			Model:      model.Model{ID: 6},
			UserId:     userId,
			IsGuest:    true,
			IsQueue:    true,
			FormalID:   int(formals[0].ID),
			Formal:     &formals[0],
			MealOption: "Pescetarian",
		},
		{
			Model:      model.Model{ID: 47},
			UserId:     userId,
			IsGuest:    false,
			FormalID:   int(formals[1].ID),
			Formal:     &formals[1],
			MealOption: "Vegetarian",
		},
		{
			Model:      model.Model{ID: 91},
			UserId:     userId,
			IsGuest:    false,
			FormalID:   int(formals[0].ID),
			Formal:     &formals[0],
			MealOption: "Normal",
		},
	}
	wantsJson := `[
		{
			"formal":{
				"id":52,
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
				"saleStart":"0001-01-01T00:00:00Z",
				"saleEnd":"0001-01-01T00:00:00Z",
				"dateTime":"2022-01-17T19:30:00Z"
			},
			"ticket":{
				"id":91,
				"createdAt":"0001-01-01T00:00:00Z",
				"updatedAt":"0001-01-01T00:00:00Z",
				"deletedAt":null,
				"isGuest":false,
				"isQueue":false,
				"option":"Normal",
				"formalId":52,
				"userId":5
			},
			"guestTickets":[
				{
					"id":6,
					"createdAt":"0001-01-01T00:00:00Z",
					"updatedAt":"0001-01-01T00:00:00Z",
					"deletedAt":null,
					"isGuest":true,
					"isQueue":true,
					"option":"Pescetarian",
					"formalId":52,
					"userId":5
				}
			]
		},
		{
			"formal":{
				"id":56,
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
				"saleStart":"0001-01-01T00:00:00Z",
				"saleEnd":"0001-01-01T00:00:00Z",
				"dateTime":"2022-01-19T19:30:00Z"
			},
			"ticket":{
				"id":47,
				"createdAt":"0001-01-01T00:00:00Z",
				"updatedAt":"0001-01-01T00:00:00Z",
				"deletedAt":null,
				"isGuest":false,
				"isQueue":false,
				"option":"Vegetarian",
				"formalId":56,
				"userId":5
			},
			"guestTickets":[
				
			]
		}
	]`
	// HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tickets", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	t.tickets.On("Get", userId).Return(tickets, nil)
	err := t.h.GetTickets(c)
	t.NoError(err)
	t.Equal(http.StatusOK, rec.Code)
	t.JSONEq(wantsJson, rec.Body.String())
	t.AssertExpectations()
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
		{Model: model.Model{ID: 1}},
		{Model: model.Model{ID: 3}},
	}
	tests := []test{
		{
			"Insufficient Groups",
			`{
				"formalId": 1,
				"ticket": {"option": "Normal"},
				"guestTickets": []	
			}`,
			model.Formal{
				Model: model.Model{ID: 1},
				Name:  "Wrong Group Formal",
				Groups: []model.Group{
					{Model: model.Model{ID: 2}},
					{Model: model.Model{ID: 4}},
				},
				SaleEnd: time.Now().AddDate(0, 0, 1),
			},
			false,
			&wants{http.StatusForbidden, "Forbidden"},
		},
		{
			"Guest Limit",
			`{
				"formalId": 2,
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]	
			}`,
			model.Formal{
				Model:      model.Model{ID: 2},
				Name:       "Wrong Number Formal",
				GuestLimit: 1,
				Groups: []model.Group{
					{Model: model.Model{ID: 3}},
					{Model: model.Model{ID: 4}},
				},
				SaleEnd: time.Now().AddDate(0, 0, 1),
			},
			false,
			&wants{http.StatusUnprocessableEntity, "Too many guest tickets requested."},
		},
		{
			"Duplicate Ticket",
			`{
				"formalId": 3,
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]	
			}`,
			model.Formal{
				Model:      model.Model{ID: 3},
				Name:       "Existing Formal",
				GuestLimit: 2,
				Groups: []model.Group{
					{Model: model.Model{ID: 3}},
					{Model: model.Model{ID: 1}},
				},
				SaleEnd: time.Now().AddDate(0, 0, 1),
			},
			true,
			&wants{http.StatusConflict, "Ticket already exists."},
		},
		{
			"Sales Closed",
			`{
				"formalId": 4,
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]	
			}`,
			model.Formal{
				Model:      model.Model{ID: 4},
				Name:       "Wrong Group Formal",
				GuestLimit: 3,
				Groups: []model.Group{
					{Model: model.Model{ID: 1}},
					{Model: model.Model{ID: 2}},
				},
				SaleEnd: time.Now().AddDate(0, 0, -3),
			},
			false,
			&wants{http.StatusUnprocessableEntity, "Sales have closed."},
		},
		{
			"Should Create",
			`{
				"formalId": 4,
				"ticket": {"option": "Normal"},
				"guestTickets": [
					{"option": "Pescetarian"},
					{"option": "Vegan"}
				]	
			}`,
			model.Formal{
				Model:      model.Model{ID: 4},
				Name:       "Wrong Group Formal",
				GuestLimit: 3,
				Groups: []model.Group{
					{Model: model.Model{ID: 1}},
					{Model: model.Model{ID: 2}},
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
			e := echo.New()
			req := httptest.NewRequest(
				http.MethodPost, "/tickets", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			t.formals.On("Find", int(test.formal.ID)).Return(test.formal, nil).Once()
			t.users.On("Groups", t.mockUser).Return(userGroups, nil)
			t.tickets.On(
				"ExistsByFormal", int(test.formal.ID), userId,
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
	t.AssertExpectations()
}

func (t *TicketSuite) TestCancelTickets() {
	// HTTP
	e := echo.New()
	req := new(http.Request)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	formalId := 7
	c.SetParamValues(strconv.Itoa(formalId))
	t.formals.On("Find", formalId).Return(model.Formal{
		SaleEnd: time.Now().AddDate(0, 0, 7),
	}, nil)
	// Mock
	t.tickets.On("DeleteByFormal", formalId, userId).Return(nil)
	err := t.h.CancelTickets(c)
	t.NoError(err)
	t.Equal(http.StatusOK, rec.Code)
	t.AssertExpectations()
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
			UserId:   userId + 1,
			FormalID: 53,
			IsGuest:  true,
			IsQueue:  false,
		}, &wants{http.StatusForbidden, "Forbidden"}},
		{"Forbid Non-Guest", model.Ticket{
			UserId:   userId,
			FormalID: 54,
			IsGuest:  false,
			IsQueue:  true,
		}, &wants{http.StatusForbidden, "Non-guest tickets must be cancelled as a group"}},
		{"Sales Closed", model.Ticket{
			UserId:   userId,
			FormalID: 55,
			IsGuest:  true,
			IsQueue:  true,
			Formal:   &model.Formal{SaleEnd: time.Now().AddDate(0, 0, -5)},
		}, &wants{http.StatusUnprocessableEntity, "Sales have closed."}},
		{"Should Cancel", model.Ticket{
			UserId:   userId,
			FormalID: 55,
			IsGuest:  true,
			IsQueue:  true,
		}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func() {
			// HTTP
			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			ticketId := 84
			c.SetParamValues(strconv.Itoa(ticketId))
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
	t.AssertExpectations()
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
			UserId:   userId + 1,
			FormalID: 53,
			IsGuest:  true,
			IsQueue:  false,
		}, "Vegetarian", &wants{http.StatusForbidden, "Forbidden"}},
		{"Sale Closed", model.Ticket{
			UserId:   userId,
			FormalID: 55,
			IsGuest:  true,
			IsQueue:  true,
			Formal: &model.Formal{
				SaleEnd: time.Now().AddDate(0, 0, -7),
			},
		}, "Normal", &wants{http.StatusUnprocessableEntity, "Sales have closed."}},
		{"Should Update", model.Ticket{
			UserId:   userId,
			FormalID: 55,
			IsGuest:  true,
			IsQueue:  true,
		}, "Pescetarian", nil},
	}
	// HTTP
	for _, test := range tests {
		t.Run(test.name, func() {
			ticketId := 95
			route := fmt.Sprintf("/tickets/%d", ticketId)
			body, err := json.Marshal(map[string]string{
				"option": test.option,
			})
			t.Require().NoError(err)
			e := echo.New()
			req := httptest.NewRequest(
				http.MethodPut, route, bytes.NewReader(body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(ticketId))
			// Mock
			test.ticket.ID = uint(ticketId)
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
	t.AssertExpectations()
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
			formalId := 56
			path := fmt.Sprintf("/formals/%d", formalId)
			e := echo.New()
			body, err := json.Marshal(map[string]string{
				"option": test.option,
			})
			t.Require().NoError(err)
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(formalId))
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
					UserId:     userId,
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
	t.AssertExpectations()
}

func TestTicketSuite(t *testing.T) {
	suite.Run(t, new(TicketSuite))
}

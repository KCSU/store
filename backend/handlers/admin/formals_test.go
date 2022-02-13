package admin_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/kcsu/store/handlers/admin"
	"github.com/kcsu/store/middleware"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type AdminFormalSuite struct {
	suite.Suite
	h       *AdminHandler
	formals *mocks.FormalStore
}

func (s *AdminFormalSuite) SetupTest() {
	// Init handler
	s.h = new(AdminHandler)
	s.formals = new(mocks.FormalStore)
	s.h.Formals = s.formals
}

func (s *AdminFormalSuite) TestGetFormals() {
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
			"guestTicketsRemaining": 56,
			"groups": []
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
			"guestTicketsRemaining": 31,
			"groups": [
				{
					"id": 2,
					"name": "Group A"
				},
				{
					"id": 4,
					"name": "Group B"
				}
			]
		}
	]`
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
			Groups: []model.Group{
				{
					Model: model.Model{ID: 2},
					Name:  "Group A",
				},
				{
					Model: model.Model{ID: 4},
					Name:  "Group B",
				},
			},
		},
	}
	// FIXME: refactor to make it easier to add cases?
	s.formals.On("All").Return(formals, nil)
	s.formals.On("TicketsRemaining", &formals[0], true).Return(uint(56))
	s.formals.On("TicketsRemaining", &formals[0], false).Return(uint(24))
	s.formals.On("TicketsRemaining", &formals[1], true).Return(uint(31))
	s.formals.On("TicketsRemaining", &formals[1], false).Return(uint(64))

	// Run test
	err := s.h.GetFormals(c)
	s.NoError(err)
	s.formals.AssertExpectations(s.T())
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminFormalSuite) TestGetFormal() {
	const expectedJSON = `{
		"id": 13,
		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z",
		"deletedAt": null,
		"name": "Test 5",
		"menu": "Another menu",
		"price": 26.3,
		"guestPrice": 12.7,
		"guestLimit": 3,
		"tickets": 0,
		"guestTickets": 0,
		"saleStart": "0001-01-01T00:00:00Z",
		"saleEnd": "0001-01-01T00:00:00Z",
		"dateTime": "0001-01-01T00:00:00Z",
		"ticketsRemaining": 20,
		"guestTicketsRemaining": 36,
		"groups": [{
			"id": 5,
			"name": "Group"
		}]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/formals/13", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("13")
	formal := model.Formal{
		Model:      model.Model{ID: 13},
		Name:       "Test 5",
		Menu:       "Another menu",
		Price:      26.3,
		GuestPrice: 12.7,
		GuestLimit: 3,
		Groups: []model.Group{{
			Model: model.Model{ID: 5},
			Name:  "Group",
		}},
	}
	s.formals.On("Find", 13).Return(formal, nil)
	s.formals.On("TicketsRemaining", &formal, true).Return(uint(36))
	s.formals.On("TicketsRemaining", &formal, false).Return(uint(20))

	err := s.h.GetFormal(c)
	s.NoError(err)
	s.formals.AssertExpectations(s.T())
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminFormalSuite) TestCreateFormal() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		body   string
		groups []int
		formal *model.Formal
		wants  *wants
	}
	tests := []test{
		{
			"Invalid Groups",
			`{
				"name": "A formal",
				"menu": "Some menu",
				"price": 0,
				"guestPrice": 9.5,
				"guestLimit": 2,
				"tickets": 50,
				"guestTickets": 20,
				"saleStart": "2022-02-10T11:30:00Z",
				"saleEnd": "2022-03-01T17:45:00Z",
				"dateTime": "2022-03-05T20:30:00Z",
				"groups": [11, 3]
			}`,
			[]int{11, 3},
			&model.Formal{
				Name:         "A formal",
				Menu:         "Some menu",
				Price:        0,
				GuestPrice:   9.5,
				GuestLimit:   2,
				Tickets:      50,
				GuestTickets: 20,
				SaleStart: time.Date(
					2022, 02, 10, 11, 30,
					0, 0, time.UTC,
				),
				SaleEnd: time.Date(
					2022, 03, 01, 17, 45,
					0, 0, time.UTC,
				),
				DateTime: time.Date(
					2022, 03, 05, 20, 30,
					0, 0, time.UTC,
				),
				Groups: []model.Group{
					{Model: model.Model{ID: 11}, Name: "Group 1"},
				},
			},
			&wants{http.StatusUnprocessableEntity, "Selected groups do not exist."},
		},
		{
			"Should Create",
			`{
				"name": "A formal",
				"menu": "Some menu",
				"price": 0,
				"guestPrice": 9.5,
				"guestLimit": 2,
				"tickets": 50,
				"guestTickets": 20,
				"saleStart": "2022-02-10T11:30:00Z",
				"saleEnd": "2022-03-01T17:45:00Z",
				"dateTime": "2022-03-05T20:30:00Z",
				"groups": [11, 3]
			}`,
			[]int{11, 3},
			&model.Formal{
				Name:         "A formal",
				Menu:         "Some menu",
				Price:        0,
				GuestPrice:   9.5,
				GuestLimit:   2,
				Tickets:      50,
				GuestTickets: 20,
				SaleStart: time.Date(
					2022, 02, 10, 11, 30,
					0, 0, time.UTC,
				),
				SaleEnd: time.Date(
					2022, 03, 01, 17, 45,
					0, 0, time.UTC,
				),
				DateTime: time.Date(
					2022, 03, 05, 20, 30,
					0, 0, time.UTC,
				),
				Groups: []model.Group{
					{Model: model.Model{ID: 11}, Name: "Group 1"},
					{Model: model.Model{ID: 3}, Name: "Group 2"},
				},
			},
			nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			e := echo.New()
			e.Validator = middleware.NewValidator()
			// HTTP
			req := httptest.NewRequest(
				http.MethodPost, "/formals", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			if test.formal != nil {
				s.formals.On("GetGroups", test.groups).Maybe().Once().Return(test.formal.Groups, nil)
			}
			if test.wants == nil {
				s.formals.On("Create", test.formal).Return(nil).Once()
			}

			// Test
			err := s.h.CreateFormal(c)
			if test.wants == nil {
				s.NoError(err)
				s.Equal(http.StatusCreated, rec.Code)
			} else {
				var he *echo.HTTPError
				if s.ErrorAs(err, &he) {
					s.Equal(test.wants.code, he.Code)
					s.Equal(test.wants.message, he.Message)
				}
			}
			s.formals.AssertExpectations(s.T())
		})
	}
}

func (s *AdminFormalSuite) TestUpdateFormal() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		body   string
		formal model.Formal
		wants  *wants
	}
	tests := []test{
		{
			"Should Update",
			`{
				"name": "Some formal",
				"menu": "Some menu",
				"price": 20,
				"guestPrice": 11.5,
				"guestLimit": 0,
				"tickets": 55,
				"guestTickets": 25,
				"saleStart": "2022-02-10T11:30:00Z",
				"saleEnd": "2022-03-01T17:45:00Z",
				"dateTime": "2022-03-05T20:30:00Z"
			}`,
			model.Formal{
				Model:        model.Model{ID: 34},
				Name:         "Some formal",
				Menu:         "Some menu",
				Price:        20,
				GuestPrice:   11.5,
				GuestLimit:   0,
				Tickets:      55,
				GuestTickets: 25,
				SaleStart: time.Date(
					2022, 02, 10, 11, 30,
					0, 0, time.UTC,
				),
				SaleEnd: time.Date(
					2022, 03, 01, 17, 45,
					0, 0, time.UTC,
				),
				DateTime: time.Date(
					2022, 03, 05, 20, 30,
					0, 0, time.UTC,
				),
			},
			nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			e := echo.New()
			e.Validator = middleware.NewValidator()
			// HTTP
			req := httptest.NewRequest(
				http.MethodPut, "/formals/34", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(test.formal.ID)))
			// Mock
			if test.wants == nil {
				s.formals.On("Update", &test.formal).Return(nil).Once()
			}

			// Test
			err := s.h.UpdateFormal(c)
			if test.wants == nil {
				s.NoError(err)
				s.Equal(http.StatusOK, rec.Code)
			} else {
				var he *echo.HTTPError
				if s.ErrorAs(err, &he) {
					s.Equal(test.wants.code, he.Code)
					s.Equal(test.wants.message, he.Message)
				}
			}
			s.formals.AssertExpectations(s.T())
		})
	}
}

func TestAdminFormalSuite(t *testing.T) {
	suite.Run(t, new(AdminFormalSuite))
}

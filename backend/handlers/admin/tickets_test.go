package admin_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	. "github.com/kcsu/store/handlers/admin"
	"github.com/kcsu/store/middleware"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminTicketSuite struct {
	suite.Suite
	h             *AdminHandler
	tickets       *mocks.TicketStore
	manualTickets *mocks.ManualTicketStore
}

func (s *AdminTicketSuite) SetupTest() {
	s.h = new(AdminHandler)
	s.tickets = new(mocks.TicketStore)
	s.manualTickets = new(mocks.ManualTicketStore)
	s.h.Tickets = s.tickets
	s.h.ManualTickets = s.manualTickets
}

func (s *AdminTicketSuite) TestCancelTicket() {
	type wants struct {
		status  int
		message string
	}
	type test struct {
		name   string
		id     int
		ticket *model.Ticket
		wants  *wants
	}
	tests := []test{
		{
			"Ticket Not Found",
			17,
			nil,
			&wants{http.StatusNotFound, "Not Found"},
		},
		{
			"Should Cancel Guest",
			13,
			&model.Ticket{
				Model:      model.Model{ID: 13},
				IsGuest:    true,
				FormalID:   1,
				UserID:     3,
				MealOption: "Vegetarian",
			},
			nil,
		},
		{
			"Should Cancel All",
			26,
			&model.Ticket{
				Model:      model.Model{ID: 26},
				IsGuest:    false,
				FormalID:   9,
				UserID:     5,
				MealOption: "Vegan",
			},
			nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			// HTTP
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(test.id))
			// Mock
			if test.ticket != nil {
				s.tickets.On("Find", test.id).Return(*test.ticket, nil).Once()
				if test.ticket.IsGuest {
					s.tickets.On("Delete", test.id).Return(nil).Once()
				} else {
					s.tickets.On("DeleteByFormal", test.ticket.FormalID, test.ticket.UserID).Return(nil).Once()
				}
			} else {
				s.tickets.On("Find", test.id).Return(model.Ticket{}, gorm.ErrRecordNotFound).Once()
			}
			// Test
			err := s.h.CancelTicket(c)
			if test.wants == nil {
				s.NoError(err)
				s.Equal(http.StatusOK, rec.Code)
			} else {
				var he *echo.HTTPError
				if s.ErrorAs(err, &he) {
					s.Equal(test.wants.status, he.Code)
					s.Equal(test.wants.message, he.Message)
				}
			}
		})
	}
	s.tickets.AssertExpectations(s.T())
}

func (s *AdminTicketSuite) TestEditTicket() {
	e := echo.New()
	e.Validator = middleware.NewValidator()
	body := `{
		"option": "Vegan"
	}`
	req := httptest.NewRequest(
		http.MethodPut, "/tickets/7", strings.NewReader(body),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("7")
	s.tickets.On(
		"Update", 7, &dto.TicketRequestDto{MealOption: "Vegan"},
	).Return(nil).Once()
	err := s.h.EditTicket(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.tickets.AssertExpectations(s.T())
}

func (s *AdminTicketSuite) TestCreateManualTicket() {
	type wants struct {
		status  int
		message string
	}
	type test struct {
		name   string
		body   string
		ticket *model.ManualTicket
		wants  *wants
	}
	tests := []test{
		{
			"Should Create Ticket",
			`{
				"option": "Vegan",
				"formalId": 1,
				"type": "guest",
				"name": "John Doe",
				"justification": "Freebie",
				"email": "jd123@cam.ac.uk"
			}`,
			&model.ManualTicket{
				FormalID:      1,
				MealOption:    "Vegan",
				Type:          "guest",
				Name:          "John Doe",
				Justification: "Freebie",
				Email:         "jd123@cam.ac.uk",
			},
			nil,
		},
		{
			"Invalid Type",
			`{
				"option": "Vegan",
				"formalId": 1,
				"type": "invalid",
				"name": "John Doe",
				"justification": "Freebie",
				"email": "jd123@cam.ac.uk"
			}`,
			nil,
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'CreateManualTicketDto.Type' Error:Field validation for 'Type' failed on the 'oneof' tag",
			},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			// HTTP
			e := echo.New()
			e.Validator = middleware.NewValidator()
			req := httptest.NewRequest(http.MethodPost, "/tickets/manual", strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			if test.ticket != nil {
				s.manualTickets.On("Create", test.ticket).Return(nil).Once()
			}
			// Test
			err := s.h.CreateManualTicket(c)
			if test.wants == nil {
				s.NoError(err)
				s.Equal(http.StatusCreated, rec.Code)
			} else {
				var he *echo.HTTPError
				if s.ErrorAs(err, &he) {
					s.Equal(test.wants.status, he.Code)
					s.Equal(test.wants.message, he.Message)
				}
			}
		})
	}
	s.manualTickets.AssertExpectations(s.T())
}

func (s *AdminTicketSuite) TestCancelManualTicket() {
	type wants struct {
		status  int
		message string
	}
	type test struct {
		name   string
		id     int
		ticket *model.ManualTicket
		wants  *wants
	}
	tests := []test{
		{
			"Should Cancel Ticket",
			7,
			&model.ManualTicket{
				Model: model.Model{ID: 7},
			},
			nil,
		},
		{
			"Should Return Not Found",
			7,
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			// HTTP
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/tickets/manual/7", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Language: go
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(test.id))
			// Mock
			if test.ticket != nil {
				s.manualTickets.On("Find", test.id).Return(*test.ticket, nil).Once()
				s.manualTickets.On("Delete", test.id).Return(nil).Once()
			} else {
				s.manualTickets.On("Find", test.id).Return(model.ManualTicket{}, gorm.ErrRecordNotFound).Once()
			}
			// Test
			err := s.h.CancelManualTicket(c)
			if test.wants == nil {
				s.NoError(err)
				s.Equal(http.StatusOK, rec.Code)
			} else {
				var he *echo.HTTPError
				if s.ErrorAs(err, &he) {
					s.Equal(test.wants.status, he.Code)
					s.Equal(test.wants.message, he.Message)
				}
			}
		})
	}
	s.manualTickets.AssertExpectations(s.T())
}

func (s *AdminTicketSuite) TestEditManualTicket() {
	type wants struct {
		status  int
		message string
	}
	type test struct {
		name   string
		id     int
		body   string
		ticket *model.ManualTicket
		wants  *wants
	}
	tests := []test{
		{
			"Should Edit Ticket",
			7,
			`{
				"option": "Vegan",
				"type": "guest",
				"name": "John Doe",
				"justification": "Freebie",
				"email": "jd1234@cam.ac.uk"
			}`,
			&model.ManualTicket{
				Model:         model.Model{ID: 7},
				MealOption:    "Vegan",
				FormalID:      3,
				Type:          "guest",
				Name:          "John Doe",
				Justification: "Freebie",
				Email:         "jd1234@cam.ac.uk",
			},
			nil,
		},
		{
			"Should Return Not Found",
			7,
			`{
				"option": "Vegan",
				"type": "guest",
				"name": "John Doe",
				"justification": "Freebie",
				"email": "jd1234@cam.ac.uk"
			}`,
			nil,
			&wants{http.StatusNotFound, "Not Found"},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			// HTTP
			e := echo.New()
			e.Validator = middleware.NewValidator()
			req := httptest.NewRequest(http.MethodPut, "/tickets/manual/7", strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(test.id))
			// Mock
			if test.ticket != nil {
				tr := model.ManualTicket{
					FormalID: test.ticket.FormalID,
				}
				tr.ID = test.ticket.ID
				s.manualTickets.On("Find", test.id).Return(tr, nil).Once()
				s.manualTickets.On("Update", test.ticket).Return(nil).Once()
			} else {
				s.manualTickets.On("Find", test.id).Return(model.ManualTicket{}, gorm.ErrRecordNotFound).Once()
			}
			// Test
			err := s.h.EditManualTicket(c)
			if test.wants == nil {
				s.NoError(err)
				s.Equal(http.StatusOK, rec.Code)
			} else {
				var he *echo.HTTPError
				if s.ErrorAs(err, &he) {
					s.Equal(test.wants.status, he.Code)
					s.Equal(test.wants.message, he.Message)
				}
			}
		})
	}
	s.manualTickets.AssertExpectations(s.T())
}

func TestAdminTicketSuite(t *testing.T) {
	suite.Run(t, new(AdminTicketSuite))
}

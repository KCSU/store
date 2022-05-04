package admin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	. "github.com/kcsu/store/handlers/admin"
	"github.com/kcsu/store/middleware"
	mocks "github.com/kcsu/store/mocks/db"
	mm "github.com/kcsu/store/mocks/middleware"
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminTicketSuite struct {
	suite.Suite
	h             *AdminHandler
	tickets       *mocks.TicketStore
	manualTickets *mocks.ManualTicketStore
	formals       *mocks.FormalStore
}

func (s *AdminTicketSuite) SetupTest() {
	s.h = new(AdminHandler)
	s.tickets = mocks.NewTicketStore(s.T())
	s.manualTickets = mocks.NewManualTicketStore(s.T())
	s.formals = mocks.NewFormalStore(s.T())
	s.h.Tickets = s.tickets
	s.h.ManualTickets = s.manualTickets
	s.h.Formals = s.formals
	// FIXME: We currently ignore calls to Access.Log
	// but this is probably a bad idea, especially here.
	accessMock := mm.NewAccess(s.T())
	accessMock.On(
		"Log",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Maybe().Return(nil)
	s.h.Access = accessMock
}

func (s *AdminTicketSuite) TestCancelTicket() {
	type wants struct {
		status  int
		message string
	}
	type test struct {
		name   string
		id     uuid.UUID
		ticket *model.Ticket
		wants  *wants
	}
	tests := []test{
		{
			"Ticket Not Found",
			uuid.New(),
			nil,
			&wants{http.StatusNotFound, "Not Found"},
		},
		{
			"Should Cancel Guest",
			uuid.MustParse("c0187157-a065-41b4-a29e-018f49934701"),
			&model.Ticket{
				Model:      model.Model{ID: uuid.MustParse("c0187157-a065-41b4-a29e-018f49934701")},
				IsGuest:    true,
				FormalID:   uuid.New(),
				UserID:     uuid.New(),
				Formal:     &model.Formal{},
				User:       &model.User{},
				MealOption: "Vegetarian",
			},
			nil,
		},
		{
			"Should Cancel All",
			uuid.MustParse("57ae970f-7f82-471b-a60c-46d4269f2e81"),
			&model.Ticket{
				Model:      model.Model{ID: uuid.MustParse("57ae970f-7f82-471b-a60c-46d4269f2e81")},
				IsGuest:    false,
				FormalID:   uuid.New(),
				UserID:     uuid.New(),
				Formal:     &model.Formal{},
				User:       &model.User{},
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
			c.SetParamValues(test.id.String())
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
}

func (s *AdminTicketSuite) TestEditTicket() {
	e := echo.New()
	e.Validator = middleware.NewValidator()
	body := `{
		"option": "Vegan"
	}`
	id := uuid.New()
	req := httptest.NewRequest(
		http.MethodPut, fmt.Sprint("/tickets/", id), strings.NewReader(body),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())
	s.tickets.On(
		"Update", id, &dto.TicketRequestDto{MealOption: "Vegan"},
	).Return(nil).Once()
	err := s.h.EditTicket(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
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
				"formalId": "eb6d25f9-14e8-4abe-90e2-80542c73a2b6",
				"type": "guest",
				"name": "John Doe",
				"justification": "Freebie",
				"email": "jd123@cam.ac.uk"
			}`,
			&model.ManualTicket{
				FormalID:      uuid.MustParse("eb6d25f9-14e8-4abe-90e2-80542c73a2b6"),
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
				"formalId": "9739dbf9-4ce3-44e5-b9e1-cdfb0ac8de14",
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
				s.formals.On("Find", mock.AnythingOfType("uuid.UUID")).Return(
					model.Formal{}, nil,
				).Once()
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
}

func (s *AdminTicketSuite) TestCancelManualTicket() {
	type wants struct {
		status  int
		message string
	}
	type test struct {
		name   string
		id     uuid.UUID
		ticket *model.ManualTicket
		wants  *wants
	}
	tests := []test{
		{
			"Should Cancel Ticket",
			uuid.MustParse("06a1832e-4231-4ccb-bc80-e41680bc5a2b"),
			&model.ManualTicket{
				Model: model.Model{ID: uuid.MustParse("06a1832e-4231-4ccb-bc80-e41680bc5a2b")},
			},
			nil,
		},
		{
			"Should Return Not Found",
			uuid.New(),
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
			req := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprint("/tickets/manual/", test.id),
				nil,
			)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Language: go
			c.SetParamNames("id")
			c.SetParamValues(test.id.String())
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
}

func (s *AdminTicketSuite) TestEditManualTicket() {
	type wants struct {
		status  int
		message string
	}
	type test struct {
		name   string
		id     uuid.UUID
		body   string
		ticket *model.ManualTicket
		wants  *wants
	}
	tests := []test{
		{
			"Should Edit Ticket",
			uuid.MustParse("0aa3b271-d66e-4b3b-af6e-e84ead601171"),
			`{
				"option": "Vegan",
				"type": "guest",
				"name": "John Doe",
				"justification": "Freebie",
				"email": "jd1234@cam.ac.uk"
			}`,
			&model.ManualTicket{
				Model:         model.Model{ID: uuid.MustParse("0aa3b271-d66e-4b3b-af6e-e84ead601171")},
				MealOption:    "Vegan",
				FormalID:      uuid.New(),
				Type:          "guest",
				Name:          "John Doe",
				Justification: "Freebie",
				Email:         "jd1234@cam.ac.uk",
			},
			nil,
		},
		{
			"Should Return Not Found",
			uuid.MustParse("fcfe1c65-36be-434c-b5ca-8644f971ac91"),
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
			req := httptest.NewRequest(
				http.MethodPut,
				fmt.Sprint("/tickets/manual/", test.id),
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.id.String())
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
}

func TestAdminTicketSuite(t *testing.T) {
	suite.Run(t, new(AdminTicketSuite))
}

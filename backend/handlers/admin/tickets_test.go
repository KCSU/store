package admin_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	. "github.com/kcsu/store/handlers/admin"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminTicketSuite struct {
	suite.Suite
	h       *AdminHandler
	tickets *mocks.TicketStore
}

func (s *AdminTicketSuite) SetupTest() {
	s.h = new(AdminHandler)
	s.tickets = new(mocks.TicketStore)
	s.h.Tickets = s.tickets
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

func TestAdminTicketSuite(t *testing.T) {
	suite.Run(t, new(AdminTicketSuite))
}

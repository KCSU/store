package admin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/kcsu/store/handlers/admin"
	"github.com/kcsu/store/middleware"
	mocks "github.com/kcsu/store/mocks/db"
	mm "github.com/kcsu/store/mocks/middleware"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminFormalSuite struct {
	suite.Suite
	h       *AdminHandler
	formals *mocks.FormalStore
}

func (s *AdminFormalSuite) SetupTest() {
	// Init handler
	s.h = new(AdminHandler)
	s.formals = mocks.NewFormalStore(s.T())
	s.h.Formals = s.formals
	// HACK: We currently ignore calls to Access.Log
	// but this is probably a bad idea.
	accessMock := mm.NewAccess(s.T())
	accessMock.On(
		"Log",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Maybe().Return(nil)
	s.h.Access = accessMock
}

func (s *AdminFormalSuite) TestGetFormals() {
	const expectedJSON = `[
		{
			"id": "a49b7a96-92f3-49a1-88bb-e99c3549fc25",
			"createdAt": "0001-01-01T00:00:00Z",
			"updatedAt": "0001-01-01T00:00:00Z",
			"deletedAt": null,
			"name": "Test 1",
			"menu": "A menu",
			"price": 21.3,
			"guestPrice": 11.6,
			"guestLimit": 0,
			"firstSaleTickets": 0,
			"firstSaleGuestTickets": 0,
			"firstSaleStart": "0001-01-01T00:00:00Z",
			"secondSaleTickets": 0,
			"secondSaleGuestTickets": 0,
			"secondSaleStart": "0001-01-01T00:00:00Z",
			"hasGuestList": true,
			"isVisible": false,
			"saleEnd": "0001-01-01T00:00:00Z",
			"dateTime": "0001-01-01T00:00:00Z",
			"ticketsRemaining": 24,
			"guestTicketsRemaining": 56,
			"groups": []
		},
		{
			"id": "9acf368c-0c7b-4ef2-85e4-d34bfae8dd2e",
			"createdAt": "0001-01-01T00:00:00Z",
			"updatedAt": "0001-01-01T00:00:00Z",
			"deletedAt": null,
			"name": "Test 2",
			"menu": "Another menu",
			"price": 15.6,
			"guestPrice": 27.2,
			"guestLimit": 0,
			"firstSaleTickets": 0,
			"firstSaleGuestTickets": 0,
			"firstSaleStart": "0001-01-01T00:00:00Z",
			"secondSaleTickets": 0,
			"secondSaleGuestTickets": 0,
			"secondSaleStart": "0001-01-01T00:00:00Z",
			"hasGuestList": false,
			"isVisible": true,
			"saleEnd": "0001-01-01T00:00:00Z",
			"dateTime": "0001-01-01T00:00:00Z",
			"ticketsRemaining": 64,
			"guestTicketsRemaining": 31,
			"groups": [
				{
					"id": "56b7ebef-23a9-47d5-88f3-aacc9898807d",
					"name": "Group A"
				},
				{
					"id": "416e9375-8b09-4b04-a459-436b078dd375",
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
			Model:        model.Model{ID: uuid.MustParse("a49b7a96-92f3-49a1-88bb-e99c3549fc25")},
			Name:         "Test 1",
			Menu:         "A menu",
			Price:        21.3,
			GuestPrice:   11.6,
			HasGuestList: true,
			IsVisible:    false,
		},
		{
			Model:        model.Model{ID: uuid.MustParse("9acf368c-0c7b-4ef2-85e4-d34bfae8dd2e")},
			Name:         "Test 2",
			Menu:         "Another menu",
			Price:        15.6,
			GuestPrice:   27.2,
			HasGuestList: false,
			IsVisible:    true,
			Groups: []model.Group{
				{
					Model: model.Model{ID: uuid.MustParse("56b7ebef-23a9-47d5-88f3-aacc9898807d")},
					Name:  "Group A",
				},
				{
					Model: model.Model{ID: uuid.MustParse("416e9375-8b09-4b04-a459-436b078dd375")},
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
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminFormalSuite) TestGetFormal() {
	const expectedJSON = `{
		"id": "c5212510-4fb2-4623-8117-9dca85ed3ea2",
		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z",
		"deletedAt": null,
		"name": "Test 5",
		"menu": "Another menu",
		"price": 26.3,
		"guestPrice": 12.7,
		"guestLimit": 3,
		"firstSaleTickets": 120,
		"firstSaleGuestTickets": 0,
		"firstSaleStart": "0001-01-01T00:00:00Z",
		"secondSaleTickets": 0,
		"secondSaleGuestTickets": 0,
		"secondSaleStart": "0001-01-01T00:00:00Z",
		"hasGuestList": true,
		"isVisible": true,
		"saleEnd": "0001-01-01T00:00:00Z",
		"dateTime": "0001-01-01T00:00:00Z",
		"queueLength": 11,
		"groups": [{
			"id": "97e9db3b-077a-4150-bc20-66b5a8490083",
			"name": "Group"
		}],
		"ticketSales": [{
			"id": "af7960bf-dca3-45e5-869e-5637db289e5e",
			"createdAt": "0001-01-01T00:00:00Z",
			"updatedAt": "0001-01-01T00:00:00Z",
			"deletedAt": null,
			"isGuest": false,
			"isQueue": false,
			"formalId": "c5212510-4fb2-4623-8117-9dca85ed3ea2",
			"option": "Vegetarian",
			"userId": "809331d0-f3ad-4cc4-bbc8-5149737ffca4",
			"userName": "Stephen Sondheim",
			"userEmail": "ss103@cam.ac.uk"
		}],
		"manualTickets": [{
			"id": "3cf7795b-858b-40d7-9541-91e46e433215",
			"createdAt": "0001-01-01T00:00:00Z",
			"updatedAt": "0001-01-01T00:00:00Z",
			"deletedAt": null,
			"formalId": "c5212510-4fb2-4623-8117-9dca85ed3ea2",
			"option": "Vegan",
			"type": "standard",
			"name": "Kara Thrace",
			"justification": "Ents Committee",
			"email": "kth123@cam.ac.uk"
		}]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/formals/c5212510-4fb2-4623-8117-9dca85ed3ea2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("c5212510-4fb2-4623-8117-9dca85ed3ea2")
	formal := model.Formal{
		Model:        model.Model{ID: uuid.MustParse("c5212510-4fb2-4623-8117-9dca85ed3ea2")},
		Name:         "Test 5",
		Menu:         "Another menu",
		Price:        26.3,
		GuestPrice:   12.7,
		GuestLimit:   3,
		HasGuestList: true,
		IsVisible:    true,
		Groups: []model.Group{{
			Model: model.Model{ID: uuid.MustParse("97e9db3b-077a-4150-bc20-66b5a8490083")},
			Name:  "Group",
		}},
		TicketSales: []model.Ticket{{
			Model: model.Model{ID: uuid.MustParse("af7960bf-dca3-45e5-869e-5637db289e5e")},
			User: &model.User{
				Model: model.Model{ID: uuid.MustParse("809331d0-f3ad-4cc4-bbc8-5149737ffca4")},
				Name:  "Stephen Sondheim",
				Email: "ss103@cam.ac.uk",
			},
			UserID:     uuid.MustParse("809331d0-f3ad-4cc4-bbc8-5149737ffca4"),
			MealOption: "Vegetarian",
			FormalID:   uuid.MustParse("c5212510-4fb2-4623-8117-9dca85ed3ea2"),
		}},
		ManualTickets: []model.ManualTicket{{
			Model:         model.Model{ID: uuid.MustParse("3cf7795b-858b-40d7-9541-91e46e433215")},
			MealOption:    "Vegan",
			FormalID:      uuid.MustParse("c5212510-4fb2-4623-8117-9dca85ed3ea2"),
			Type:          "standard",
			Name:          "Kara Thrace",
			Justification: "Ents Committee",
			Email:         "kth123@cam.ac.uk",
		}},
		FirstSaleTickets: 120,
	}
	s.formals.On("FindWithTickets", formal.ID).Return(formal, nil)
	s.formals.On("GetQueueLength", formal.ID).Return(11, nil)
	err := s.h.GetFormal(c)
	s.NoError(err)
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
		groups []uuid.UUID
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
				"firstSaleTickets": 50,
				"firstSaleGuestTickets": 20,
				"firstSaleStart": "2022-02-10T11:30:00Z",
				"secondSaleStart": "2022-02-10T11:30:00Z",
				"saleEnd": "2022-03-01T17:45:00Z",
				"dateTime": "2022-03-05T20:30:00Z",
				"hasGuestList": true,
				"isVisible": false,
				"groups": [
					"d9577854-0f8b-4350-ae42-4a5572913444",
					"609f4ef2-d516-4281-abb1-98fc687fd991"
				]
			}`,
			[]uuid.UUID{
				uuid.MustParse("d9577854-0f8b-4350-ae42-4a5572913444"),
				uuid.MustParse("609f4ef2-d516-4281-abb1-98fc687fd991"),
			},
			&model.Formal{
				Name:                  "A formal",
				Menu:                  "Some menu",
				Price:                 0,
				GuestPrice:            9.5,
				GuestLimit:            2,
				FirstSaleTickets:      50,
				FirstSaleGuestTickets: 20,
				HasGuestList:          true,
				IsVisible:             false,
				FirstSaleStart: time.Date(
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
				Groups: []model.Group{{
					Model: model.Model{ID: uuid.MustParse("609f4ef2-d516-4281-abb1-98fc687fd991")},
					Name:  "Group 1",
				}},
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
				"firstSaleTickets": 50,
				"firstSaleGuestTickets": 20,
				"firstSaleStart": "2022-02-10T11:30:00Z",
				"secondSaleTickets": 10,
				"secondSaleGuestTickets": 5,
				"secondSaleStart": "2022-02-11T11:30:00Z",
				"hasGuestList": true,
				"isVisible": true,
				"saleEnd": "2022-03-01T17:45:00Z",
				"dateTime": "2022-03-05T20:30:00Z",
				"groups": [
					"d9577854-0f8b-4350-ae42-4a5572913444",
					"609f4ef2-d516-4281-abb1-98fc687fd991"
				]
			}`,
			[]uuid.UUID{
				uuid.MustParse("d9577854-0f8b-4350-ae42-4a5572913444"),
				uuid.MustParse("609f4ef2-d516-4281-abb1-98fc687fd991"),
			},
			&model.Formal{
				Name:                  "A formal",
				Menu:                  "Some menu",
				Price:                 0,
				GuestPrice:            9.5,
				GuestLimit:            2,
				FirstSaleTickets:      50,
				FirstSaleGuestTickets: 20,
				HasGuestList:          true,
				IsVisible:             true,
				FirstSaleStart: time.Date(
					2022, 02, 10, 11, 30,
					0, 0, time.UTC,
				),
				SecondSaleTickets:      10,
				SecondSaleGuestTickets: 5,
				SecondSaleStart: time.Date(
					2022, 02, 11, 11, 30,
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
					{
						Model: model.Model{
							ID: uuid.MustParse("d9577854-0f8b-4350-ae42-4a5572913444"),
						},
						Name: "Group 1",
					},
					{
						Model: model.Model{
							ID: uuid.MustParse("609f4ef2-d516-4281-abb1-98fc687fd991"),
						},
						Name: "Group 2",
					},
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
				"firstSaleTickets": 55,
				"firstSaleGuestTickets": 25,
				"firstSaleStart": "2022-02-10T11:30:00Z",
				"secondSaleTickets": 10,
				"secondSaleGuestTickets": 5,
				"secondSaleStart": "2022-02-11T11:30:00Z",
				"hasGuestList": true,
				"isVisible": false,
				"saleEnd": "2022-03-01T17:45:00Z",
				"dateTime": "2022-03-05T20:30:00Z"
			}`,
			model.Formal{
				Model:                 model.Model{ID: uuid.MustParse("3008f056-8e1e-4971-af88-27e9146da1ae")},
				Name:                  "Some formal",
				Menu:                  "Some menu",
				Price:                 20,
				GuestPrice:            11.5,
				GuestLimit:            0,
				FirstSaleTickets:      55,
				FirstSaleGuestTickets: 25,
				HasGuestList:          true,
				IsVisible:             false,
				FirstSaleStart: time.Date(
					2022, 02, 10, 11, 30,
					0, 0, time.UTC,
				),
				SecondSaleTickets:      10,
				SecondSaleGuestTickets: 5,
				SecondSaleStart: time.Date(
					2022, 02, 11, 11, 30,
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
				http.MethodPut, "/formals/3008f056-8e1e-4971-af88-27e9146da1ae", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.formal.ID.String())
			// Mock
			if test.wants == nil {
				s.formals.On("Find", test.formal.ID).Once().Return(model.Formal{}, nil)
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
		})
	}
}

func (s *AdminFormalSuite) TestDeleteFormal() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		id     uuid.UUID
		formal model.Formal
		wants  *wants
	}
	tests := []test{
		{
			"Should Delete",
			uuid.MustParse("08d84aba-ce1b-49b7-8946-eca7e4f95aeb"),
			model.Formal{
				Model: model.Model{ID: uuid.MustParse("08d84aba-ce1b-49b7-8946-eca7e4f95aeb")},
				Name:  "Some formal",
			},
			nil,
		},
		{
			"Formal Not Found",
			uuid.MustParse("25bebb0e-a5ab-4344-9b9d-349bda25f669"),
			model.Formal{},
			&wants{http.StatusNotFound, "Not Found"},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			e := echo.New()
			e.Validator = middleware.NewValidator()
			// HTTP
			req := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/formals/%s", test.id.String()),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.id.String())
			// Mock
			if test.wants == nil {
				s.formals.On("Find", test.id).Return(
					test.formal, nil,
				).Once()
				s.formals.On("Delete", &test.formal).Return(nil).Once()
			} else {
				s.formals.On("Find", test.id).Return(
					model.Formal{}, gorm.ErrRecordNotFound,
				).Once()
			}

			// Test
			err := s.h.DeleteFormal(c)
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
		})
	}
}

func (s *AdminFormalSuite) TestUpdateFormalGroups() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		body   string
		ids    []uuid.UUID
		formal model.Formal
		groups []model.Group
		wants  *wants
	}
	id := uuid.MustParse("bc6a54fa-0ccb-48f2-9ba7-3ab7e15e053d")
	tests := []test{
		{
			"Invalid Groups",
			`[
				"ada9b70c-be1d-4ad4-a852-337c2df26184",
				"6a0cf98f-05ea-403d-a46d-3ca59c89b4a9",
				"5cec7779-c8cb-482d-8deb-235c537c31da"
			]`,
			[]uuid.UUID{
				uuid.MustParse("ada9b70c-be1d-4ad4-a852-337c2df26184"),
				uuid.MustParse("6a0cf98f-05ea-403d-a46d-3ca59c89b4a9"),
				uuid.MustParse("5cec7779-c8cb-482d-8deb-235c537c31da"),
			},
			model.Formal{
				Model: model.Model{ID: id},
				Name:  "My Formal",
			},
			[]model.Group{{
				Model: model.Model{ID: uuid.MustParse("6a0cf98f-05ea-403d-a46d-3ca59c89b4a9")},
				Name:  "A Group",
			}},
			&wants{
				http.StatusUnprocessableEntity,
				"Selected groups do not exist.",
			},
		},
		{
			"Should Update",
			`[
				"ada9b70c-be1d-4ad4-a852-337c2df26184",
				"6a0cf98f-05ea-403d-a46d-3ca59c89b4a9"
			]`,
			[]uuid.UUID{
				uuid.MustParse("ada9b70c-be1d-4ad4-a852-337c2df26184"),
				uuid.MustParse("6a0cf98f-05ea-403d-a46d-3ca59c89b4a9"),
			},
			model.Formal{
				Model: model.Model{ID: id},
				Name:  "My Formal",
			},
			[]model.Group{
				{
					Model: model.Model{ID: uuid.MustParse("ada9b70c-be1d-4ad4-a852-337c2df26184")},
					Name:  "A Group",
				},
				{
					Model: model.Model{ID: uuid.MustParse("6a0cf98f-05ea-403d-a46d-3ca59c89b4a9")},
					Name:  "B Group",
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
				http.MethodPut, "/formals/bc6a54fa-0ccb-48f2-9ba7-3ab7e15e053d/groups", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.formal.ID.String())

			s.formals.On("FindWithGroups", test.formal.ID).Return(test.formal, nil).Once()
			s.formals.On("GetGroups", test.ids).Return(test.groups, nil).Once()
			// Mock
			if test.wants == nil {
				s.formals.On(
					"UpdateGroups", test.formal, test.groups,
				).Return(nil).Once()
			}

			// Test
			err := s.h.UpdateFormalGroups(c)
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
		})
	}
}

func TestAdminFormalSuite(t *testing.T) {
	suite.Run(t, new(AdminFormalSuite))
}

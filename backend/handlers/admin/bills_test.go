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
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminBillSuite struct {
	suite.Suite
	h     *AdminHandler
	bills *mocks.BillStore
}

func (s *AdminBillSuite) SetupTest() {
	s.h = new(AdminHandler)
	s.bills = mocks.NewBillStore(s.T())
	s.h.Bills = s.bills
}

func (s *AdminBillSuite) TestGetBills() {
	const expectedJSON = `[
		{
			"id": "580cf6b7-b799-4199-bb0e-1e281caa6945",
			"createdAt": "0001-01-01T00:00:00Z",
			"updatedAt": "0001-01-01T00:00:00Z",
			"deletedAt": null,
			"name": "Bill 1",
			"start": "2022-01-22T00:00:00Z",
			"end": "2022-05-22T00:00:00+01:00",
			"formals": null
		}
	]`
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock database
	bills := []model.Bill{{
		Model: model.Model{ID: uuid.MustParse("580cf6b7-b799-4199-bb0e-1e281caa6945")},
		Name:  "Bill 1",
		Start: time.Date(2022, 1, 22, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 5, 22, 0, 0, 0, 0, time.FixedZone("Europe/London", 3600)),
	}}
	s.bills.On("Get").Return(bills, nil)
	// Run test
	err := s.h.GetBills(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminBillSuite) TestGetBill() {
	const expectedJSON = `{
  		"id": "af63ca4e-7f54-45bf-aab0-2d971f08222a",
  		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z",
		"deletedAt": null,
		"name": "Bill 1",
		"start": "2022-01-22T00:00:00Z",
		"end": "2022-05-22T00:00:00+01:00",
		"formals": [
			{
				"id": "92851ad0-d9ea-4223-970c-496d260b9905",
				"createdAt": "0001-01-01T00:00:00Z",
				"updatedAt": "0001-01-01T00:00:00Z",
				"deletedAt": null,
				"name": "Test Formal",
				"menu": "This is a menu",
				"price": 15,
				"guestPrice": 30,
				"guestLimit": 2,
				"tickets": 10,
				"guestTickets": 5,
				"saleStart": "2022-04-23T14:00:00+01:00",
				"saleEnd": "2022-07-21T22:30:00+01:00",
				"dateTime": "2022-08-17T15:30:00+01:00",
				"billId": "af63ca4e-7f54-45bf-aab0-2d971f08222a"
			}
		]
	}`
	// Init HTTP
	id := uuid.MustParse("af63ca4e-7f54-45bf-aab0-2d971f08222a")
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/bills/", id), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())
	// Mock database
	bill := model.Bill{
		Model: model.Model{ID: id},
		Name:  "Bill 1",
		Start: time.Date(2022, 1, 22, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 5, 22, 0, 0, 0, 0, time.FixedZone("Europe/London", 3600)),
		Formals: []model.Formal{
			{
				Model:        model.Model{ID: uuid.MustParse("92851ad0-d9ea-4223-970c-496d260b9905")},
				Name:         "Test Formal",
				Menu:         "This is a menu",
				Price:        15,
				GuestPrice:   30,
				GuestLimit:   2,
				Tickets:      10,
				GuestTickets: 5,
				SaleStart:    time.Date(2022, 4, 23, 14, 0, 0, 0, time.FixedZone("Europe/London", 3600)),
				SaleEnd:      time.Date(2022, 7, 21, 22, 30, 0, 0, time.FixedZone("Europe/London", 3600)),
				DateTime:     time.Date(2022, 8, 17, 15, 30, 0, 0, time.FixedZone("Europe/London", 3600)),
				BillID:       &id,
			},
		},
	}
	s.bills.On("FindWithFormals", id).Return(bill, nil)
	// Run test
	err := s.h.GetBill(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminBillSuite) TestCreateBill() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name  string
		body  string
		bill  *model.Bill
		wants *wants
	}
	tests := []test{
		{
			"Should Create",
			`{
				"name": "Bill 1",
				"start": "2022-01-22",
				"end": "2022-05-22"
			}`,
			&model.Bill{
				Name:  "Bill 1",
				Start: time.Date(2022, 1, 22, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, 5, 22, 0, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"Invalid End",
			`{
				"name": "Bill 1",
				"start": "2022-01-22",
				"end": "09-05-2023"
			}`,
			nil,
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'BillDto.End' Error:Field validation for 'End' failed on the 'datetime' tag",
			},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			// Init HTTP
			e := echo.New()
			e.Validator = middleware.NewValidator()
			req := httptest.NewRequest(http.MethodPost, "/bills", strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock database
			if test.bill != nil {
				s.bills.On("Create", test.bill).Return(nil)
			}
			// Run test
			err := s.h.CreateBill(c)
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

func (s *AdminBillSuite) TestUpdateBill() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name  string
		body  string
		valid bool
		bill  *model.Bill
		wants *wants
	}
	id := uuid.New()
	tests := []test{
		{
			"Should Update",
			`{
				"name": "Bill 1",
				"start": "2022-01-22",
				"end": "2022-05-22"
			}`,
			true,
			&model.Bill{
				Model: model.Model{ID: id},
				Name:  "Bill 1",
				Start: time.Date(2022, 1, 22, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, 5, 22, 0, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"Not found",
			`{
				"name": "Bill 1",
				"start": "2022-01-22",
				"end": "2022-05-22"
			}`,
			true,
			nil,
			&wants{http.StatusNotFound, "Not Found"},
		},
		{
			"Invalid start",
			`{
				"name": "Bill 1",
				"start": "2022-01-22T00:00:00Z",
				"end": "2022-05-22"
			}`,
			false,
			nil,
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'BillDto.Start' Error:Field validation for 'Start' failed on the 'datetime' tag",
			},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			// Init HTTP
			e := echo.New()
			e.Validator = middleware.NewValidator()
			req := httptest.NewRequest(
				http.MethodPut,
				fmt.Sprint("/bills/", id),
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(id.String())
			// Mock database
			if test.bill != nil {
				s.bills.On("Find", id).Return(
					model.Bill{
						Model: model.Model{ID: id},
						Name:  "Test",
					},
					nil,
				).Once()
				s.bills.On("Update", test.bill).Return(nil).Once()
			} else if test.valid {
				s.bills.On("Find", id).Return(model.Bill{}, gorm.ErrRecordNotFound).Once()
			}
			// Run test
			err := s.h.UpdateBill(c)
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

func (s *AdminBillSuite) TestDeleteBill() {
	e := echo.New()
	id := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: id},
		Name:  "Test",
	}
	route := fmt.Sprint("/bills/", id)
	req := httptest.NewRequest(http.MethodDelete, route, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())
	// Mock
	s.bills.On("Find", id).Return(bill, nil).Once()
	s.bills.On("Delete", &bill).Return(nil).Once()
	// Test
	err := s.h.DeleteBill(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *AdminBillSuite) TestAddBillFormals() {
	fid1, fid2 := uuid.New(), uuid.New()
	body := fmt.Sprintf(`{"formalIds": ["%s", "%s"]}`, fid1.String(), fid2.String())
	billId := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: billId},
		Name:  "Test",
	}
	e := echo.New()
	e.Validator = middleware.NewValidator()
	req := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprint("/bills/", "/formals"),
		strings.NewReader(body),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(billId.String())
	// Mock database
	s.bills.On("Find", billId).Return(bill, nil).Once()
	s.bills.On("AddFormals", &bill, []uuid.UUID{fid1, fid2}).Return(nil).Once()
	// Run test
	err := s.h.AddBillFormals(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *AdminBillSuite) TestRemoveBillFormal() {
	formalId := uuid.New()
	billId := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: billId},
		Name:  "Test",
	}
	e := echo.New()
	e.Validator = middleware.NewValidator()
	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprint("/bills/", billId.String(), "/formals/", formalId.String()),
		nil,
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id", "formalId")
	c.SetParamValues(billId.String(), formalId.String())
	// Mock database
	s.bills.On("Find", billId).Return(bill, nil).Once()
	s.bills.On("RemoveFormal", &bill, formalId).Return(nil).Once()
	// Run test
	err := s.h.RemoveBillFormal(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *AdminBillSuite) TestGetBillStats() {
	bill := model.Bill{
		Model: model.Model{ID: uuid.New()},
		Name:  "Test",
	}
	fbs := []model.FormalCostBreakdown{
		{
			FormalID:   uuid.New(),
			Name:       "Formal 1",
			Price:      10,
			GuestPrice: 20,
			DateTime:   time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			Standard:   11,
			Guest:      21,
		},
		{
			FormalID:   uuid.New(),
			Name:       "Formal 2",
			Price:      21.2,
			GuestPrice: 31.3,
			DateTime:   time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC),
			Standard:   12,
			Guest:      15,
		},
	}
	ubs := []model.UserCostBreakdown{
		{
			Email: "abc123@cam.ac.uk",
			Cost:  14,
		},
		{
			Email: "def456@cam.ac.uk",
			Cost:  73.4,
		},
	}
	expectedJSON := `{
		"formals": [
			{
				"formalId": "` + fbs[0].FormalID.String() + `",
				"formalName": "Formal 1",
				"price": 10,
				"guestPrice": 20,
				"dateTime": "2020-01-01T00:00:00Z",
				"standard": 11,
				"guest": 21
			},
			{
				"formalId": "` + fbs[1].FormalID.String() + `",
				"formalName": "Formal 2",
				"price": 21.2,
				"guestPrice": 31.3,
				"dateTime": "2020-02-01T00:00:00Z",
				"standard": 12,
				"guest": 15
			}
		],
		"users": [
			{
				"userEmail": "abc123@cam.ac.uk",
				"cost": 14
			},
			{
				"userEmail": "def456@cam.ac.uk",
				"cost": 73.4
			}
		]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/bills/", bill.ID.String(), "/stats"), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(bill.ID.String())
	// Mock database
	s.bills.On("Find", bill.ID).Return(bill, nil).Once()
	s.bills.On("GetCostBreakdown", &bill).Return(fbs, nil).Once()
	s.bills.On("GetCostBreakdownByUser", &bill).Return(ubs, nil).Once()
	// Run test
	err := s.h.GetBillStats(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminBillSuite) TestGetBillFormalStatsCSV() {
	bill := model.Bill{
		Model: model.Model{ID: uuid.New()},
		Name:  "Test",
	}
	fbs := []model.FormalCostBreakdown{
		{
			FormalID:   uuid.New(),
			Name:       "Formal 1",
			Price:      10,
			GuestPrice: 20,
			DateTime:   time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			Standard:   11,
			Guest:      21,
		},
		{
			FormalID:   uuid.New(),
			Name:       "Formal 2",
			Price:      21.2,
			GuestPrice: 31.3,
			DateTime:   time.Date(2020, time.February, 1, 0, 0, 0, 0, time.UTC),
			Standard:   12,
			Guest:      15,
		},
	}
	expectedBody := strings.ReplaceAll(
		`Formal,Date,King's Tickets,King's Price,Guest Tickets,Guest Price,Total
		Formal 1,Jan 1 2020,11,10.00,21,20.00,530.00
		Formal 2,Feb 1 2020,12,21.20,15,31.30,723.90
		Total,,23,,36,,1253.90`,
		"\t", "",
	)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/bills/", bill.ID.String(), "/stats/csv"), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(bill.ID.String())
	// Mock database
	s.bills.On("Find", bill.ID).Return(bill, nil).Once()
	s.bills.On("GetCostBreakdown", &bill).Return(fbs, nil).Once()
	// Run test
	err := s.h.GetBillFormalStatsCSV(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal(expectedBody, strings.TrimSpace(rec.Body.String()))
	s.Equal(`attachment; filename="formal_costs.csv"`, rec.Header().Get("Content-Disposition"))
}

func (s *AdminBillSuite) TestGetBillUserStatsCSV() {
	bill := model.Bill{
		Model: model.Model{ID: uuid.New()},
		Name:  "Test",
	}
	ubs := []model.UserCostBreakdown{
		{
			Email: "abc123@cam.ac.uk",
			Cost:  14,
		},
		{
			Email: "def456@cam.ac.uk",
			Cost:  73.4,
		},
	}
	expectedBody := strings.ReplaceAll(
		`CRSID,Total
		abc123,14.00
		def456,73.40
		Total,87.40`,
		"\t", "",
	)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/bills/", bill.ID.String(), "/stats/csv"), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(bill.ID.String())
	// Mock database
	s.bills.On("Find", bill.ID).Return(bill, nil).Once()
	s.bills.On("GetCostBreakdownByUser", &bill).Return(ubs, nil).Once()
	// Run test
	err := s.h.GetBillUserStatsCSV(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal(expectedBody, strings.TrimSpace(rec.Body.String()))
	s.Equal(`attachment; filename="user_costs.csv"`, rec.Header().Get("Content-Disposition"))
}

func TestBillSuite(t *testing.T) {
	suite.Run(t, new(AdminBillSuite))
}

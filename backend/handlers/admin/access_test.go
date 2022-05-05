package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	. "github.com/kcsu/store/handlers/admin"
	"github.com/kcsu/store/middleware"
	mocks "github.com/kcsu/store/mocks/middleware"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const expectedJSON = `[
	{
		"id": "c1fda423-0ed7-48a5-ba45-3b9246d83b7c",
		"email": "abc123@cam.ac.uk",
		"message": "ABC did XYZ",
		"createdAt": "0001-01-01T00:00:00Z",
		"metadata": {
			"test": "1"
		}
	},
	{
		"id": "502a58ba-3c44-490d-8bd7-0ef45da6b152",
		"email": "def123@cam.ac.uk",
		"message": "DEF did XYZ",
		"createdAt": "0001-01-01T00:00:00Z",
		"metadata": {
			"test": "2"
		}
	}
]`

func TestGetAccessLogs(t *testing.T) {
	// Init handler
	h := new(AdminHandler)
	a := mocks.NewAccess(t)
	h.Access = a
	// HTTP
	e := echo.New()
	e.Validator = middleware.NewValidator()
	req := httptest.NewRequest(http.MethodGet, "/logs?page=3&size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock database
	a.On("Get", 3, 10).Return([]model.AccessLog{
		{
			ID:       uuid.MustParse("c1fda423-0ed7-48a5-ba45-3b9246d83b7c"),
			Email:    "abc123@cam.ac.uk",
			Message:  "ABC did XYZ",
			Metadata: []byte(`{"test": "1"}`),
		},
		{
			ID:       uuid.MustParse("502a58ba-3c44-490d-8bd7-0ef45da6b152"),
			Email:    "def123@cam.ac.uk",
			Message:  "DEF did XYZ",
			Metadata: []byte(`{"test": "2"}`),
		},
	}, nil)
	// Run
	err := h.GetAccess(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expectedJSON, rec.Body.String())
}

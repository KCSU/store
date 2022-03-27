package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kcsu/store/middleware"
	am "github.com/kcsu/store/mocks/auth"
	um "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRBAC(t *testing.T) {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name        string
		resource    string
		action      string
		permissions []model.Permission
		wants       *wants
	}
	tests := []test{
		{
			"No Match",
			"formals",
			"write",
			[]model.Permission{
				{
					Resource: "formals",
					Action:   "read",
				},
				{
					Resource: "groups",
					Action:   "write",
				},
				{
					Resource: "*",
					Action:   "delete",
				},
				{
					Resource: "permissions",
					Action:   "*",
				},
			},
			&wants{http.StatusForbidden, "Forbidden"},
		},
		{
			"Wildcard Match Action",
			"groups",
			"delete",
			[]model.Permission{
				{
					Resource: "tickets",
					Action:   "write",
				},
				{
					Resource: "groups",
					Action:   "*",
				},
			},
			nil,
		},
		{
			"Wildcard Match Resource",
			"billing",
			"read",
			[]model.Permission{
				{
					Resource: "*",
					Action:   "read",
				},
				{
					Resource: "groups",
					Action:   "write",
				},
			},
			nil,
		},
		{
			"Match",
			"tickets",
			"write",
			[]model.Permission{
				{
					Resource: "tickets",
					Action:   "write",
				},
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			auth := new(am.Auth)
			users := new(um.UserStore)
			rbac := NewRBAC(RbacConfig{
				Auth:  auth,
				Users: users,
			})
			middleware := rbac.M(test.resource, test.action)
			// HTTP
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			auth.On("GetUserId", c).Return(12).Once()
			user := model.User{
				Model: model.Model{ID: 12},
				Name:  "James Holden",
				Email: "jh123@cam.ac.uk",
			}
			users.On("Find", 12).Return(user, nil).Once()
			users.On("Permissions", &user).Return(test.permissions, nil).Once()
			// Test
			h := middleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})
			err := h(c)
			if test.wants == nil {
				assert.NoError(t, err)
			} else {
				var he *echo.HTTPError
				if assert.ErrorAs(t, err, &he) {
					assert.Equal(t, test.wants.code, he.Code)
					assert.Equal(t, test.wants.message, he.Message)
				}
			}
			auth.AssertExpectations(t)
			users.AssertExpectations(t)
		})
	}

}

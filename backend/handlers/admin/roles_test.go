package admin_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/kcsu/store/handlers/admin"
	"github.com/kcsu/store/middleware"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminRoleSuite struct {
	suite.Suite
	h     *AdminHandler
	roles *mocks.RoleStore
}

func (s *AdminRoleSuite) SetupTest() {
	s.h = new(AdminHandler)
	s.roles = new(mocks.RoleStore)
	s.h.Roles = s.roles
}

func (s *AdminRoleSuite) TestGetRoles() {
	const expectedJSON = `[
		{
			"id":45,
			"createdAt":"0001-01-01T00:00:00Z",
			"updatedAt":"0001-01-01T00:00:00Z",
			"deletedAt":null,
			"name":"Admin",
			"permissions":[
				{
					"id":23,
					"resource":"groups",
					"action":"read"
				},
				{
					"id":11,
					"resource":"formals",
					"action":"write"
				}
			]
		},
		{
			"id":4,
			"createdAt":"0001-01-01T00:00:00Z",
			"updatedAt":"0001-01-01T00:00:00Z",
			"deletedAt":null,
			"name":"Something",
			"permissions":[
				{
					"id":31,
					"resource":"tickets",
					"action":"*"
				}
			]
		}
	]`
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock DB
	roles := []model.Role{
		{
			Model: model.Model{ID: 45},
			Name:  "Admin",
			Permissions: []model.Permission{
				{
					ID:       23,
					RoleID:   45,
					Resource: "groups",
					Action:   "read",
				},
				{
					ID:       11,
					RoleID:   45,
					Resource: "formals",
					Action:   "write",
				},
			},
		},
		{
			Model: model.Model{ID: 4},
			Name:  "Something",
			Permissions: []model.Permission{
				{
					ID:       31,
					RoleID:   4,
					Resource: "tickets",
					Action:   "*",
				},
			},
		},
	}
	s.roles.On("Get").Return(roles, nil)
	// Run test
	err := s.h.GetRoles(c)
	s.NoError(err)
	s.roles.AssertExpectations(s.T())
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminRoleSuite) TestGetUserRoles() {
	const expectedJSON = `[
		{
		  "userEmail": "abc123@cam.ac.uk",
		  "userName": "A. Bell",
		  "roleName": "Admin",
		  "userId": 26,
		  "roleId": 45
		}
	]`
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock DB
	userRoles := []model.UserRole{{
		Role: model.Role{
			Model: model.Model{ID: 45},
			Name:  "Admin",
		},
		User: model.User{
			Model: model.Model{ID: 26},
			Email: "abc123@cam.ac.uk",
			Name:  "A. Bell",
		},
		RoleID: 45,
		UserID: 26,
	}}
	s.roles.On("GetUserRoles").Return(userRoles, nil)
	// Run test
	err := s.h.GetUserRoles(c)
	s.NoError(err)
	s.roles.AssertExpectations(s.T())
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminRoleSuite) TestCreatePermission() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name       string
		body       string
		valid      bool
		role       *model.Role
		permission *model.Permission
		wants      *wants
	}
	tests := []test{
		{
			"Role Not Found",
			`{
				"resource": "formals",
				"action": "read",
				"roleId": 5
			}`,
			true,
			nil,
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
		{
			"Invalid Resource",
			`{
				"resource": "*s",
				"action": "read",
				"roleId": 5
			}`,
			false,
			nil,
			nil,
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'PermissionDto.Resource' Error:Field validation for 'Resource' failed on the 'alpha|eq=*' tag",
			},
		},
		{
			"Invalid Action",
			`{
				"resource": "*",
				"action": "1abc",
				"roleId": 5
			}`,
			false,
			nil,
			nil,
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'PermissionDto.Action' Error:Field validation for 'Action' failed on the 'alpha|eq=*' tag",
			},
		},
		{
			"Should Create",
			`{
				"resource": "formals",
				"action": "*",
				"roleId": 5
			}`,
			true,
			&model.Role{
				Name:  "My Role",
				Model: model.Model{ID: 5},
			},
			&model.Permission{
				Resource: "formals",
				Action:   "*",
				RoleID:   5,
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
				http.MethodPost, "/permissions", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			if test.valid {
				if test.role != nil {
					s.roles.On("Find", 5).Return(*test.role, nil).Once()
				} else {
					s.roles.On("Find", 5).Return(model.Role{}, gorm.ErrRecordNotFound).Once()
				}
				if test.permission != nil {
					s.roles.On("CreatePermission", test.permission).
						Return(nil).Once()
				}
			}

			// Test
			err := s.h.CreatePermission(c)
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
			s.roles.AssertExpectations(s.T())
		})
	}
}

func TestAdminRoleSuite(t *testing.T) {
	suite.Run(t, new(AdminRoleSuite))
}

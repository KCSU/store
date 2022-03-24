package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kcsu/store/handlers/admin"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
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

func TestAdminRoleSuite(t *testing.T) {
	suite.Run(t, new(AdminRoleSuite))
}

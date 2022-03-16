package admin_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/kcsu/store/handlers/admin"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type AdminGroupSuite struct {
	suite.Suite
	h      *AdminHandler
	groups *mocks.GroupStore
}

func (s *AdminGroupSuite) SetupTest() {
	// Init handler
	s.h = new(AdminHandler)
	s.groups = new(mocks.GroupStore)
	s.h.Groups = s.groups
}

func (s *AdminGroupSuite) TestGetGroups() {
	const expectedJSON = `[
		{
			"id":1,
			"createdAt":"0001-01-01T00:00:00Z",
			"updatedAt":"0001-01-01T00:00:00Z",
			"deletedAt":null,
			"name":"Group A",
			"type":"inst",
			"lookup":"GRPA"
		},
		{
			"id":51,
			"createdAt":"0001-01-01T00:00:00Z",
			"updatedAt":"0001-01-01T00:00:00Z",
			"deletedAt":null,
			"name":"Group B",
			"type":"group",
			"lookup":"GRPB"
		}
	]`
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock database
	groups := []model.Group{
		{
			Model:  model.Model{ID: 1},
			Name:   "Group A",
			Type:   "inst",
			Lookup: "GRPA",
		},
		{
			Model:  model.Model{ID: 51},
			Name:   "Group B",
			Type:   "group",
			Lookup: "GRPB",
		},
	}
	s.groups.On("Get").Return(groups, nil)
	// Run test
	err := s.h.GetGroups(c)
	s.NoError(err)
	s.groups.AssertExpectations(s.T())
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminGroupSuite) TestGetGroup() {
	const expectedJSON = `{
		"id":34,
		"createdAt":"0001-01-01T00:00:00Z",
		"updatedAt":"0001-01-01T00:00:00Z",
		"deletedAt":null,
		"name":"My Group",
		"type":"inst",
		"lookup":"MGRP",
		"users":[
			{
				"groupId":34,
				"userEmail":"abc123@cam.ac.uk",
				"createdAt":"0001-01-01T00:00:00Z"
			},
			{
				"groupId":34,
				"userEmail":"def456@cam.ac.uk",
				"createdAt":"0001-01-01T00:00:00Z"
			}
		]
	}`
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/groups/34", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("34")
	// Mock database
	group := model.Group{
		Model:  model.Model{ID: 34},
		Name:   "My Group",
		Type:   "inst",
		Lookup: "MGRP",
		GroupUsers: []model.GroupUser{
			{
				GroupID:   34,
				UserEmail: "abc123@cam.ac.uk",
			},
			{
				GroupID:   34,
				UserEmail: "def456@cam.ac.uk",
			},
		},
	}
	s.groups.On("Find", 34).Return(group, nil)

	err := s.h.GetGroup(c)
	s.NoError(err)
	log.Print(rec.Body.String())
	s.groups.AssertExpectations(s.T())
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func TestAdminGroupSuite(t *testing.T) {
	suite.Run(t, new(AdminGroupSuite))
}

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
	lm "github.com/kcsu/store/mocks/lookup"
	mm "github.com/kcsu/store/mocks/middleware"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminGroupSuite struct {
	suite.Suite
	h      *AdminHandler
	groups *mocks.GroupStore
	lookup *lm.Lookup
}

func (s *AdminGroupSuite) SetupTest() {
	// Init handler
	s.h = new(AdminHandler)
	s.groups = mocks.NewGroupStore(s.T())
	s.lookup = lm.NewLookup(s.T())
	s.h.Groups = s.groups
	s.h.Lookup = s.lookup
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

func (s *AdminGroupSuite) TestGetGroups() {
	const expectedJSON = `[
		{
			"id": "6e503329-d12f-466f-b985-a78bb1fec364",
			"createdAt":"0001-01-01T00:00:00Z",
			"updatedAt":"0001-01-01T00:00:00Z",
			"deletedAt":null,
			"name":"Group A",
			"type":"inst",
			"lookup":"GRPA"
		},
		{
			"id": "89c89641-5bdd-4770-b40f-0e25305e3530",
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
			Model:  model.Model{ID: uuid.MustParse("6e503329-d12f-466f-b985-a78bb1fec364")},
			Name:   "Group A",
			Type:   "inst",
			Lookup: "GRPA",
		},
		{
			Model:  model.Model{ID: uuid.MustParse("89c89641-5bdd-4770-b40f-0e25305e3530")},
			Name:   "Group B",
			Type:   "group",
			Lookup: "GRPB",
		},
	}
	s.groups.On("Get").Return(groups, nil)
	// Run test
	err := s.h.GetGroups(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminGroupSuite) TestGetGroup() {
	const expectedJSON = `{
		"id": "a2140fcf-89a8-4c90-bb24-c85f21aeb2e0",
		"createdAt":"0001-01-01T00:00:00Z",
		"updatedAt":"0001-01-01T00:00:00Z",
		"deletedAt":null,
		"name":"My Group",
		"type":"inst",
		"lookup":"MGRP",
		"users":[
			{
				"groupId": "a2140fcf-89a8-4c90-bb24-c85f21aeb2e0",
				"userEmail":"abc123@cam.ac.uk",
				"isManual":false,
				"createdAt":"0001-01-01T00:00:00Z"
			},
			{
				"groupId": "a2140fcf-89a8-4c90-bb24-c85f21aeb2e0",
				"userEmail":"def456@cam.ac.uk",
				"isManual":true,
				"createdAt":"0001-01-01T00:00:00Z"
			}
		]
	}`
	// Init HTTP
	e := echo.New()
	id := uuid.MustParse("a2140fcf-89a8-4c90-bb24-c85f21aeb2e0")
	req := httptest.NewRequest(http.MethodGet, "/groups/a2140fcf-89a8-4c90-bb24-c85f21aeb2e0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	// Mock database
	group := model.Group{
		Model:  model.Model{ID: id},
		Name:   "My Group",
		Type:   "inst",
		Lookup: "MGRP",
		GroupUsers: []model.GroupUser{
			{
				GroupID:   id,
				UserEmail: "abc123@cam.ac.uk",
				IsManual:  false,
			},
			{
				GroupID:   id,
				UserEmail: "def456@cam.ac.uk",
				IsManual:  true,
			},
		},
	}
	s.groups.On("Find", id).Return(group, nil)

	err := s.h.GetGroup(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(expectedJSON, rec.Body.String())
}

func (s *AdminGroupSuite) TestCreateGroup() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name  string
		body  string
		group model.Group
		wants *wants
	}
	tests := []test{
		{
			"Missing Lookup",
			`{
				"name": "A group",
				"type": "inst"
			}`,
			model.Group{},
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'AdminGroupDto.Lookup' Error:Field validation for 'Lookup' failed on the 'required_unless' tag",
			},
		},
		{
			"Should Create: manual",
			`{
				"name": "Manual Group",
				"type": "manual"
			}`,
			model.Group{
				Name: "Manual Group",
				Type: "manual",
			},
			nil,
		},
		{
			"Should Create: institution",
			`{
				"name": "My Group",
				"type": "inst",
				"lookup": "MYGRP"
			}`,
			model.Group{
				Name:   "My Group",
				Type:   "inst",
				Lookup: "MYGRP",
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
				http.MethodPost, "/groups", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			if test.wants == nil {
				s.groups.On("Create", &test.group).Return(nil).Once()
			}
			// Test
			err := s.h.CreateGroup(c)
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

func (s *AdminGroupSuite) TestUpdateGroup() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name  string
		body  string
		group model.Group
		wants *wants
	}
	tests := []test{
		{
			"Missing Lookup",
			`{
				"name": "Another Group",
				"type": "group"
			}`,
			model.Group{},
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'AdminGroupDto.Lookup' Error:Field validation for 'Lookup' failed on the 'required_unless' tag",
			},
		},
		{
			"Should Update: manual",
			`{
				"name": "Manual Group",
				"type": "manual"
			}`,
			model.Group{
				Model: model.Model{ID: uuid.New()},
				Name:  "Manual Group",
				Type:  "manual",
			},
			nil,
		},
		{
			"Should Update: institution",
			`{
				"name": "My Group",
				"type": "inst",
				"lookup": "MYGRP"
			}`,
			model.Group{
				Model:  model.Model{ID: uuid.New()},
				Name:   "My Group",
				Type:   "inst",
				Lookup: "MYGRP",
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
				http.MethodPut,
				fmt.Sprint("/groups/", test.group.ID),
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.group.ID.String())
			// Mock
			if test.wants == nil {
				s.groups.On("Find", test.group.ID).Return(
					model.Group{
						Model:  model.Model{ID: test.group.ID},
						Name:   "Initial",
						Type:   "initial",
						Lookup: "initial",
					}, nil,
				).Once()
				s.groups.On("Update", &test.group).Return(nil).Once()
			}
			// Test
			err := s.h.UpdateGroup(c)
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

func (s *AdminGroupSuite) TestDeleteGroup() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name  string
		id    uuid.UUID
		group model.Group
		wants *wants
	}
	tests := []test{
		{
			"Group Not Found",
			uuid.New(),
			model.Group{},
			&wants{http.StatusNotFound, "Not Found"},
		},
		{
			"Should Delete",
			uuid.MustParse("1109c678-e0f2-4675-94eb-d0b4bb61c862"),
			model.Group{
				Model:  model.Model{ID: uuid.MustParse("1109c678-e0f2-4675-94eb-d0b4bb61c862")},
				Name:   "Group",
				Type:   "inst",
				Lookup: "GRP",
			},
			nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			e := echo.New()
			// HTTP
			route := fmt.Sprint("/groups/", test.id)
			req := httptest.NewRequest(
				http.MethodDelete, route, nil,
			)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.id.String())
			// Mock
			if test.wants == nil {
				s.groups.On("Find", test.id).Return(
					test.group, nil,
				).Once()
				s.groups.On("Delete", &test.group).Return(nil).Once()
			} else {
				s.groups.On("Find", test.id).Return(
					model.Group{}, gorm.ErrRecordNotFound,
				).Once()
			}
			// Test
			err := s.h.DeleteGroup(c)
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

func (s *AdminGroupSuite) TestAddGroupUser() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name      string
		body      string
		email     string
		findGroup bool
		group     *model.Group
		wants     *wants
	}
	id := uuid.New()
	tests := []test{
		{
			"Invalid email",
			`{
				"userEmail": "hello@test"
			}`,
			"",
			false,
			nil,
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'GroupUserDto.Email' Error:Field validation for 'Email' failed on the 'email' tag",
			},
		},
		{
			"Non-existent group",
			`{
				"userEmail": "abc123@cam.ac.uk"
			}`,
			"abc123@cam.ac.uk",
			true,
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
		{
			"Should add user",
			`{
				"userEmail": "def456@cam.ac.uk"
			}`,
			"def456@cam.ac.uk",
			true,
			&model.Group{
				Model: model.Model{ID: id},
				Name:  "My Group",
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
				http.MethodPost,
				fmt.Sprint("/groups/", id, "/users"),
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(id.String())
			// Mock
			if test.findGroup {
				if test.group == nil {
					s.groups.On("Find", id).Once().Return(model.Group{}, gorm.ErrRecordNotFound)
				} else {
					s.groups.On("Find", id).Once().Return(*test.group, nil)
				}
			}
			if test.wants == nil {
				s.groups.On("AddUser", test.group, test.email).Once().Return(nil)
			}

			// Test
			err := s.h.AddGroupUser(c)
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

func (s *AdminGroupSuite) TestRemoveGroupUser() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name      string
		body      string
		email     string
		findGroup bool
		group     *model.Group
		wants     *wants
	}
	id := uuid.New()
	tests := []test{
		{
			"Invalid email",
			`{
				"userEmail": "hello@test"
			}`,
			"",
			false,
			nil,
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'GroupUserDto.Email' Error:Field validation for 'Email' failed on the 'email' tag",
			},
		},
		{
			"Non-existent group",
			`{
				"userEmail": "abc123@cam.ac.uk"
			}`,
			"abc123@cam.ac.uk",
			true,
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
		{
			"Should remove user",
			`{
				"userEmail": "def456@cam.ac.uk"
			}`,
			"def456@cam.ac.uk",
			true,
			&model.Group{
				Model: model.Model{ID: id},
				Name:  "My Group",
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
				http.MethodPost,
				fmt.Sprint("/groups/", id, "/users"),
				strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(id.String())
			// Mock
			if test.findGroup {
				if test.group == nil {
					s.groups.On("Find", id).Once().Return(model.Group{}, gorm.ErrRecordNotFound)
				} else {
					s.groups.On("Find", id).Once().Return(*test.group, nil)
				}
			}
			if test.wants == nil {
				s.groups.On("RemoveUser", test.group, test.email).Once().Return(nil)
			}

			// Test
			err := s.h.RemoveGroupUser(c)
			if test.wants == nil {
				s.NoError(err)
				s.Equal(http.StatusNoContent, rec.Code)
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

func (s *AdminGroupSuite) TestLookupGroupUsers() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name  string
		group model.Group
		users []model.GroupUser
		wants *wants
	}
	tests := []test{
		{
			"Should Replace",
			model.Group{
				Model:  model.Model{ID: uuid.New()},
				Name:   "My Group",
				Type:   "inst",
				Lookup: "MYGRP",
			},
			[]model.GroupUser{
				{UserEmail: "abc123@cam.ac.uk"},
				{UserEmail: "def456@cam.ac.uk"},
			},
			nil,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			e := echo.New()
			// HTTP
			route := fmt.Sprint("/groups/", test.group.ID, "/users/lookup")
			req := httptest.NewRequest(
				http.MethodPost, route, nil,
			)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.group.ID.String())
			// Mock
			s.groups.On("Find", test.group.ID).Return(
				test.group, nil,
			).Once()
			if test.wants == nil {
				s.lookup.On("ProcessGroup", test.group).Return(nil).Once()
			}
			// Test
			err := s.h.LookupGroupUsers(c)
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

func TestAdminGroupSuite(t *testing.T) {
	suite.Run(t, new(AdminGroupSuite))
}

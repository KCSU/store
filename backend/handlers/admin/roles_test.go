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
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AdminRoleSuite struct {
	suite.Suite
	h     *AdminHandler
	roles *mocks.RoleStore
	users *mocks.UserStore
}

func (s *AdminRoleSuite) SetupTest() {
	s.h = new(AdminHandler)
	s.roles = new(mocks.RoleStore)
	s.users = new(mocks.UserStore)
	s.h.Roles = s.roles
	s.h.Users = s.users
}

func (s *AdminRoleSuite) TestGetRoles() {
	const expectedJSON = `[
		{
			"id": "bec928b2-fd8e-4e46-ab6a-811c41fd5260",
			"createdAt":"0001-01-01T00:00:00Z",
			"updatedAt":"0001-01-01T00:00:00Z",
			"deletedAt":null,
			"name":"Admin",
			"permissions":[
				{
					"id": "7e3ef068-9e8f-48c4-9f63-9d016f815b64",
					"resource":"groups",
					"action":"read"
				},
				{
					"id": "29c37fe8-6e5b-46d9-a988-40a3ffe7d58a",
					"resource":"formals",
					"action":"write"
				}
			]
		},
		{
			"id": "bc71fbe4-de50-49b0-b7aa-381c0da566ef",
			"createdAt":"0001-01-01T00:00:00Z",
			"updatedAt":"0001-01-01T00:00:00Z",
			"deletedAt":null,
			"name":"Something",
			"permissions":[
				{
					"id": "e0bd7d92-4dd6-4a97-b4b2-d9ef30f89d55",
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
			Model: model.Model{ID: uuid.MustParse("bec928b2-fd8e-4e46-ab6a-811c41fd5260")},
			Name:  "Admin",
			Permissions: []model.Permission{
				{
					ID:       uuid.MustParse("7e3ef068-9e8f-48c4-9f63-9d016f815b64"),
					RoleID:   uuid.MustParse("bec928b2-fd8e-4e46-ab6a-811c41fd5260"),
					Resource: "groups",
					Action:   "read",
				},
				{
					ID:       uuid.MustParse("29c37fe8-6e5b-46d9-a988-40a3ffe7d58a"),
					RoleID:   uuid.MustParse("bec928b2-fd8e-4e46-ab6a-811c41fd5260"),
					Resource: "formals",
					Action:   "write",
				},
			},
		},
		{
			Model: model.Model{ID: uuid.MustParse("bc71fbe4-de50-49b0-b7aa-381c0da566ef")},
			Name:  "Something",
			Permissions: []model.Permission{
				{
					ID:       uuid.MustParse("e0bd7d92-4dd6-4a97-b4b2-d9ef30f89d55"),
					RoleID:   uuid.MustParse("bc71fbe4-de50-49b0-b7aa-381c0da566ef"),
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
		  "userId": "e14792f2-6faa-4378-89b0-e2261bdd8dc4",
		  "roleId": "5b77d20e-0f14-4680-b8b5-cc962ddd8eb9"
		}
	]`
	// Init HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock DB
	userId := uuid.MustParse("e14792f2-6faa-4378-89b0-e2261bdd8dc4")
	roleId := uuid.MustParse("5b77d20e-0f14-4680-b8b5-cc962ddd8eb9")
	userRoles := []model.UserRole{{
		Role: model.Role{
			Model: model.Model{ID: roleId},
			Name:  "Admin",
		},
		User: model.User{
			Model: model.Model{ID: userId},
			Email: "abc123@cam.ac.uk",
			Name:  "A. Bell",
		},
		RoleID: roleId,
		UserID: userId,
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
	roleId := uuid.MustParse("b7780926-570a-48df-a4a9-09970f26db55")
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
				"roleId": "b7780926-570a-48df-a4a9-09970f26db55"
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
				"roleId": "b7780926-570a-48df-a4a9-09970f26db55"
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
				"roleId": "b7780926-570a-48df-a4a9-09970f26db55"
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
				"roleId": "b7780926-570a-48df-a4a9-09970f26db55"
			}`,
			true,
			&model.Role{
				Name:  "My Role",
				Model: model.Model{ID: roleId},
			},
			&model.Permission{
				Resource: "formals",
				Action:   "*",
				RoleID:   roleId,
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
					s.roles.On("Find", roleId).Return(*test.role, nil).Once()
				} else {
					s.roles.On("Find", roleId).Return(model.Role{}, gorm.ErrRecordNotFound).Once()
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
		})
	}
	s.roles.AssertExpectations(s.T())
}

func (s *AdminRoleSuite) TestDeletePermission() {
	e := echo.New()
	id := uuid.New()
	route := fmt.Sprint("/permissions/", id)
	req := httptest.NewRequest(
		http.MethodDelete, route, nil,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())
	// Mock
	s.roles.On("DeletePermission", id).Return(nil).Once()
	// Test
	err := s.h.DeletePermission(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.roles.AssertExpectations(s.T())
}

func (s *AdminRoleSuite) TestCreateRole() {
	body := `{"name": "Admin"}`
	e := echo.New()
	e.Validator = middleware.NewValidator()
	req := httptest.NewRequest(
		http.MethodPost, "/roles", strings.NewReader(body),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock
	s.roles.On("Create", &model.Role{Name: "Admin"}).Return(nil).Once()
	// Test
	err := s.h.CreateRole(c)
	s.NoError(err)
	s.Equal(http.StatusCreated, rec.Code)
	s.roles.AssertExpectations(s.T())
}

func (s *AdminRoleSuite) TestUpdateRole() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name  string
		body  string
		role  model.Role
		wants *wants
	}
	tests := []test{
		{
			"Name Too Short",
			`{"name": "Ad"}`,
			model.Role{},
			&wants{
				http.StatusUnprocessableEntity,
				"Key: 'RoleDto.Name' Error:Field validation for 'Name' failed on the 'min' tag",
			},
		},
		{
			"Should Update",
			`{"name": "Admin"}`,
			model.Role{
				Model: model.Model{ID: uuid.New()},
				Name:  "Admin",
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
				http.MethodPut, "/roles/17", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.role.ID.String())
			// Mock
			if test.wants == nil {
				s.roles.On("Find", test.role.ID).Return(
					model.Role{
						Model: model.Model{ID: test.role.ID},
						Name:  "initial",
					}, nil,
				).Once()
				s.roles.On("Update", &test.role).Return(nil).Once()
			}
			// Test
			err := s.h.UpdateRole(c)
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
	s.roles.AssertExpectations(s.T())
}

func (s *AdminRoleSuite) TestDeleteRole() {
	e := echo.New()
	id := uuid.New()
	role := model.Role{
		Model: model.Model{ID: id},
		Name:  "Admin",
	}
	route := fmt.Sprint("/roles/", id)
	req := httptest.NewRequest(
		http.MethodDelete, route, nil,
	)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())
	// Mock
	s.roles.On("Find", id).Return(role, nil).Once()
	s.roles.On("Delete", &role).Return(nil).Once()
	// Test
	err := s.h.DeleteRole(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.roles.AssertExpectations(s.T())
}

func (s *AdminRoleSuite) TestAddUserRole() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		body   string
		roleId uuid.UUID
		email  string
		role   *model.Role
		user   *model.User
		wants  *wants
	}
	tests := []test{
		{
			"Role Not Found",
			`{
				"roleId": "72fa0314-f29c-4ed0-8654-c50626e4d481",
				"email": "abc123@cam.ac.uk"
			}`,
			uuid.MustParse("72fa0314-f29c-4ed0-8654-c50626e4d481"),
			"abc123@cam.ac.uk",
			nil,
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
		{
			"User Not Found",
			`{
				"roleId": "72fa0314-f29c-4ed0-8654-c50626e4d481",
				"email": "abc123@cam.ac.uk"
			}`,
			uuid.MustParse("72fa0314-f29c-4ed0-8654-c50626e4d481"),
			"abc123@cam.ac.uk",
			&model.Role{
				Model: model.Model{ID: uuid.MustParse("72fa0314-f29c-4ed0-8654-c50626e4d481")},
				Name:  "Admin",
			},
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
		{
			"Should Add",
			`{
				"roleId": "d8cb8bde-6b6f-489e-916c-4830a5bd605c",
				"email": "def456@cam.ac.uk"
			}`,
			uuid.MustParse("d8cb8bde-6b6f-489e-916c-4830a5bd605c"),
			"def456@cam.ac.uk",
			&model.Role{
				Model: model.Model{ID: uuid.MustParse("d8cb8bde-6b6f-489e-916c-4830a5bd605c")},
				Name:  "Admin",
			},
			&model.User{
				Model: model.Model{ID: uuid.New()},
				Name:  "James Holden",
				Email: "def456@cam.ac.uk",
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
				http.MethodPost, "/roles/users", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			if test.role != nil {
				s.roles.On("Find", test.roleId).Return(*test.role, nil).Once()
				if test.user != nil {
					s.users.On("FindByEmail", test.email).
						Return(*test.user, nil).Once()
				} else {
					s.users.On("FindByEmail", test.email).
						Return(model.User{}, gorm.ErrRecordNotFound).Once()
				}
			} else {
				s.roles.On("Find", test.roleId).
					Return(model.Role{}, gorm.ErrRecordNotFound).Once()
			}
			if test.wants == nil {
				s.roles.On("AddUserRole", test.role, test.user).
					Return(nil).Once()
			}
			// Test
			err := s.h.AddUserRole(c)
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
	s.roles.AssertExpectations(s.T())
	s.users.AssertExpectations(s.T())
}

func (s *AdminRoleSuite) TestRemoveUserRole() {
	type wants struct {
		code    int
		message string
	}
	type test struct {
		name   string
		body   string
		roleId uuid.UUID
		email  string
		role   *model.Role
		user   *model.User
		wants  *wants
	}
	tests := []test{
		{
			"Role Not Found",
			`{
				"roleId": "6f1816f1-bf17-49d0-8a28-30437413d414",
				"email": "abc123@cam.ac.uk"
			}`,
			uuid.MustParse("6f1816f1-bf17-49d0-8a28-30437413d414"),
			"abc123@cam.ac.uk",
			nil,
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
		{
			"User Not Found",
			`{
				"roleId": "8de2f033-8716-40d1-a19e-53b1a5104d9a",
				"email": "hij123@cam.ac.uk"
			}`,
			uuid.MustParse("8de2f033-8716-40d1-a19e-53b1a5104d9a"),
			"hij123@cam.ac.uk",
			&model.Role{
				Model: model.Model{ID: uuid.MustParse("8de2f033-8716-40d1-a19e-53b1a5104d9a")},
				Name:  "Admin",
			},
			nil,
			&wants{
				http.StatusNotFound,
				"Not Found",
			},
		},
		{
			"Should Remove",
			`{
				"roleId": "b8f179a5-48de-49e5-b2a2-8e9bdc316aab",
				"email": "def456@cam.ac.uk"
			}`,
			uuid.MustParse("b8f179a5-48de-49e5-b2a2-8e9bdc316aab"),
			"def456@cam.ac.uk",
			&model.Role{
				Model: model.Model{ID: uuid.MustParse("b8f179a5-48de-49e5-b2a2-8e9bdc316aab")},
				Name:  "Admin",
			},
			&model.User{
				Model: model.Model{ID: uuid.New()},
				Name:  "James Holden",
				Email: "def456@cam.ac.uk",
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
				http.MethodPost, "/roles/users", strings.NewReader(test.body),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// Mock
			if test.role != nil {
				s.roles.On("Find", test.roleId).Return(*test.role, nil).Once()
				if test.user != nil {
					s.users.On("FindByEmail", test.email).
						Return(*test.user, nil).Once()
				} else {
					s.users.On("FindByEmail", test.email).
						Return(model.User{}, gorm.ErrRecordNotFound).Once()
				}
			} else {
				s.roles.On("Find", test.roleId).
					Return(model.Role{}, gorm.ErrRecordNotFound).Once()
			}
			if test.wants == nil {
				s.roles.On("RemoveUserRole", test.role, test.user).
					Return(nil).Once()
			}
			// Test
			err := s.h.RemoveUserRole(c)
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
	s.roles.AssertExpectations(s.T())
	s.users.AssertExpectations(s.T())
}

func TestAdminRoleSuite(t *testing.T) {
	suite.Run(t, new(AdminRoleSuite))
}

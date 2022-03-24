package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/kcsu/store/auth"
	. "github.com/kcsu/store/handlers"
	am "github.com/kcsu/store/mocks/auth"
	um "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/quasoft/memstore"
	"github.com/stretchr/testify/suite"
)

const jwtSecret = "jwtSecret"
const redirect = "https://example.com/redirect"

type AuthSuite struct {
	suite.Suite
	h     *Handler
	auth  *am.Auth
	users *um.UserStore
}

func (a *AuthSuite) SetupTest() {
	a.auth = new(am.Auth)
	a.users = new(um.UserStore)
	a.h = &Handler{
		Auth:  a.auth,
		Users: a.users,
	}
	a.h.Config.JwtSecret = jwtSecret
	a.h.Config.OauthRedirectUrl = redirect
}

func (a *AuthSuite) TestGetUser() {
	// HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Mock
	user := model.User{
		Name:           "Kara Thrace",
		Email:          "kt494@cam.ac.uk",
		ProviderUserId: "123456",
	}
	user.ID = 5
	userJson := `{
		"id":5,
		"createdAt":"0001-01-01T00:00:00Z",
		"updatedAt":"0001-01-01T00:00:00Z",
		"deletedAt":null,
		"name":"Kara Thrace",
		"email":"kt494@cam.ac.uk",
		"groups": [
			{
				"id": 34,
				"name": "Battlestar"
			}
		],
		"permissions": [
			{
				"id": 12,
				"resource": "formals",
				"action": "read"
			}
		]
	}`
	a.auth.On("GetUserId", c).Return(int(user.ID))
	a.users.On("Find", int(user.ID)).Return(user, nil)
	a.users.On("Groups", &user).Return([]model.Group{{
		Model: model.Model{ID: 34},
		Name:  "Battlestar",
	}}, nil)
	a.users.On("Permissions", &user).Return([]model.Permission{{
		ID:       12,
		Resource: "formals",
		Action:   "read",
	}}, nil)
	// Test
	err := a.h.GetUser(c)
	a.NoError(err)
	a.Equal(http.StatusOK, rec.Code)
	a.JSONEq(userJson, rec.Body.String())
	a.auth.AssertExpectations(a.T())
	a.users.AssertExpectations(a.T())
}

func (a *AuthSuite) TestAuthCallback() {
	// HTTP
	e := echo.New()
	credential := "4815162342"
	f := make(url.Values)
	f.Set("credential", credential)
	// FIXME: form params for CSRF?
	req := httptest.NewRequest(
		http.MethodPost, "/", strings.NewReader(f.Encode()),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess := memstore.NewMemStore()
	c.Set("_session_store", sess)
	// Mock
	user := model.User{
		Name:           "John Locke",
		Email:          "jl815@cam.ac.uk",
		ProviderUserId: "109632",
	}
	user.ID = 415
	oauthUser := goth.User{
		Name:   user.Name,
		Email:  user.Email,
		UserID: user.ProviderUserId,
	}
	// a.auth.On("VerifyGoogleCsrfToken", c).Return(nil)
	a.auth.On(
		"CompleteUserAuth", c,
	).Return(oauthUser, nil)
	a.users.On("FindOrCreate", &oauthUser).Return(user, nil)
	// Test
	err := a.h.AuthCallback(c)
	a.NoError(err)
	a.Equal(http.StatusTemporaryRedirect, rec.Code)
	location, err := rec.Result().Location()
	a.NoError(err)
	a.Equal(redirect, location.String())
	// Check JWT
	session, err := sess.Get(req, "__session")
	a.NoError(err)
	a.NotNil(session.Values["_token"])
	token, err := jwt.ParseWithClaims(
		session.Values["_token"].(string), &auth.JwtClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	a.NoError(err)
	a.True(token.Valid)
	claims := token.Claims.(*auth.JwtClaims)
	a.Equal(strconv.Itoa(int(user.ID)), claims.Subject)
	a.Equal(user.Name, claims.Name)
	a.Equal(user.Email, claims.Email)
	a.auth.AssertExpectations(a.T())
	a.users.AssertExpectations(a.T())
}

func (a *AuthSuite) TestEmailConflict() {
	// HTTP
	e := echo.New()
	credential := "4815162342"
	f := make(url.Values)
	f.Set("credential", credential)
	// FIXME: form params for CSRF?
	req := httptest.NewRequest(
		http.MethodPost, "/", strings.NewReader(f.Encode()),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	oauthUser := goth.User{
		Name:   "Naomi Nagata",
		Email:  "nng56@cam.ac.uk",
		UserID: "175295",
	}
	a.auth.On(
		"CompleteUserAuth", c,
	).Return(oauthUser, nil)
	a.users.On("FindOrCreate", &oauthUser).Return(
		model.User{}, errors.New("invalid data"),
	)
	a.users.On("Exists", oauthUser.Email).Return(true, nil)
	// Test
	err := a.h.AuthCallback(c)

	var he *echo.HTTPError
	if a.ErrorAs(err, &he) {
		a.Equal(he.Code, http.StatusConflict)
		a.Equal(he.Message, "email is taken")
	}
	a.auth.AssertExpectations(a.T())
	a.users.AssertExpectations(a.T())
}

func (a *AuthSuite) TestLogout() {
	// HTTP
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess := memstore.NewMemStore()
	s, err := sess.New(req, "__session")
	a.Require().NoError(err)
	s.Values["_token"] = "token"
	s.Save(req, rec)
	c.Set("_session_store", sess)
	// Test
	err = a.h.Logout(c)
	a.NoError(err)
	a.Equal(http.StatusOK, rec.Code)
	s, err = sess.Get(req, "__session")
	a.NoError(err)
	a.NotContains(s.Values, "_token")
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

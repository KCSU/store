package handlers_test

import (
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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const jwtSecret = "jwtSecret"

type AuthSuite struct {
	suite.Suite
	h     *Handler
	auth  *am.Auth
	users *um.UserStore
}

func (a *AuthSuite) SetupSuite() {
	a.auth = new(am.Auth)
	a.users = new(um.UserStore)
	a.h = &Handler{
		Auth:  a.auth,
		Users: a.users,
	}
	a.h.Config.JwtSecret = jwtSecret
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
		"email":"kt494@cam.ac.uk"
	}`
	a.auth.On("GetUserId", c).Return(int(user.ID))
	a.users.On("Find", int(user.ID)).Return(user, nil)
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
	// Mock
	user := model.User{
		Name:           "John Locke",
		Email:          "jl815@cam.ac.uk",
		ProviderUserId: "109632",
	}
	user.ID = 415
	userJson := `{
		"id": 415,
		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z",
		"deletedAt": null,
		"name": "John Locke",
		"email": "jl815@cam.ac.uk"
	}`
	oauthUser := auth.OauthUser{
		Name:   user.Name,
		Email:  user.Email,
		UserID: user.ProviderUserId,
	}
	a.auth.On("VerifyGoogleCsrfToken", c).Return(nil)
	a.auth.On(
		"VerifyIdToken", credential, mock.Anything,
	).Return(&oauthUser, nil)
	a.users.On("FindOrCreate", &oauthUser).Return(user, nil)
	// Test
	err := a.h.AuthCallback(c)
	a.NoError(err)
	a.Equal(http.StatusOK, rec.Code)
	a.JSONEq(userJson, rec.Body.String())
	// Check JWT
	if a.Len(rec.Result().Cookies(), 1) {
		cookie := rec.Result().Cookies()[0]
		a.Equal(cookie.Name, "_token")
		token, err := jwt.ParseWithClaims(
			cookie.Value, &auth.JwtClaims{},
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
	}
	a.auth.AssertExpectations(a.T())
	a.users.AssertExpectations(a.T())
}

// TODO: test email conflict failure

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

package auth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/labstack/echo/v4"
	"github.com/quasoft/memstore"
	"github.com/stretchr/testify/suite"
)

type OAuthSuite struct {
	suite.Suite
	auth Auth
}

func (a *OAuthSuite) SetupTest() {
	store := memstore.NewMemStore([]byte{})
	a.auth = Init(&config.Config{}, store)
}

func (a *OAuthSuite) TestGetAuthUrl() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	rawUrl, err := a.auth.GetAuthUrl(c)
	a.Require().NoError(err)
	url, err := url.Parse(rawUrl)
	a.Require().NoError(err)
	a.Equal("accounts.google.com", url.Hostname())
	a.Equal("cam.ac.uk", url.Query().Get("hd"))
	a.Equal("email profile openid", url.Query().Get("scope"))
}

// TODO:
// func (a *OAuthSuite) TestCompleteUserAuth() {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	rec := httptest.NewRecorder()
// }

func TestOAuthSuite(t *testing.T) {
	suite.Run(t, new(OAuthSuite))
}

package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kcsu/store/config"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

const hostedDomain = "cam.ac.uk"

type IdTokenValidator func(context.Context, string, string) (*idtoken.Payload, error)

type Auth struct {
	clientId       string
	TokenValidator IdTokenValidator
}

type OauthUser struct {
	UserID string
	Name   string
	Email  string
}

func Init(c *config.Config) *Auth {
	return &Auth{
		clientId:       c.OauthClientKey,
		TokenValidator: idtoken.Validate,
	}
}

func (auth *Auth) VerifyGoogleCsrfToken(c echo.Context) error {
	cookie, err := c.Cookie("g_csrf_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No CSRF token in Cookie")
	}
	token := c.FormValue("g_csrf_token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No CSRF token in post body")
	}
	if cookie.Value != token {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to verify double submit cookie")
	}
	return nil
}

func (auth *Auth) VerifyIdToken(token string, c context.Context) (*OauthUser, error) {
	id, err := auth.TokenValidator(c, token, auth.clientId)
	if err != nil {
		return nil, err
	}
	// TODO: type checking
	if id.Claims["hd"] != hostedDomain {
		return nil, fmt.Errorf("invalid hosted domain: should be %s", hostedDomain)
	}
	if !strings.HasSuffix(id.Claims["email"].(string), "@"+hostedDomain) {
		return nil, fmt.Errorf("invalid email domain: should be %s", hostedDomain)
	}
	return &OauthUser{
		UserID: id.Subject,
		Name:   id.Claims["name"].(string),
		Email:  id.Claims["email"].(string),
	}, nil
}

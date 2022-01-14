package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/kcsu/store/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/idtoken"
)

func getCsrfError(form_token string, cookie_token string) error {
	e := echo.New()
	var req *http.Request
	if form_token != "" {
		f := make(url.Values)
		f.Set("g_csrf_token", form_token)
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	} else {
		req = httptest.NewRequest(http.MethodPost, "/", nil)
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	if cookie_token != "" {
		req.AddCookie(&http.Cookie{
			Name:  "g_csrf_token",
			Value: cookie_token,
		})
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	a := auth.Auth{}
	return a.VerifyGoogleCsrfToken(c)
}

func TestCsrfCookie(t *testing.T) {
	err := getCsrfError("abc123", "")
	assert.EqualError(t, err, "code=400, message=No CSRF token in Cookie")
}

func TestCsrfForm(t *testing.T) {
	err := getCsrfError("", "abc123")
	assert.EqualError(t, err, "code=400, message=No CSRF token in post body")
}

func TestCsrfEqual(t *testing.T) {
	err := getCsrfError("abc123", "def456")
	assert.EqualError(t, err, "code=400, message=Failed to verify double submit cookie")
}

func TestCsrfValid(t *testing.T) {
	err := getCsrfError("abc123", "abc123")
	assert.NoError(t, err)
}

func mockTokenValidator(payload idtoken.Payload) auth.IdTokenValidator {
	return func(c context.Context, s1, s2 string) (*idtoken.Payload, error) {
		return &payload, nil
	}
}

func TestVerifyHostedDomain(t *testing.T) {
	a := auth.Auth{
		TokenValidator: mockTokenValidator(idtoken.Payload{
			Claims: map[string]interface{}{
				// This should fail!!
				"hd":    "ox.ac.uk",
				"email": "abc12@cam.ac.uk",
			},
		}),
	}
	_, err := a.VerifyIdToken("", context.Background())
	assert.EqualError(t, err, "invalid hosted domain: should be cam.ac.uk")
}

func TestVerifyEmailSuffix(t *testing.T) {
	a := auth.Auth{
		TokenValidator: mockTokenValidator(idtoken.Payload{
			Claims: map[string]interface{}{
				"hd":    "cam.ac.uk",
				"email": "abc12@ox.ac.uk",
			},
		}),
	}
	_, err := a.VerifyIdToken("", context.Background())
	assert.EqualError(t, err, "invalid email domain: should be cam.ac.uk")
}

func TestVerifyUser(t *testing.T) {
	u := auth.OauthUser{
		Email:  "abc12@cam.ac.uk",
		UserID: "123456",
		Name:   "James Holden",
	}
	a := auth.Auth{
		TokenValidator: mockTokenValidator(idtoken.Payload{
			Subject: u.UserID,
			Claims: map[string]interface{}{
				"hd":    "cam.ac.uk",
				"email": u.Email,
				"name":  u.Name,
			},
		}),
	}
	user, err := a.VerifyIdToken("", context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, user, &u)
	}
}

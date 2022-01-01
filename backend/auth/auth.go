package auth

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/kcsu/store/config"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	sessionName  = "_oauth_session"
	providerName = "google"
	hostedDomain = "cam.ac.uk"
)

type Auth struct {
	goth.Provider
	Store sessions.Store
}

func Init(c *config.Config) *Auth {
	// TODO: supply cookie store from middleware
	store := sessions.NewCookieStore([]byte(c.CookieSecret))
	provider := google.New(c.OauthClientKey, c.OauthSecretKey, c.OauthCallbackUrl)
	// Should this be a config value?
	provider.SetHostedDomain(hostedDomain)
	return &Auth{
		Provider: provider,
		Store:    store,
	}
}

func (auth *Auth) GetAuthUrl(c echo.Context) (string, error) {
	state := auth.GenerateState()
	sess, err := auth.BeginAuth(state)
	if err != nil {
		return "", err
	}
	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	err = auth.StoreInSession(providerName, sess.Marshal(), c)

	if err != nil {
		return "", err
	}

	return url, err
}

func (auth *Auth) CompleteUserAuth(c echo.Context) (goth.User, error) {
	value, err := auth.GetFromSession(providerName, c.Request())
	if err != nil {
		return goth.User{}, err
	}
	// TODO: defer Logout(res, req)
	sess, err := auth.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	err = validateState(c, sess)
	if err != nil {
		return goth.User{}, err
	}

	user, err := auth.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	params := c.QueryParams()
	if params.Encode() == "" && c.Request().Method == http.MethodPost {
		params, err = c.FormParams()
		if err != nil {
			return goth.User{}, err
		}
	}

	// get new token and retry fetch
	_, err = sess.Authorize(auth, params)
	if err != nil {
		return goth.User{}, err
	}

	err = auth.StoreInSession(providerName, sess.Marshal(), c)

	if err != nil {
		return goth.User{}, err
	}

	gu, err := auth.FetchUser(sess)

	// verify hd and email claims
	if gu.RawData["hd"] != hostedDomain {
		return goth.User{}, fmt.Errorf("invalid hosted domain: should be %s", hostedDomain)
	}
	if !strings.HasSuffix(gu.Email, "@"+hostedDomain) {
		return goth.User{}, fmt.Errorf("invalid email domain: should be %s", hostedDomain)
	}

	return gu, err
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateState(c echo.Context, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	reqState := GetState(c)

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != reqState) {
		return errors.New("state token mismatch")
	}
	return nil
}

// NOTE: do we need to use request state if provided?
// See gothic docs.
// Generate a random base64-encoded nonce so that the state on
// the auth URL is unguessable, preventing CSRF attacks, as described in
//
// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
func (auth *Auth) GenerateState() string {
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("auth: source of randomness unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func GetState(c echo.Context) string {
	params := c.QueryParams()
	if params.Encode() == "" && c.Request().Method == http.MethodPost {
		return c.FormValue("state")
	}
	return params.Get("state")
}

func (auth *Auth) GetFromSession(key string, req *http.Request) (string, error) {
	session, _ := auth.Store.Get(req, sessionName)
	value, err := getSessionValue(session, key)
	if err != nil {
		return "", errors.New("could not find a matching session for this request")
	}

	return value, nil
}

func (auth *Auth) StoreInSession(key string, value string, c echo.Context) error {
	session, _ := auth.Store.New(c.Request(), sessionName)

	if err := updateSessionValue(session, key, value); err != nil {
		return err
	}

	return session.Save(c.Request(), c.Response())
}

func getSessionValue(session *sessions.Session, key string) (string, error) {
	value := session.Values[key]
	if value == nil {
		return "", fmt.Errorf("could not find a matching session for this request")
	}

	rdata := strings.NewReader(value.(string))
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func updateSessionValue(session *sessions.Session, key, value string) error {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return err
	}
	if err := gz.Flush(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	session.Values[key] = b.String()
	return nil
}

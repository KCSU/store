package auth

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/kcsu/store/config"
	"golang.org/x/oauth2"
)

type Auth struct {
	oauth2.Config
	ClientKey   string
	Secret      string
	CallbackURL string
	Store       sessions.Store
}

func Init(c *config.Config) *Auth {
	// TODO: implement
	return &Auth{
		Config: oauth2.Config{},
		// ClientKey: ,

	}
}

const (
	sessionName  = "_oauth_session"
	providerName = "google"
)

func (auth *Auth) GetAuthUrl(req *http.Request, res http.ResponseWriter) (string, error) {
	state := auth.GenerateState()
	sess := auth.BeginAuth(state)
	url := sess.AuthURL
	err := auth.StoreInSession(providerName, sess.Marshal(), req, res)

	if err != nil {
		return "", err
	}

	return url, nil
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

func (auth *Auth) BeginAuth(state string) *Session {
	url := auth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	session := &Session{
		AuthURL: url,
	}
	return session
}

func (auth *Auth) StoreInSession(key string, value string, req *http.Request, res http.ResponseWriter) error {
	session, _ := auth.Store.New(req, sessionName)

	if err := updateSessionValue(session, key, value); err != nil {
		return err
	}

	return session.Save(req, res)
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

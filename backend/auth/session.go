package auth

import (
	"github.com/gorilla/sessions"
	"github.com/kcsu/store/config"
)

func InitSessionStore(c *config.Config) sessions.Store {
	store := sessions.NewCookieStore([]byte(c.CookieSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   !c.Debug, // Secure in production
	}
	return store
}

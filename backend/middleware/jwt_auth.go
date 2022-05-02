package middleware

import (
	"errors"

	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Generate the error in the event of unauthenticated user
func jwtErrorHandler(e error) error {
	return echo.ErrUnauthorized
}

func extractJwt(c echo.Context) ([]string, error) {
	sess, err := session.Get("__session", c)
	if err != nil {
		return nil, err
	}
	token, ok := sess.Values["_token"].(string)
	if !ok {
		return nil, errors.New("no string token")
	}
	return []string{token}, nil
}

// Ensure a valid JWT is present in the Authorization header
func JWTAuth(c *config.Config) echo.MiddlewareFunc {
	jwtConfig := middleware.JWTConfig{
		Claims:           &auth.JwtClaims{},
		SigningKey:       []byte(c.JwtSecret),
		TokenLookupFuncs: []middleware.ValuesExtractor{extractJwt},
		ErrorHandler:     jwtErrorHandler,
	}
	return middleware.JWTWithConfig(jwtConfig)
}

package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/labstack/echo-contrib/session"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Generate the error in the event of unauthenticated user
func jwtErrorHandler(_ echo.Context, e error) error {
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
	jwtConfig := echojwt.Config{
		NewClaimsFunc:    func(echo.Context) jwt.Claims { return &auth.JwtClaims{} },
		SigningKey:       []byte(c.JwtSecret),
		TokenLookupFuncs: []middleware.ValuesExtractor{extractJwt},
		ErrorHandler:     jwtErrorHandler,
	}
	return echojwt.WithConfig(jwtConfig)
}

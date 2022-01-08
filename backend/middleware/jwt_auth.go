package middleware

import (
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/config"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
)

// Generate the error in the event of unauthenticated user
func jwtErrorHandler(e error) error {
	return echo.ErrUnauthorized
}

// Ensure a valid JWT is present in the Authorization header
func JWTAuth(c *config.Config) echo.MiddlewareFunc {
	jwtConfig := em.JWTConfig{
		Claims:       &auth.JwtClaims{},
		SigningKey:   []byte(c.JwtSecret),
		ErrorHandler: jwtErrorHandler,
	}
	return em.JWTWithConfig(jwtConfig)
}

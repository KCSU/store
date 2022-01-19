package auth

import (
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Custom claims for user JWT
type JwtClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
	// TODO: admin stuff
}

// Load claims from the current context
//
// Requires authentication middleware
func (auth *GoogleAuth) GetClaims(c echo.Context) *JwtClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtClaims)
	return claims
}

// Load the user's id from the current context
//
// Requires authentication middleware
func (auth *GoogleAuth) GetUserId(c echo.Context) int {
	claims := auth.GetClaims(c)
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		panic(err) // FIXME: hmmmmm...
	}
	return id
}

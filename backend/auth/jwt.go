package auth

import (
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
	// TODO: admin stuff
	// id? or is that iss/aud or something
}

func GetClaims(c echo.Context) *JwtClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtClaims)
	return claims
}

func GetUserId(c echo.Context) int {
	claims := GetClaims(c)
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		panic(err) // FIXME: hmmmmm...
	}
	return id
}

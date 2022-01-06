package auth

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
	// TODO: admin stuff
	// id? or is that iss/aud or something
}

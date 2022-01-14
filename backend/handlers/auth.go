package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kcsu/store/auth"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUser(c echo.Context) error {
	userId := auth.GetUserId(c)
	user, err := h.users.Find(userId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

// OAuth2 callback route handler
func (h *Handler) AuthCallback(c echo.Context) error {
	// Fetch the OAuth2 user data
	err := h.auth.VerifyGoogleCsrfToken(c)
	if err != nil {
		return err
	}

	authUser, err := h.auth.VerifyIdToken(c.FormValue("credential"), c)
	if err != nil {
		return err
	}

	// Create or fetch the user in the database
	user, err := h.users.FindOrCreate(authUser)
	if err != nil {
		// Ensure there is no email address conflict
		exists, exerr := h.users.Exists(authUser.Email)
		if exerr != nil {
			return exerr
		}
		if exists {
			return echo.NewHTTPError(http.StatusConflict, "email is taken")
		}
		return err
	}

	// Create JWT for login
	claims := &auth.JwtClaims{
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Subject:   strconv.Itoa(int(user.ID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte(h.config.JwtSecret))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "_token"
	cookie.Value = t
	// TODO: short-lived, refresh tokens
	cookie.Expires = time.Now().Add(24 * time.Hour)
	// Cookie is secure in production
	cookie.Secure = !h.config.Debug
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}

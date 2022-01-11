package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kcsu/store/auth"
	"github.com/labstack/echo/v4"
)

// OAuth2 callback route handler
func (h *Handler) AuthCallback(c echo.Context) error {
	// Fetch the OAuth2 user data
	gothUser, err := h.auth.CompleteUserAuth(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create or fetch the user in the database
	user, err := h.users.FindOrCreate(gothUser)
	if err != nil {
		// Ensure there is no email address conflict
		exists, exerr := h.users.Exists(gothUser.Email)
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

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// Redirect to the (google) OAuth2 provider
func (h *Handler) AuthRedirect(c echo.Context) error {
	url, err := h.auth.GetAuthUrl(c)
	if err != nil {
		return echo.ErrBadRequest
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}
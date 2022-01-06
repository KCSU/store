package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kcsu/store/auth"
	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthCallback(c echo.Context) error {
	gothUser, err := h.auth.CompleteUserAuth(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.users.FindOrCreate(gothUser)
	if err != nil {
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

func (h *Handler) AuthRedirect(c echo.Context) error {
	url, err := h.auth.GetAuthUrl(c)
	if err != nil {
		return echo.ErrBadRequest
	}
	// gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

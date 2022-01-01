package handlers

import (
	"net/http"

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
			return echo.ErrInternalServerError
		}
		if exists {
			return echo.NewHTTPError(http.StatusConflict, "email is taken")
		}
		return echo.ErrInternalServerError
	}
	// TODO: create JWT or session cookie for login
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) AuthRedirect(c echo.Context) error {
	url, err := h.auth.GetAuthUrl(c)
	if err != nil {
		return echo.ErrBadRequest
	}
	// gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

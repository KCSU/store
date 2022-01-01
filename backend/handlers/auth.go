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
	// TODO: handle user in database
	return c.JSON(http.StatusOK, gothUser)
}

func (h *Handler) AuthRedirect(c echo.Context) error {
	url, err := h.auth.GetAuthUrl(c)
	if err != nil {
		return echo.ErrBadRequest
	}
	// gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthCallback(c echo.Context) error {
	return nil
}

func (h *Handler) AuthRedirect(c echo.Context) error {
	url, err := h.auth.GetAuthUrl(c.Request(), c.Response())
	if err != nil {
		return echo.ErrBadRequest
	}
	// gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

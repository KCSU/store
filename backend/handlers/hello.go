package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (ah *AdminHandler) GetRoles(c echo.Context) error {
	roles, err := ah.Roles.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &roles)
}

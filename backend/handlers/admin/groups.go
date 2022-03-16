package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Fetch a list of all groups
func (ah *AdminHandler) GetGroups(c echo.Context) error {
	groups, err := ah.Groups.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &groups)
}
